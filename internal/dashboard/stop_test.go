package dashboard

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// startServing runs Serve on a free port with a published stop token; returns the port and the Serve result channel.
func startServing(t *testing.T, roots string) (int, string, chan error) {
	t.Helper()
	cfgDir := t.TempDir()
	userConfigDir = func() (string, error) { return cfgDir, nil }
	t.Cleanup(func() { userConfigDir = os.UserConfigDir })
	ln, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	tok, _, err := newStopToken(port)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- Serve(ln, Config{Roots: roots, SnapshotDir: t.TempDir(), Stale: DefaultStale, DoneDecay: time.Minute, Port: port, StopToken: tok})
	}()
	return port, tok, done
}

func post(t *testing.T, port int, token string, hdr map[string]string) int {
	t.Helper()
	req, _ := http.NewRequest("POST", "http://127.0.0.1:"+itoa(port)+"/api/stop", nil)
	if token != "" {
		req.Header.Set(stopHeader, token)
	}
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

func itoa(n int) string { return strconv.Itoa(n) }

func TestConfigInState(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	get := func(cfg Config) jsonConfig {
		cfg.SnapshotDir = t.TempDir()
		var st struct{ Config jsonConfig }
		if err := json.Unmarshal(do(Handler(cfg), "GET", "127.0.0.1:1", "/api/state").Body.Bytes(), &st); err != nil {
			t.Fatal(err)
		}
		return st.Config
	}
	c := get(Config{Roots: " " + a + ";;" + b + " ", SpritesDir: "", Port: 4799})
	if len(c.Roots) != 2 || c.Roots[0] != a || c.Roots[1] != b || c.Sprites != "" || c.Port != 4799 {
		t.Fatalf("config: %+v", c)
	}
	w := do(Handler(Config{SnapshotDir: t.TempDir()}), "GET", "127.0.0.1:1", "/api/state")
	if !strings.Contains(w.Body.String(), `"roots":[]`) {
		t.Fatalf("roots must be [] when unset: %s", w.Body.String())
	}
	if c := get(Config{SpritesDir: "rel/sprites"}); !filepath.IsAbs(c.Sprites) {
		t.Fatalf("sprites not absolute: %q", c.Sprites)
	}
}

func TestStopGuards(t *testing.T) {
	h := Handler(Config{SnapshotDir: t.TempDir(), StopToken: "secret"})
	if c := do(h, "GET", "127.0.0.1:1", "/api/stop").Code; c != 405 {
		t.Errorf("GET /api/stop: %d", c)
	}
	if c := do(h, "OPTIONS", "127.0.0.1:1", "/api/stop").Code; c != 405 {
		t.Errorf("OPTIONS /api/stop: %d", c)
	}
	if w := do(h, "OPTIONS", "127.0.0.1:1", "/api/stop"); w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("CORS header present")
	}
	for _, p := range []string{"/", "/api/state", "/sprite/frente", "/nada"} {
		if c := do(h, "POST", "127.0.0.1:1", p).Code; c != 405 {
			t.Errorf("POST %s: %d", p, c)
		}
	}
	if c := do(h, "POST", "evil.example", "/api/stop").Code; c != 403 {
		t.Errorf("foreign Host: %d", c)
	}
	// empty configured token never matches an empty header
	if c := do(Handler(Config{SnapshotDir: t.TempDir()}), "POST", "127.0.0.1:1", "/api/stop").Code; c != 403 {
		t.Errorf("no token configured: %d", c)
	}
}

func TestStopTokenAndShutdown(t *testing.T) {
	port, tok, done := startServing(t, "")
	if c := post(t, port, "", nil); c != 403 {
		t.Fatalf("missing token: %d", c)
	}
	if c := post(t, port, "wrong", nil); c != 403 {
		t.Fatalf("wrong token: %d", c)
	}
	if c := post(t, port, tok, map[string]string{"Origin": "http://evil.example"}); c != 403 {
		t.Fatalf("origin with right token: %d", c)
	}
	if c := post(t, port, tok, nil); c != 204 {
		t.Fatalf("right token: %d", c)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not exit")
	}
	if portOpen("127.0.0.1:" + itoa(port)) {
		t.Fatal("port still open")
	}
}

func TestStopCommand(t *testing.T) {
	port, _, done := startServing(t, "")
	var out, errb bytes.Buffer
	if rc := Stop(port, &out, &errb); rc != 0 || !strings.Contains(out.String(), "stopped dashboard on port") {
		t.Fatalf("rc=%d out=%q err=%q", rc, out.String(), errb.String())
	}
	<-done
	out.Reset()
	if rc := Stop(port, &out, &errb); rc != 0 || !strings.Contains(out.String(), "no dashboard on port") {
		t.Fatalf("idempotent: rc=%d out=%q", rc, out.String())
	}
}

func portOf(t *testing.T, srv *httptest.Server) int {
	t.Helper()
	return srv.Listener.Addr().(*net.TCPAddr).Port
}

