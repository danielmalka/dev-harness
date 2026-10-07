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
	if skipped(err) {
		s.Warnings = append(s.Warnings, "poses.json not a regular file or over 1 MiB, using default map")
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

// MaxFrames caps the frames of one pose.
const MaxFrames = 60

// effective is the pose that supplies files: the pose itself when any file exists, else frente (jevmon behaviour).
// ok is false for an unknown name (e.g. "../x") or no file anywhere. Paths are names only, nothing is read here.
func (s *Sprites) effective(pose string) (paths []string, name string, ok bool) {
	if !KnownPose(pose) {
		return nil, "", false
	}
	for _, p := range []string{pose, "frente"} {
		paths = s.ownerPaths(p)
		if len(paths) > 0 || embeddedHas(p) {
			return paths, p, true
		}
	}
	return nil, "", false
}

func embeddedHas(pose string) bool {
	_, err := embedded.ReadFile("web/sprites/" + pose + ".png")
	return err == nil
}

// ownerPaths lists the owner's files for a pose: <pose>.png, else <pose>_00.png... (at most MaxFrames), by glob only.
func (s *Sprites) ownerPaths(pose string) []string {
	if s.dir == "" {
		return nil
	}
	single := filepath.Join(s.dir, pose+".png")
	if st, err := os.Lstat(single); err == nil && st.Mode().IsRegular() {
		return []string{single}
	}
	files, _ := filepath.Glob(filepath.Join(s.dir, pose+"_[0-9][0-9].png"))
	sort.Strings(files)
	var out []string
	for _, f := range files {
		if m := frameName.FindStringSubmatch(filepath.Base(f)); m != nil && m[1] == pose {
			out = append(out, f)
		}
	}
	if len(out) > MaxFrames {
		out = out[:MaxFrames]
	}
	return out
}

// FrameCount counts the frames of a pose without reading any file.
func (s *Sprites) FrameCount(pose string) int {
	paths, _, ok := s.effective(pose)
	switch {
	case !ok:
		return 0
	case len(paths) > 0:
		return len(paths)
	}
	return 1
}

// Frame returns frame i of a pose, reading only that file. An owner file that cannot be read (symlink, too big)
// falls back to the kit's example of the same pose for frame 0.
func (s *Sprites) Frame(pose string, i int) ([]byte, bool) {
	paths, name, ok := s.effective(pose)
	if !ok || i < 0 {
		return nil, false
	}
	if i < len(paths) {
		if b, err := readCapped(paths[i]); err == nil {
			return b, true
		}
	}
	if i == 0 {
		if b, err := embedded.ReadFile("web/sprites/" + name + ".png"); err == nil {
			return b, true
		}
	}
	return nil, false
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
