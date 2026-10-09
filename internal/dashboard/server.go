package dashboard

import (
	"context"
	"crypto/subtle"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

//go:embed web/index.html web/metrics.js
var pageFS embed.FS

// Marker header lets --detach recognise a running dh dashboard on the port.
const markerHeader = "X-Dh-Dashboard"

const DefaultPort = 4747

// stopHeader carries the stop token (custom header, so browsers must preflight it).
const stopHeader = "X-Dh-Stop-Token"

type Config struct {
	Home        string // harness home; config.yaml is re-read on every state build
	SnapshotDir string
	Stale       time.Duration
	DoneDecay   time.Duration
	Now         func() time.Time // tests only; nil = time.Now
	CacheTTL    time.Duration    // built state is reused this long (Serve sets 1s when zero)
	Port        int              // listening port, reported in /api/state config
	StopToken   string           // secret for POST /api/stop; empty = stop disabled (always 403)
	OnStop      func()           // called after the 204 of an accepted stop (Serve wires it to Shutdown)
}

// Handler is the whole HTTP surface: loopback Host only, GET only except one write path, POST /api/stop.
// That path only asks the process to exit, and is safe against a hostile web page because: it needs the
// X-Dh-Stop-Token header, whose value lives in a file only the user can read; a custom header forces a CORS
// preflight (OPTIONS), which is always 405 and never carries CORS headers, so a browser never sends the POST;
// the loopback Host check blocks DNS rebinding; and any request carrying an Origin header (browsers always send
// one, the CLI never does) is refused even with a valid token.
func Handler(cfg Config) http.Handler {
	sprites := &spriteHolder{}
	sprites.refresh(func() string { d, _ := spritesDirOf(cfg.Home); return d }())
	cache := &stateCache{ttl: cfg.CacheTTL}
	page, _ := pageFS.ReadFile("web/index.html")
	metricsJS, _ := pageFS.ReadFile("web/metrics.js")
	memory := memoryH(cfg)
	// own routing instead of ServeMux: the mux 307-redirects unclean paths ("/sprite/../x"); here they are plain 404s
	state := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(cache.get(func() state { return buildState(cfg, sprites) }))
	}
	spriteH := func(w http.ResponseWriter, r *http.Request) {
		pose := strings.TrimPrefix(r.URL.Path, "/sprite/")
		idx := 0
		if f := r.URL.Query().Get("f"); f != "" {
			var err error
			if idx, err = strconv.Atoi(f); err != nil {
				http.NotFound(w, r)
				return
			}
		}
		b, ok := sprites.current().Frame(pose, idx) // Frame rejects any name outside PoseNames, so "../x" is 404
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(b)
	}
	stop := func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(stopHeader)
		// defence in depth: browsers always send Origin, the CLI never does
		if r.Header.Get("Origin") != "" || cfg.StopToken == "" || subtle.ConstantTimeCompare([]byte(got), []byte(cfg.StopToken)) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusNoContent) // body is never read
		if cfg.OnStop != nil {
			go cfg.OnStop()
		}
	}
	return guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/stop":
			stop(w, r)
		case r.URL.Path == "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(page)
		case r.URL.Path == "/api/state":
			state(w, r)
		case r.URL.Path == "/api/memory":
			memory(w, r)
		case r.URL.Path == "/metrics.js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache")
			_, _ = w.Write(metricsJS)
		case strings.HasPrefix(r.URL.Path, "/sprite/"):
			spriteH(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

// guard enforces loopback Host (DNS rebinding) and GET only, and sets the common headers.
// ponytail: no X-Frame-Options and no frame-ancestors on purpose (PRD-012 R5, iframe/webview embed).
func guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set(markerHeader, "1")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'")
		if !loopbackHost(r.Host) {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		want := http.MethodGet
		if r.URL.Path == "/api/stop" {
			want = http.MethodPost
		}
		if r.Method != want {
			h.Set("Allow", want)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loopbackHost compares the host (port stripped) against the exact list; [::1] is deliberately not accepted.
func loopbackHost(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	host = strings.ToLower(host)
	return host == "127.0.0.1" || host == "localhost"
}

// stateCache keeps the built state for ttl so rapid or concurrent polls do not rescan the disk.
type stateCache struct {
	mu   sync.Mutex
	ttl  time.Duration
	at   time.Time
	last *state
}

func (c *stateCache) get(build func() state) state {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.last != nil && time.Since(c.at) < c.ttl {
		return *c.last
	}
	s := build()
	c.last, c.at = &s, time.Now()
	return s
}

type jsonBar struct {
	PRD       string `json:"prd"`
	Done      int    `json:"done"`
	Blocked   int    `json:"blocked"`
	Total     int    `json:"total"`
	NoTickets bool   `json:"no_tickets"`
}

type jsonProject struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Mode      string    `json:"mode"`
	Harness   string    `json:"harness"`
	Open      []jsonBar `json:"open"`
	Delivered []string  `json:"delivered"`
	Metrics   Metrics   `json:"metrics"`
}

type jsonSession struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CWD        string `json:"cwd"`
	Agent      string `json:"agent"`
	Project    string `json:"project"`
	Activity   string `json:"activity"`
	ActivityAt string `json:"activity_at,omitempty"`
	UpdatedAt  string `json:"updated_at"`
}

type jsonLimits struct {
	OK             bool     `json:"ok"`
	FiveHour       *float64 `json:"five_hour"`
	SevenDay       *float64 `json:"seven_day"`
	FiveHourAgeSec int64    `json:"five_hour_age_sec"`
	SevenDayAgeSec int64    `json:"seven_day_age_sec"`
}

type jsonAvatar struct {
	State  string `json:"state"`
	Pose   string `json:"pose"`
	FPS    int    `json:"fps"`
	Frames int    `json:"frames"`
}

type jsonConfig struct {
	Home    string `json:"home"`
	Sprites string `json:"sprites"`
	Port    int    `json:"port"`
}

type state struct {
	Config   jsonConfig    `json:"config"`
	Projects []jsonProject `json:"projects"`
	Sessions []jsonSession `json:"sessions"`
	Limits   jsonLimits    `json:"limits"`
	Avatar   jsonAvatar    `json:"avatar"`
	Warnings []string      `json:"warnings"`
}

// spritesDirOf reads dashboard.sprites from <home>/config.yaml. A leading "~/" expands to the user home; the
// result must be an absolute existing directory, else "" (default set) plus a warning. Unset or unreadable config = "".
func spritesDirOf(home string) (string, string) {
	if home == "" {
		return "", ""
	}
	c, err := harness.LoadConfig(home) // an invalid file is reported by harness.Projects
	if err != nil || c.Sprites == "" {
		return "", ""
	}
	v := c.Sprites
	if strings.HasPrefix(v, "~/") {
		if uh, err := os.UserHomeDir(); err == nil && uh != "" {
			v = filepath.Join(uh, v[2:])
		}
	}
	if fi, err := os.Stat(v); !filepath.IsAbs(v) || err != nil || !fi.IsDir() {
		return "", fmt.Sprintf("dashboard.sprites %q is not an absolute existing directory: using the default set", c.Sprites)
	}
	return filepath.Clean(v), ""
}

// spriteHolder reloads the sprite set only when dashboard.sprites changes between state builds.
type spriteHolder struct {
	mu  sync.Mutex
	dir string
	s   *Sprites
}

func (h *spriteHolder) refresh(dir string) *Sprites {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.s == nil || dir != h.dir {
		h.dir, h.s = dir, LoadSprites(dir)
	}
	return h.s
}

func (h *spriteHolder) current() *Sprites {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.s
}

func buildState(cfg Config, holder *spriteHolder) state {
	now := time.Now()
	if cfg.Now != nil {
		now = cfg.Now()
	}
	spDir, spWarn := spritesDirOf(cfg.Home)
	sprites := holder.refresh(spDir)
	st := state{Config: jsonConfig{Home: cfg.Home, Sprites: spDir, Port: cfg.Port}, Projects: []jsonProject{}, Sessions: []jsonSession{}, Warnings: []string{}}
	var projects []harness.Project
	if cfg.Home == "" {
		st.Warnings = append(st.Warnings, "harness home unknown: no projects to show.")
	} else {
		var warns []string
		projects, warns = harness.Projects(cfg.Home)
		st.Warnings = append(st.Warnings, warns...)
	}
	if spWarn != "" {
		st.Warnings = append(st.Warnings, spWarn)
	}
	st.Warnings = append(st.Warnings, sprites.Warnings...)
	for _, p := range projects {
		pr := projectProgress(p.Name, p.Path, p.Harness)
		jp := jsonProject{Name: p.Name, Path: p.Path, Mode: p.Mode, Harness: p.Harness, Open: []jsonBar{}, Delivered: pr.Delivered, Metrics: pr.Metrics}
		pr.Warnings = append(pr.Warnings, fileMetrics(&jp.Metrics, p.Name, p.Harness)...)
		if jp.Delivered == nil {
			jp.Delivered = []string{}
		}
		for _, b := range pr.Open {
			jp.Open = append(jp.Open, jsonBar(b))
		}
		st.Warnings = append(st.Warnings, pr.Warnings...)
		st.Projects = append(st.Projects, jp)
	}
	all, skipped := snapshot.LoadDirMax(cfg.SnapshotDir, MaxSnapshots)
	if skipped > 0 {
		st.Warnings = append(st.Warnings, fmt.Sprintf("More than %d snapshot files: reading the newest %d.", MaxSnapshots, MaxSnapshots))
	}
	open := OpenSessions(all, projects, now, cfg.Stale, cfg.DoneDecay)
	for _, s := range open {
		js := jsonSession{ID: s.SessionID, Name: s.SessionName, CWD: s.CWD, Agent: s.Agent, Project: s.Project,
			Activity: s.Activity, UpdatedAt: s.UpdatedAt.UTC().Format(time.RFC3339)}
		if !s.ActivityAt.IsZero() {
			js.ActivityAt = s.ActivityAt.UTC().Format(time.RFC3339)
		}
		st.Sessions = append(st.Sessions, js)
	}
	lim := LatestLimits(all, now)
	st.Limits = jsonLimits{OK: lim.OK, FiveHour: lim.FiveHour, SevenDay: lim.SevenDay,
		FiveHourAgeSec: int64(lim.FiveHourAge.Seconds()), SevenDayAgeSec: int64(lim.SevenDayAge.Seconds())}
	as := AvatarState(open, lim)
	spec := sprites.PoseFor(as)
	st.Avatar = jsonAvatar{State: as, Pose: spec.Pose, FPS: spec.FPS, Frames: sprites.FrameCount(spec.Pose)}
	return st
}

// Listen binds 127.0.0.1 only; a busy port gives a readable error.
func Listen(port int) (net.Listener, error) {
	ln, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	if err != nil {
		// ponytail: message match covers Windows, where EADDRINUSE is not the socket error code
		if errors.Is(err, syscall.EADDRINUSE) || strings.Contains(err.Error(), "in use") || strings.Contains(err.Error(), "Only one usage") {
			return nil, fmt.Errorf("port %d is busy: pick another with --port", port)
		}
		return nil, fmt.Errorf("cannot listen on 127.0.0.1:%d: %w", port, err)
	}
	return ln, nil
}

// Serve blocks until ln closes.
func Serve(ln net.Listener, cfg Config) error {
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = time.Second
	}
	var srv *http.Server
	done := make(chan struct{})
	var once sync.Once
	cfg.OnStop = func() {
		once.Do(func() {
			defer close(done)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
		})
	}
	srv = &http.Server{Handler: Handler(cfg), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		<-done // Serve returns at once on Shutdown; wait so the 204 is flushed before the process exits
		return nil
	}
	return err
}
