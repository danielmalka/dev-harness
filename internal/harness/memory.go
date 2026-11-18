package harness

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// MemoryFile is the record the hash guard protects since PRD-034; FactsFile joined it in PRD-035 (0.47.0). EPOCHAL.md and RISKS.md stay outside.
const MemoryFile = "MEMORY.md"

// FactsFile is the system facts record of a consumer project (PRD-035), next to REPO-MAP.md in the harness dir.
const FactsFile = "FACTS.md"

// RecordFiles is the closed list of files `dh memory hash` and `dh memory write` accept in --file (PRD-035 R11). Any other name is refused.
var RecordFiles = []string{MemoryFile, FactsFile}

// CheckRecordFile refuses a --file value that is not on the closed list, naming both valid files.
func CheckRecordFile(name string) error {
	for _, f := range RecordFiles {
		if name == f {
			return nil
		}
	}
	return fmt.Errorf("--file aceita só %s ou %s, recebi %q", RecordFiles[0], RecordFiles[1], name)
}

// MemoryAbsent is the digest value of a MEMORY.md that does not exist yet.
const MemoryAbsent = "absent"

var shaHex = regexp.MustCompile(`^[0-9a-f]{64}$`)

// MemoryChangedError is returned by WriteMemory and WriteFile when the file is no longer the one the caller read. File is empty for MEMORY.md.
type MemoryChangedError struct{ Expected, Found, File string }

func (e *MemoryChangedError) Error() string {
	name := e.File
	if name == "" {
		name = MemoryFile
	}
	return fmt.Sprintf("%s mudou desde que foi lido (esperado %s, achei %s): alguém escreveu no meio. "+
		"Nada foi gravado. Leia o %s de novo e refaça a edição", name, e.Expected, e.Found, name)
}

// IsMemoryChanged reports whether err is a MemoryChangedError.
func IsMemoryChanged(err error) bool {
	var c *MemoryChangedError
	return errors.As(err, &c)
}

// MemoryHash returns the sha256 of <dir>/MEMORY.md as lower-case hex, or MemoryAbsent when there is no file.
func MemoryHash(dir string) (string, error) { return FileHash(dir, MemoryFile) }

// FileHash is MemoryHash for any file of the closed list RecordFiles.
func FileHash(dir, name string) (string, error) {
	if err := CheckRecordFile(name); err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(dir, name))
	if errors.Is(err, os.ErrNotExist) {
		return MemoryAbsent, nil
	}
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// WriteMemory is a compare-and-swap on <dir>/MEMORY.md; see WriteFile.
func WriteMemory(dir, expect string, content []byte) (string, error) {
	return WriteFile(dir, MemoryFile, expect, content)
}

// WriteFile is a compare-and-swap on <dir>/<name>, where name is on the closed list RecordFiles. Under the lock <name>.lock it hashes the current file and
// compares with expect (a sha256 hex or MemoryAbsent); only when equal it replaces the file (temp file in the same folder, flushed, renamed) and
// returns the new digest. A different digest writes nothing and returns *MemoryChangedError. Empty content is refused.
// Residual (PRD-034 R7): a writer that bypasses this function and lands between the compare and the rename is lost.
func WriteFile(dir, name, expect string, content []byte) (string, error) {
	if err := CheckRecordFile(name); err != nil {
		return "", err
	}
	if expect != MemoryAbsent && !shaHex.MatchString(expect) {
		return "", fmt.Errorf("--expect deve ser um sha256 (64 dígitos hexadecimais minúsculos) ou %q, recebi %q", MemoryAbsent, expect)
	}
	if len(content) == 0 {
		return "", fmt.Errorf("recusei gravar um %s vazio", name)
	}
	release, err := lockNamed(dir, name+".lock")
	if err != nil {
		return "", err
	}
	defer release()
	found, err := FileHash(dir, name)
	if err != nil {
		return "", err
	}
	if found != expect {
		file := ""
		if name != MemoryFile {
			file = name
		}
		return "", &MemoryChangedError{Expected: expect, Found: found, File: file}
	}
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	tmp := filepath.Join(dir, "."+name+".tmp-"+hex.EncodeToString(nonce))
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
	if err != nil {
		return "", err
	}
	_, werr := f.Write(content)
	if werr == nil {
		werr = f.Sync()
	}
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		os.Remove(tmp)
		return "", werr
	}
	if err := os.Rename(tmp, filepath.Join(dir, name)); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return hashBytes(content), nil
}
