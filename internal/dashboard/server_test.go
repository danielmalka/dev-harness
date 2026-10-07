package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testHandler(t *testing.T, sprites string) http.Handler {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "proj", ".harness"), 0o755); err != nil {
		t.Fatal(err)
	}
	return Handler(Config{Roots: root, SpritesDir: sprites, SnapshotDir: t.TempDir(),
		Stale: DefaultStale, DoneDecay: 5 * time.Minute})
}

func do(h http.Handler, method, host, path string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://"+host+path, nil)
	r.Host = host
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestHostGuard(t *testing.T) {
	h := testHandler(t, "")
	for _, host := range []string{"exemplo.com", "exemplo.com:7777", "127.0.0.1.evil.com", "0.0.0.0:7777", "[::1]:7777", ""} {
		if c := do(h, "GET", host, "/api/state").Code; c < 400 || c > 499 {
			t.Errorf("host %q: got %d, want 4xx", host, c)
		}
	}
	for _, host := range []string{"127.0.0.1", "127.0.0.1:7777", "localhost", "LOCALHOST:7777"} {
		if c := do(h, "GET", host, "/api/state").Code; c != 200 {
			t.Errorf("host %q: got %d, want 200", host, c)
		}
	}
}

func TestMethodsAreGetOnly(t *testing.T) {
	h := testHandler(t, "")
	for _, p := range []string{"/", "/api/state", "/sprite/frente", "/sprite/../x", "/nada"} {
		for _, m := range []string{"POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
			if c := do(h, m, "127.0.0.1:7777", p).Code; c != 405 {
				t.Errorf("%s %s: got %d, want 405", m, p, c)
			}
		}
	}
}

func TestStateHasAllBlocks(t *testing.T) {
	w := do(testHandler(t, ""), "GET", "127.0.0.1:7777", "/api/state")
	var m map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"projects", "sessions", "limits", "avatar", "warnings"} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing %q in %s", k, w.Body)
		}
	}
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Error("content type")
	}
}

func TestPageAndHeaders(t *testing.T) {
	w := do(testHandler(t, ""), "GET", "localhost:7777", "/")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "<html") {
		t.Fatalf("page: %d", w.Code)
	}
	if w.Header().Get("X-Frame-Options") != "" {
		t.Error("X-Frame-Options must be absent (R5)")
	}
	if strings.Contains(strings.ToLower(w.Header().Get("Content-Security-Policy")), "frame-ancestors") {
		t.Error("frame-ancestors must be absent (R5)")
	}
	if w.Header().Get(markerHeader) == "" {
		t.Error("marker header")
	}
	if c := do(testHandler(t, ""), "GET", "localhost:7777", "/nada").Code; c != 404 {
		t.Errorf("unknown path: %d", c)
	}
}

func TestSpriteRoute(t *testing.T) {
	partial := t.TempDir()
	if err := os.WriteFile(filepath.Join(partial, "duvida.png"), []byte("OWNER"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, dir := range map[string]string{"none": "", "empty": t.TempDir(), "partial": partial} {
		h := testHandler(t, dir)
		for _, pose := range []string{"frente", "duvida", "puto"} {
			w := do(h, "GET", "127.0.0.1:7777", "/sprite/"+pose)
			if w.Code != 200 || w.Header().Get("Content-Type") != "image/png" {
				t.Errorf("%s/%s: %d", name, pose, w.Code)
			}
		}
		for _, bad := range []string{"../x", "..%2fx", "x", "frente/..", "frente.png", "", "frente?f=9", "frente?f=a", "%2e%2e/x"} {
			if c := do(h, "GET", "127.0.0.1:7777", "/sprite/"+bad).Code; c != 404 && c != 301 {
				t.Errorf("%s: /sprite/%s got %d, want 404", name, bad, c)
			}
		}
	}
	w := do(testHandler(t, partial), "GET", "127.0.0.1:7777", "/sprite/duvida")
	if w.Body.String() != "OWNER" {
		t.Error("owner folder must win over the kit example")
	}
}

func TestListenBusy(t *testing.T) {
	ln, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	if host := ln.Addr().String(); !strings.HasPrefix(host, "127.0.0.1:") {
		t.Fatalf("bound to %s", host)
	}
	port := ln.Addr().(interface{ String() string }).String()
	var p int
	for i := strings.LastIndex(port, ":") + 1; i < len(port); i++ {
		p = p*10 + int(port[i]-'0')
	}
	if _, err := Listen(p); err == nil || !strings.Contains(err.Error(), "busy") {
		t.Fatalf("want busy error, got %v", err)
	}
}
