package dashboard

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

const SpritesEnv = "DH_DASHBOARD_SPRITES"

//go:embed web/sprites/*.png
var embedded embed.FS

// PoseNames are jevmon's 19 poses; only these names are ever served.
var PoseNames = []string{"frente", "caneca", "animado", "notebook", "confiante", "assustado", "puto", "triste",
	"joinha", "feliz", "costas", "lateral", "fone", "celular", "duvida", "costas2", "mostrando", "explicando", "apontando"}

type PoseSpec struct {
	Pose string `json:"pose"`
	FPS  int    `json:"fps"`
}

var defaultTable = map[string]PoseSpec{
	StateWaiting: {Pose: "duvida"}, StateError: {Pose: "puto"}, StateWorking: {Pose: "celular"},
	StateDone: {Pose: "mostrando"}, StateAttn: {Pose: "apontando"}, StateIdle: {Pose: "frente"},
}

type Sprites struct {
	dir      string
	table    map[string]PoseSpec
	Warnings []string
}

// LoadSprites reads the optional poses.json of dir (jevmon format {"<state>":{"pose":..,"fps":0..60}}).
// An invalid file is ignored with a warning and the default map applies. dir == "" uses only the embedded set.
func LoadSprites(dir string) *Sprites {
	s := &Sprites{dir: dir, table: map[string]PoseSpec{}}
	for k, v := range defaultTable {
		s.table[k] = v
	}
	if dir == "" {
		return s
	}
	b, err := readCapped(filepath.Join(dir, "poses.json"))
	if errors.Is(err, os.ErrNotExist) {
		return s // absent is normal
	}
	if err == errTooLarge {
		s.Warnings = append(s.Warnings, "poses.json over 1 MiB, using default map")
		return s
	}
	if err != nil {
		s.Warnings = append(s.Warnings, "poses.json unreadable, using default map")
		return s
	}
	var over map[string]PoseSpec
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&over); err != nil {
		s.Warnings = append(s.Warnings, "poses.json invalid JSON, using default map")
		return s
	}
	for st, spec := range over {
		_, ours := statePriority[st]
		if !ours && !jevmonOnlyState[st] {
			s.Warnings = append(s.Warnings, fmt.Sprintf("poses.json unknown state %q, using default map", st))
			return s
		}
		if !KnownPose(spec.Pose) || spec.FPS < 0 || spec.FPS > 60 {
			s.Warnings = append(s.Warnings, fmt.Sprintf("poses.json invalid entry for %q, using default map", st))
			return s
		}
	}
	for st, spec := range over {
		if _, ours := statePriority[st]; ours {
			s.table[st] = spec
		}
	}
	return s
}

// jevmon states the dashboard does not use; accepted so a shared poses.json stays valid.
var jevmonOnlyState = map[string]bool{"sem_conexao": true, "silenciado": true, "voz_escuta": true, "voz_fala": true}

func KnownPose(p string) bool {
	for _, n := range PoseNames {
		if n == p {
			return true
		}
	}
	return false
}

// PoseFor returns the pose and fps for an avatar state (idle if the state is unknown).
func (s *Sprites) PoseFor(state string) PoseSpec {
	if p, ok := s.table[state]; ok {
		return p
	}
	return s.table[StateIdle]
}

var frameName = regexp.MustCompile(`^(.+)_(\d\d)\.png$`)

// Frames returns the PNG frames of a known pose: the owner's folder first (<pose>.png or <pose>_00.png...),
// then the kit's example of the same name, then frente. false for an unknown name (e.g. "../x") or no file anywhere.
func (s *Sprites) Frames(pose string) ([][]byte, bool) {
	if !KnownPose(pose) {
		return nil, false
	}
	if fr, ok := s.frames(pose); ok {
		return fr, true
	}
	if pose != "frente" { // jevmon behaviour: a pose with no file anywhere shows frente
		return s.frames("frente")
	}
	return nil, false
}

func (s *Sprites) frames(pose string) ([][]byte, bool) {
	if s.dir != "" {
		if b, err := readCapped(filepath.Join(s.dir, pose+".png")); err == nil {
			return [][]byte{b}, true
		}
		files, _ := filepath.Glob(filepath.Join(s.dir, pose+"_[0-9][0-9].png"))
		sort.Strings(files)
		var out [][]byte
		for _, f := range files {
			if m := frameName.FindStringSubmatch(filepath.Base(f)); m != nil && m[1] == pose {
				if b, err := readCapped(f); err == nil {
					out = append(out, b)
				}
			}
		}
		if len(out) > 0 {
			return out, true
		}
	}
	if b, err := embedded.ReadFile("web/sprites/" + pose + ".png"); err == nil {
		return [][]byte{b}, true
	}
	return nil, false
}

// Frame returns frame i of a pose (same rules as Frames).
func (s *Sprites) Frame(pose string, i int) ([]byte, bool) {
	fr, ok := s.Frames(pose)
	if !ok || i < 0 || i >= len(fr) {
		return nil, false
	}
	return fr[i], true
}

// EmbeddedBytes is the total size of the embedded example PNGs.
func EmbeddedBytes() int {
	n := 0
	es, _ := embedded.ReadDir("web/sprites")
	for _, e := range es {
		if i, err := e.Info(); err == nil {
			n += int(i.Size())
		}
	}
	return n
}
