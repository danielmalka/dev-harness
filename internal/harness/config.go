package harness

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Config is the modeled part of <home>/config.yaml. Every other line is preserved on save, including a legacy `language:` line (the option no longer exists: ADR-011).
type Config struct {
	Mode         string            // default setup mode: "repo" or "global"
	Sprites      string            // dashboard.sprites
	Worktree     string            // "always" or "on-request"; not checked on load
	ReviewersRaw string            // the reviewers: block (header included), verbatim, never parsed
	Projects     map[string]string // absolute repository path -> name
}

const maxConfigBytes = 1 << 20

var keyRe = lazyRe(`^[A-Za-z_][A-Za-z0-9_-]*$`)

// seg is a top-level key line plus everything after it up to the next top-level key.
// The first seg may have an empty key (leading comments).
type seg struct {
	key   string
	val   string // inline scalar after "key:", raw
	line  int    // 1-based line number of the key
	lines []string
}

func cfgErr(line int, format string, a ...any) error {
	return fmt.Errorf("config.yaml linha %d: %s", line, fmt.Sprintf(format, a...))
}

func eol(s string) string { return strings.TrimRight(s, "\r\n") }

func isBlank(s string) bool {
	t := strings.TrimSpace(s)
	return t == "" || strings.HasPrefix(t, "#")
}

// split cuts the file into segments and rejects everything outside the supported subset.
func split(data string) ([]*seg, error) {
	var segs []*seg
	var cur *seg
	seen := map[string]bool{}
	raw := strings.SplitAfter(data, "\n")
	for i, l := range raw {
		if l == "" {
			continue
		}
		n := i + 1
		c := eol(l)
		switch {
		case isBlank(c):
		case c[0] == ' ' || c[0] == '\t':
			if strings.Contains(c[:len(c)-len(strings.TrimLeft(c, " \t"))], "\t") {
				return nil, cfgErr(n, "indentação com tab não é suportada")
			}
			if cur == nil || cur.key == "" {
				return nil, cfgErr(n, "linha indentada fora de um bloco")
			}
		default:
			idx := strings.Index(c, ":")
			if idx < 0 || !keyRe().MatchString(c[:idx]) {
				return nil, cfgErr(n, "esperado 'chave:' no nível de topo")
			}
			key := c[:idx]
			if seen[key] {
				return nil, cfgErr(n, "chave duplicada %q", key)
			}
			seen[key] = true
			cur = &seg{key: key, val: strings.TrimSpace(c[idx+1:]), line: n}
			segs = append(segs, cur)
		}
		if cur == nil {
			cur = &seg{}
			segs = append(segs, cur)
		}
		cur.lines = append(cur.lines, l)
	}
	return segs, nil
}

// scalar parses a quoted or plain single-line value (trailing comment allowed).
func scalar(s string, n int) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" || s[0] == '#' {
		return "", nil
	}
	switch s[0] {
	case '"':
		end := closeDouble(s)
		if end < 0 {
			return "", cfgErr(n, "aspas duplas sem fechar")
		}
		v, err := strconv.Unquote(s[:end+1])
		if err != nil {
			return "", cfgErr(n, "valor entre aspas duplas inválido")
		}
		return v, restOK(s[end+1:], n)
	case '\'':
		end := strings.Index(s[1:], "'")
		if end < 0 {
			return "", cfgErr(n, "aspas simples sem fechar")
		}
		return s[1 : end+1], restOK(s[end+2:], n)
	case '{', '[', '|', '>', '&', '*', '!':
		return "", cfgErr(n, "sintaxe de valor não suportada %q", s[:1])
	}
	if i := strings.Index(s, " #"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	if strings.Contains(s, ":") {
		return "", cfgErr(n, "o valor contém ':': ponha entre aspas")
	}
	return s, nil
}

func restOK(rest string, n int) error {
	if !isBlank(rest) {
		return cfgErr(n, "texto inesperado depois do valor entre aspas")
	}
	return nil
}

// closeDouble returns the index of the closing quote of s (s[0] == '"'), or -1.
func closeDouble(s string) int {
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return -1
}

