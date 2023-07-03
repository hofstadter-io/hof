package container

import (
	"strconv"
	"strings"
	"time"
)

type Image struct {
	CreatedAt   time.Time
	Labels      map[string]string
	ParentID    string
	Digest      string
	ID          string
	Names       []string
	RepoDigests []string
	History     []string
	RepoTags    []string
	VirtualSize int
	SharedSize  int
	Containers  int
	Size        int
	Created     int
}

type Container struct {
	Size       string
	Ports      string
	Namespaces string
	Labels     map[string]string
	ID         string
	PodName    string
	CreatedAt  string
	Image      string
	ImageID    string
	State      string
	Pod        string
	Status     string
	Names      []string
	Mounts     []any
	Networks   []any
	Command    []string
	Pid        int
	ExitedAt   int64
	ExitCode   int
	Created    int
	StartedAt  int
	IsInfra    bool
	Exited     bool
	AutoRemove bool
}

func (c Container) PortList() []int {
	var (
		parts = strings.Split(c.Ports, ",")
		ls    = make([]int, 0, len(parts))
	)

	for _, p := range parts {
		pp := strings.Split(p, "/")
		if len(pp) != 2 {
			continue
		}

		i, err := strconv.Atoi(pp[0])
		if err != nil {
			continue
		}

		ls = append(ls, i)
	}

	return ls
}
