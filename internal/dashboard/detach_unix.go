//go:build !windows

package dashboard

import (
	"os/exec"
	"syscall"
)

func detachAttrs(c *exec.Cmd) { c.SysProcAttr = &syscall.SysProcAttr{Setsid: true} }
