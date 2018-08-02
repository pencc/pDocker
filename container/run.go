package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

func Run(tty bool, command string) {
	syscall.Mount("proc", "/proc", "proc", 0, "")
	parent := NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(1)
}
