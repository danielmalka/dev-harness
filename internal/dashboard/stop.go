package dashboard

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// tokenPath is the stop-token file of the dashboard on port: <home>/dashboard/stop-<port>.token
// (0700 dir, 0600 file); on Windows the mode bits do not restrict access, but the user profile does.
func tokenPath(port int) (string, error) {
	dir := harness.DashboardDir()
	if dir == "" {
		return "", errors.New("cannot determine the harness home: set DH_HOME, HOME or USERPROFILE")
	}
	return filepath.Join(dir, "stop-"+strconv.Itoa(port)+".token"), nil
}

// newStopToken creates a fresh token and publishes it atomically. The temp file is random-named and O_EXCL
// (os.CreateTemp), so a planted symlink is never written through; the dir is forced to 0700 in case it pre-existed
// looser. The file is never deleted on exit: a stale 0600 file is harmless (the next start on the port overwrites it,
// and no server accepts a stale value), and not deleting avoids a stopping process erasing its successor's token.
func newStopToken(port int) (token, path string, err error) {
	if path, err = tokenPath(port); err != nil {
		return "", "", err
	}
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return "", "", err
	}
	if err = os.Chmod(dir, 0o700); err != nil { // repairs a pre-existing lax dir; no-op on Windows
		return "", "", err
	}
	f, err := os.CreateTemp(dir, ".stop-*.tmp")
	if err != nil {
		return "", "", err
	}
	tmp := f.Name()
	if _, err = f.WriteString(token); err == nil {
		err = f.Chmod(0o600)
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(tmp, path)
	}
	if err != nil {
		_ = os.Remove(tmp)
		return "", "", err
	}
	return token, path, nil
}

// Stop implements `dh dashboard --stop`: asks the dashboard on port to exit; never kills a process.
func Stop(port int, stdout, stderr io.Writer) int {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	if !portOpen(addr) {
		fmt.Fprintf(stdout, "no dashboard on port %d\n", port)
		return 0
	}
	c := http.Client{Timeout: 2 * time.Second}
	other := func() int {
		fmt.Fprintf(stderr, "port %d is held by another program, not a dh dashboard\n", port)
		return 1
	}
	resp, err := c.Get("http://" + addr + "/api/state")
	if err != nil {
		return other()
	}
	resp.Body.Close()
	if resp.Header.Get(markerHeader) == "" {
		return other()
	}
	path, err := tokenPath(port)
	if err != nil {
		fmt.Fprintln(stderr, "cannot locate the stop token:", err)
		return 1
	}
	tok, err := snapshot.ReadCapped(path) // refuses symlink, FIFO and oversize
	if err != nil {
		fmt.Fprintf(stderr, "stop token not found (%s): the dashboard on port %d was not started by this user or predates --stop; stop it manually\n", path, port)
		return 1
	}
	req, _ := http.NewRequest(http.MethodPost, "http://"+addr+"/api/stop", nil)
	req.Header.Set(stopHeader, strings.TrimSpace(string(tok)))
	resp, err = c.Do(req)
	if err != nil {
		fmt.Fprintln(stderr, "stop request failed:", err)
		return 1
	}
	resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNoContent:
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		fmt.Fprintf(stderr, "the dh dashboard on port %d predates --stop (0.19.x or older); stop it manually once\n", port)
		return 1
	case http.StatusForbidden:
		fmt.Fprintln(stderr, "stop token rejected")
		return 1
	default:
		fmt.Fprintf(stderr, "unexpected answer to stop: %s\n", resp.Status)
		return 1
	}
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		if !portOpen(addr) {
			fmt.Fprintf(stdout, "stopped dashboard on port %d\n", port)
			return 0
		}
	}
	fmt.Fprintf(stderr, "dashboard on port %d did not stop in time\n", port)
	return 1
}

func portOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
