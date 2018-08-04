package cgroups

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MEMORYMP = "memory"
const CPUSETMP = "cpuset"
const CPUSHAREMP = "cpu,cpuacct"

const MPPATH = "/proc/self/mountinfo"

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

func FindCgroupMountPoint(subsystem string) string {
	if MEMORYMP != subsystem && CPUSETMP != subsystem && CPUSHAREMP != subsystem {
		log.Warning("not support MP")
		return ""
	}

	f, err := os.Open(MPPATH)
	if err != nil {
		log.Error(err)
		return ""
	}

	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if strings.Contains(line, subsystem) {
			return strings.Split(line, " ")[4]
		}
		if nil != err {
			if err != io.EOF {
				log.Error(err)
			}
			return ""
		}
	}
	return ""
}

/**
 * limit cgroup subsystem resources, as cpuset、cpushare、memory
 *
 * @Param cType: 0 - memory  1 - cpuSet  2 - cpuShare
 * @Param limitValue: the value to be writtened to cgroup subsystem
 * @return
 */
func LimitCgroupResources(cType int, limitValue string, pid int) bool {
	if cType < 0 || cType > 2 {
		log.Error("cgroup type not supported")
		return false
	}

	var mp string
	var regStr string
	var mpType string
	var err error

	if "" == limitValue {
		log.Error("get limit value failed!")
		return false
	} else if 0 == cType {
		regStr = `^\d{1,}m`
		mpType = MEMORYMP
	} else if 1 == cType {
		regStr = `^\d{1,}`
		mpType = CPUSETMP
	} else if 2 == cType {
		regStr = `^\d{1,}`
		mpType = CPUSHAREMP
	}

	if r := regexp.MustCompile(regStr); !r.MatchString(limitValue) {
		log.Error("cg param format wrong")
		return false
	}

	if mp = FindCgroupMountPoint(mpType); "" == mp {
		log.Errorf("mount point %s not found", mpType)
		return false
	}

	os.Mkdir(mp+"/pDocker", 0766)

	var f *os.File

	if 0 == cType {
		f, err = os.OpenFile(mp+"/pDocker/"+"memory.limit_in_bytes", os.O_WRONLY, 0666)
		if nil != err {
			log.Error(err)
			return false
		}
	} else if 1 == cType {
		f, err = os.OpenFile(mp+"/pDocker/"+"cpuset.cpus", os.O_WRONLY, 0666)
		if nil != err {
			log.Error(err)
			return false
		}
	} else if 2 == cType {
		f, err = os.OpenFile(mp+"/pDocker/"+"cpu.shares", os.O_WRONLY, 0666)
		if nil != err {
			log.Error(err)
			return false
		}
	}

	defer f.Close()

	n, err1 := io.WriteString(f, limitValue)
	if nil != err1 {
		log.Error(err1)
		log.Errorf("limitValue:%s, n:%d", limitValue, n)
		return false
	}

	log.Infof("path:%s, limitValue:%s", mp+"/pDocker/", limitValue)

	f1, err2 := os.OpenFile(mp+"/pDocker/"+"tasks", os.O_WRONLY, 0666)
	if nil != err2 {
		log.Error(err)
		return false
	}
	n, err1 = io.WriteString(f1, strconv.Itoa(pid))
	if nil != err1 {
		log.Error(err1)
		log.Error(n)
		return false
	}

	return true
}