func TestStopRefusesForeignAndOldServers(t *testing.T) {
	plain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer plain.Close()
	old := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(markerHeader, "1")
		if r.URL.Path == "/api/stop" {
			http.NotFound(w, r)
		}
	}))
	defer old.Close()
	cfgDir := t.TempDir()
	userConfigDir = func() (string, error) { return cfgDir, nil }
	t.Cleanup(func() { userConfigDir = os.UserConfigDir })
	if _, _, err := newStopToken(portOf(t, old)); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if rc := Stop(portOf(t, plain), &out, &errb); rc != 1 || !strings.Contains(errb.String(), "held by another program") {
		t.Fatalf("plain: rc=%d err=%q", rc, errb.String())
	}
	errb.Reset()
	if rc := Stop(portOf(t, old), &out, &errb); rc != 1 || !strings.Contains(errb.String(), "predates --stop") {
		t.Fatalf("old: rc=%d err=%q", rc, errb.String())
	}
}

func TestRestartOverwritesToken(t *testing.T) {
	cfgDir := t.TempDir()
	userConfigDir = func() (string, error) { return cfgDir, nil }
	t.Cleanup(func() { userConfigDir = os.UserConfigDir })
	t1, p, err := newStopToken(4801)
	if err != nil {
		t.Fatal(err)
	}
	t2, _, err := newStopToken(4801)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(p)
	if t1 == t2 || string(got) != t2 {
		t.Fatalf("restart must overwrite with a new token: %q %q %q", t1, t2, got)
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(p), ".stop-*")); len(left) != 0 {
		t.Fatalf("temp files left: %v", left)
	}
}

func TestTokenPublishIgnoresPlantedSymlinkAndTightensDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no unix symlink/mode semantics")
	}
	cfgDir := t.TempDir()
	userConfigDir = func() (string, error) { return cfgDir, nil }
	t.Cleanup(func() { userConfigDir = os.UserConfigDir })
	p, _ := tokenPath(4802)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	_ = os.Chmod(dir, 0o777)
	victim := filepath.Join(t.TempDir(), "victim")
	for _, n := range []string{"stop-4802.token.tmp", ".stop-1.tmp", ".stop-0.tmp"} {
		_ = os.Symlink(victim, filepath.Join(dir, n))
	}
	if _, _, err := newStopToken(4802); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(victim); err == nil {
		t.Fatal("wrote through a planted symlink")
	}
	if fi, _ := os.Stat(dir); fi.Mode().Perm() != 0o700 {
		t.Fatalf("dir mode %v", fi.Mode().Perm())
	}
	if fi, _ := os.Stat(p); fi.Mode().Perm() != 0o600 {
		t.Fatalf("file mode %v", fi.Mode().Perm())
	}
}

func TestStopRefusesSymlinkedTokenFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need privileges on windows")
	}
	port, tok, done := startServing(t, "")
	p, _ := tokenPath(port)
	real := p + ".real"
	_ = os.Rename(p, real)
	if err := os.Symlink(real, p); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	if rc := Stop(port, &out, &errb); rc != 1 || !strings.Contains(errb.String(), "stop token not found") {
		t.Fatalf("rc=%d err=%q", rc, errb.String())
	}
	_ = post(t, port, tok, nil)
	<-done
}

func TestDoubleStopNoPanic(t *testing.T) {
	port, tok, done := startServing(t, "")
	res := make(chan int, 2)
	for i := 0; i < 2; i++ {
		go func() {
			req, _ := http.NewRequest("POST", "http://127.0.0.1:"+itoa(port)+"/api/stop", nil)
			req.Header.Set(stopHeader, tok)
			if resp, err := http.DefaultClient.Do(req); err == nil {
				resp.Body.Close()
				res <- resp.StatusCode
			} else {
				res <- 0 // second request may lose the race with the listener closing
			}
		}()
	}
	<-res
	<-res
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not exit")
	}
}

func TestOriginWithValidTokenKeepsServerUp(t *testing.T) {
	port, tok, done := startServing(t, "")
	if c := post(t, port, tok, map[string]string{"Origin": "http://evil.example"}); c != 403 {
		t.Fatalf("got %d", c)
	}
	select {
	case <-done:
		t.Fatal("server stopped")
	case <-time.After(200 * time.Millisecond):
	}
	_ = post(t, port, tok, nil)
	<-done
}

func TestStopWithoutTokenFile(t *testing.T) {
	port, tok, done := startServing(t, "")
	p, _ := tokenPath(port)
	_ = os.Remove(p)
	var out, errb bytes.Buffer
	if rc := Stop(port, &out, &errb); rc != 1 || !strings.Contains(errb.String(), "stop token not found") {
		t.Fatalf("rc=%d err=%q", rc, errb.String())
	}
	_ = post(t, port, tok, nil) // stop with the server's real token
	<-done
}

func TestStopWithRejectedToken(t *testing.T) {
	port, tok, done := startServing(t, "")
	p, _ := tokenPath(port)
	_ = os.WriteFile(p, []byte("not-the-token"), 0o600)
	var out, errb bytes.Buffer
	if rc := Stop(port, &out, &errb); rc != 1 || !strings.Contains(errb.String(), "stop token rejected") {
		t.Fatalf("rc=%d err=%q", rc, errb.String())
	}
	_ = post(t, port, tok, nil) // stop with the server's real token
	<-done
}
