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

// MemoryFile is the record the guard protects (PRD-034): only MEMORY.md, not EPOCHAL.md or RISKS.md.
const MemoryFile = "MEMORY.md"

// MemoryAbsent is the digest value of a MEMORY.md that does not exist yet.
const MemoryAbsent = "absent"

var shaHex = regexp.MustCompile(`^[0-9a-f]{64}$`)

// MemoryChangedError is returned by WriteMemory when the file is no longer the one the caller read.
type MemoryChangedError struct{ Expected, Found string }

func (e *MemoryChangedError) Error() string {
	return fmt.Sprintf("%s mudou desde que foi lido (esperado %s, achei %s): alguém escreveu no meio. "+
		"Nada foi gravado. Leia o %s de novo e refaça a edição", MemoryFile, e.Expected, e.Found, MemoryFile)
}

// IsMemoryChanged reports whether err is a MemoryChangedError.
func IsMemoryChanged(err error) bool {
	var c *MemoryChangedError
	return errors.As(err, &c)
}

// MemoryHash returns the sha256 of <dir>/MEMORY.md as lower-case hex, or MemoryAbsent when there is no file.
func MemoryHash(dir string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, MemoryFile))
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

// WriteMemory is a compare-and-swap on <dir>/MEMORY.md. Under the lock MEMORY.md.lock it hashes the current file and compares with
// expect (a sha256 hex or MemoryAbsent); only when equal it replaces the file (temp file in the same folder, flushed, renamed) and
// returns the new digest. A different digest writes nothing and returns *MemoryChangedError. Empty content is refused.
// Residual (PRD-034 R7): a writer that bypasses this function and lands between the compare and the rename is lost.
func WriteMemory(dir, expect string, content []byte) (string, error) {
	if expect != MemoryAbsent && !shaHex.MatchString(expect) {
		return "", fmt.Errorf("--expect deve ser um sha256 (64 dígitos hexadecimais minúsculos) ou %q, recebi %q", MemoryAbsent, expect)
	}
	if len(content) == 0 {
		return "", fmt.Errorf("recusei gravar um %s vazio", MemoryFile)
	}
	release, err := lockNamed(dir, MemoryFile+".lock")
	if err != nil {
		return "", err
	}
	defer release()
	found, err := MemoryHash(dir)
	if err != nil {
		return "", err
	}
	if found != expect {
		return "", &MemoryChangedError{Expected: expect, Found: found}
	}
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	tmp := filepath.Join(dir, "."+MemoryFile+".tmp-"+hex.EncodeToString(nonce))
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
	if err := os.Rename(tmp, filepath.Join(dir, MemoryFile)); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return hashBytes(content), nil
}
