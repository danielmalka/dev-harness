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

	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// Run implements `dh dashboard [--port N] [--stale D] [--done-decay D] [--detach|--stop]`.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dh dashboard", flag.ContinueOnError)
	fs.SetOutput(stderr)
	port := fs.Int("port", DefaultPort, "porta TCP em 127.0.0.1")
	stale := fs.Duration("stale", DefaultStale, "uma sessão cujo snapshot é mais velho que isto não aparece como aberta")
	decay := fs.Duration("done-decay", snapshot.DefaultDoneDecay, "por quanto tempo 'done' aparece antes de virar ocioso")
	detach := fs.Bool("detach", false, "inicia em segundo plano, imprime a URL e sai (não faz nada se já estiver rodando)")
	stop := fs.Bool("stop", false, "pede ao dashboard da --port que saia (nunca mata um processo)")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "uso: dh dashboard [--port N] [--stale D] [--done-decay D] [--detach | --stop]")
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
	if *stop && *detach {
		fmt.Fprintln(stderr, "--stop e --detach não podem ser usados juntos")
		fs.Usage()
		return 1
	}
	if *stop {
		return Stop(*port, stdout, stderr)
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
	home, err := harness.EnsureHome()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	cfg := Config{Home: home, SnapshotDir: snapshot.SnapshotDir(),
		Stale: *stale, DoneDecay: *decay, Port: *port}
	if tok, _, err := newStopToken(*port); err != nil {
		fmt.Fprintln(stderr, "aviso: --stop desativado, não foi possível gravar o token de parada:", err)
	} else {
		cfg.StopToken = tok
	}
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
		fmt.Fprintln(stderr, "não foi possível localizar o próprio binário:", err)
		return 1
	}
	cmd := exec.Command(self, argv...)
	detachAttrs(cmd) // stdio stays nil = /dev/null
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(stderr, "não foi possível iniciar o dashboard:", err)
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
			fmt.Fprintf(stderr, "o dashboard não iniciou na porta %d (ocupada por outro programa?)\n", port)
			return 1
		case <-time.After(100 * time.Millisecond):
		}
		if probe(port) {
			fmt.Fprintln(stdout, url)
			return 0
		}
	}
	fmt.Fprintln(stderr, "o dashboard não respondeu a tempo na porta", port)
	return 1
}
