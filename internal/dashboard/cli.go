package dashboard

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// Run implements `dh dashboard [--port N] [--stale D] [--done-decay D] [--detach]`.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dh dashboard", flag.ContinueOnError)
	fs.SetOutput(stderr)
	port := fs.Int("port", DefaultPort, "TCP port on 127.0.0.1")
	stale := fs.Duration("stale", DefaultStale, "a session whose snapshot is older than this is not shown as open")
	decay := fs.Duration("done-decay", snapshot.DefaultDoneDecay, "how long 'done' shows before reading as idle")
	detach := fs.Bool("detach", false, "start in the background, print the URL and exit (no-op if already running)")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: dh dashboard [--port N] [--stale D] [--done-decay D] [--detach]")
		fmt.Fprintln(stderr, "env: DH_DASHBOARD_ROOTS (folders separated by ';'), DH_DASHBOARD_SPRITES (optional sprite folder)")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() != 0 || *port < 1 || *port > 65535 || *stale <= 0 || *decay <= 0 {
		fs.Usage()
		return 1
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/", *port)
	if *detach {
		return startDetached(*port, url, childArgs(*port, *stale, *decay), stdout, stderr)
	}
	ln, err := Listen(*port)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, url)
	cfg := Config{Roots: os.Getenv(RootsEnv), SpritesDir: os.Getenv(SpritesEnv), SnapshotDir: snapshot.SnapshotDir(),
		Stale: *stale, DoneDecay: *decay}
	if err := Serve(ln, cfg); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

// isDashboard reports whether the port answers as a dh dashboard.
func isDashboard(port int) bool {
	c := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := c.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/api/state")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.Header.Get(markerHeader) != ""
}

// childArgs rebuilds the child argv from parsed values, so no spelling of --detach can survive into it.
func childArgs(port int, stale, decay time.Duration) []string {
	return []string{"dashboard", "--port", strconv.Itoa(port), "--stale", stale.String(), "--done-decay", decay.String()}
}

func startDetached(port int, url string, argv []string, stdout, stderr io.Writer) int {
	return startDetachedWith(port, url, argv, isDashboard, stdout, stderr)
}

func startDetachedWith(port int, url string, argv []string, probe func(int) bool, stdout, stderr io.Writer) int {
	if probe(port) {
		fmt.Fprintln(stdout, url)
		return 0
	}
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(stderr, "cannot locate own binary:", err)
		return 1
	}
	cmd := exec.Command(self, argv...)
	detachAttrs(cmd) // stdio stays nil = /dev/null
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(stderr, "cannot start dashboard:", err)
		return 1
	}
	exited := make(chan struct{})
	go func() { _ = cmd.Wait(); close(exited) }()
	for i := 0; i < 40; i++ {
		select {
		case <-exited: // child died: a concurrent --detach may have won the port
			if probe(port) {
				fmt.Fprintln(stdout, url)
				return 0
			}
			fmt.Fprintf(stderr, "dashboard did not start on port %d (busy by another program?)\n", port)
			return 1
		case <-time.After(100 * time.Millisecond):
		}
		if probe(port) {
			fmt.Fprintln(stdout, url)
			return 0
		}
	}
	fmt.Fprintln(stderr, "dashboard did not answer in time on port", port)
	return 1
}
