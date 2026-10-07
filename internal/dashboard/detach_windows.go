//go:build windows

package dashboard

import (
	"os/exec"
	"syscall"
)

const detachedProcess = 0x00000008

func detachAttrs(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{CreationFlags: detachedProcess | syscall.CREATE_NEW_PROCESS_GROUP}
}