// entry parses `  "path": name` (path quoted or plain) into path and name.
func entry(line string, n int) (path, name string, err error) {
	s := strings.TrimSpace(eol(line))
	var rest string
	switch {
	case strings.HasPrefix(s, `"`):
		end := closeDouble(s)
		if end < 0 {
			return "", "", cfgErr(n, "aspas duplas sem fechar")
		}
		if path, err = strconv.Unquote(s[:end+1]); err != nil {
			return "", "", cfgErr(n, "caminho entre aspas duplas inválido")
		}
		rest = s[end+1:]
	case strings.HasPrefix(s, "'"):
		end := strings.Index(s[1:], "'")
		if end < 0 {
			return "", "", cfgErr(n, "aspas simples sem fechar")
		}
		path, rest = s[1:end+1], s[end+2:]
	default:
		return "", "", cfgErr(n, "os caminhos de projeto devem estar entre aspas")
	}
	rest = strings.TrimLeft(rest, " ")
	if !strings.HasPrefix(rest, ":") {
		return "", "", cfgErr(n, "esperado ':' depois do caminho do projeto")
	}
	if name, err = scalar(rest[1:], n); err != nil {
		return "", "", err
	}
	if path == "" || name == "" {
		return "", "", cfgErr(n, "a entrada de projeto precisa de caminho e nome")
	}
	return path, name, nil
}

func parseConfig(data string) (Config, []*seg, error) {
	c := Config{}
	segs, err := split(data)
	if err != nil {
		return c, nil, err
	}
	for _, s := range segs {
		body := s.lines[1:]
		if s.key == "" {
			continue
		}
		switch s.key {
		case "mode", "worktree":
			v, err := scalar(s.val, s.line)
			if err != nil {
				return c, nil, err
			}
			for i, l := range body {
				if !isBlank(l) {
					return c, nil, cfgErr(s.line+1+i, "%s aceita um valor só", s.key)
				}
			}
			if s.key == "mode" && v != "" && v != "repo" && v != "global" {
				return c, nil, cfgErr(s.line, "mode deve ser repo ou global")
			}
			if s.key == "mode" {
				c.Mode = v
			} else {
				c.Worktree = v
			}
		case "dashboard":
			for i, l := range body {
				t := strings.TrimSpace(eol(l))
				if isBlank(l) || !strings.HasPrefix(t, "sprites:") {
					continue
				}
				if c.Sprites != "" {
					return c, nil, cfgErr(s.line+1+i, "chave duplicada \"sprites\"")
				}
				v, err := scalar(strings.TrimPrefix(t, "sprites:"), s.line+1+i)
				if err != nil {
					return c, nil, err
				}
				c.Sprites = v
			}
		case "reviewers":
			c.ReviewersRaw = strings.TrimRight(strings.Join(s.lines, ""), "\r\n \t")
			if s.val == "" && len(strings.TrimSpace(strings.Join(body, ""))) == 0 {
				c.ReviewersRaw = "" // header only, nothing to paste
			}
		case "projects":
			c.Projects = map[string]string{}
			if !isBlank(s.val) && s.val != "{}" {
				return c, nil, cfgErr(s.line, "projects deve ser um bloco de linhas \"caminho\": nome")
			}
			for i, l := range body {
				if isBlank(l) {
					continue
				}
				p, name, err := entry(l, s.line+1+i)
				if err != nil {
					return c, nil, err
				}
				if _, dup := c.Projects[p]; dup {
					return c, nil, cfgErr(s.line+1+i, "caminho de projeto duplicado %q", p)
				}
				c.Projects[p] = name
			}
		}
	}
	return c, segs, nil
}

func readConfig(home string) (string, error) {
	p := filepath.Join(home, "config.yaml")
	fi, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if !fi.Mode().IsRegular() || fi.Size() > maxConfigBytes {
		return "", fmt.Errorf("%s não é um arquivo comum de até %d bytes", p, maxConfigBytes)
	}
	b, err := os.ReadFile(p)
	return string(b), err
}

// LoadConfig reads <home>/config.yaml. A missing file is Config{} and nil; an invalid file
// is a readable error with a line number.
func LoadConfig(home string) (Config, error) {
	data, err := readConfig(home)
	if err != nil {
		return Config{}, err
	}
	c, _, err := parseConfig(data)
	if err != nil {
		return Config{}, err
	}
	if c.Projects == nil {
		c.Projects = map[string]string{}
	}
	return c, nil
}

// lazyRe compiles pattern on first use, so a package-level regexp costs nothing at init.
func lazyRe(pattern string) func() *regexp.Regexp {
	return sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(pattern) })
}
