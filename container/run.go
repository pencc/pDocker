package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"pDocker/cgroups"
	"syscall"
)

func Run(tty bool, command string, config cgroups.ResourceConfig) {
	syscall.Mount("proc", "/proc", "proc", 0, "")
	parent := NewParentProcess(tty, command)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	if "" != config.MemoryLimit {
		cgroups.LimitCgroupResources(0, config.MemoryLimit, parent.Process.Pid)
	}

	if "" != config.CpuShare {
		cgroups.LimitCgroupResources(2, config.CpuShare, parent.Process.Pid)
	}

	if "" != config.CpuSet {
		cgroups.LimitCgroupResources(1, config.CpuSet, parent.Process.Pid)
	}

	parent.Wait()
	os.Exit(1)
}
