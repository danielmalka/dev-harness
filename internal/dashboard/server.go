package dashboard

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

//go:embed web/index.html
var pageFS embed.FS

// Marker header lets --detach recognise a running dh dashboard on the port.
const markerHeader = "X-Dh-Dashboard"

const DefaultPort = 4747

type Config struct {
	Roots       string // DH_DASHBOARD_ROOTS value
	SpritesDir  string // DH_DASHBOARD_SPRITES value
	SnapshotDir string
	Stale       time.Duration
	DoneDecay   time.Duration
	Now         func() time.Time // tests only; nil = time.Now
	CacheTTL    time.Duration    // built state is reused this long (Serve sets 1s when zero)
}

// Handler is the whole HTTP surface: GET only, loopback Host only, no write path anywhere.
func Handler(cfg Config) http.Handler {
	sprites := LoadSprites(cfg.SpritesDir)
	cache := &stateCache{ttl: cfg.CacheTTL}
	page, _ := pageFS.ReadFile("web/index.html")
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
		b, ok := sprites.Frame(pose, idx) // Frame rejects any name outside PoseNames, so "../x" is 404
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(b)
	}
	return guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(page)
		case r.URL.Path == "/api/state":
			state(w, r)
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
		if r.Method != http.MethodGet {
			h.Set("Allow", "GET")
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
	Open      []jsonBar `json:"open"`
	Delivered []string  `json:"delivered"`
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

type state struct {
	Projects []jsonProject `json:"projects"`
	Sessions []jsonSession `json:"sessions"`
	Limits   jsonLimits    `json:"limits"`
	Avatar   jsonAvatar    `json:"avatar"`
	Warnings []string      `json:"warnings"`
}

func buildState(cfg Config, sprites *Sprites) state {
	now := time.Now()
	if cfg.Now != nil {
		now = cfg.Now()
	}
	st := state{Projects: []jsonProject{}, Sessions: []jsonSession{}, Warnings: []string{}}
	projects, warn := Projects(cfg.Roots)
	if warn != "" {
		st.Warnings = append(st.Warnings, warn)
	}
	st.Warnings = append(st.Warnings, sprites.Warnings...)
	for _, p := range projects {
		pr := ProjectProgress(p.Path)
		jp := jsonProject{Name: p.Name, Path: p.Path, Open: []jsonBar{}, Delivered: pr.Delivered}
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
	fr, _ := sprites.Frames(spec.Pose)
	st.Avatar = jsonAvatar{State: as, Pose: spec.Pose, FPS: spec.FPS, Frames: len(fr)}
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
	srv := &http.Server{Handler: Handler(cfg), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
