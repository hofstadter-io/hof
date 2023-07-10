package container

import (
	"regexp"
	"strconv"
	"strings"
)

type Image struct {
	Labels       map[string]string
	CreatedAt    string
	CreatedSince string
	Size         string
	ParentID     string
	Digest       string
	ID           string
	VirtualSize  string
	UniqueSize   string
	Containers   string
	RepoDigests  []string
	RepoTags     []string
	Names        []string
	History      []string
}

type Container struct {
	Labels     string
	CreatedAt  string
	Namespaces string
	Ports      string
	ID         string
	PodName    string
	Command    string
	Image      string
	ImageID    string
	State      string
	Pod        string
	Status     string
	Size       string
	Names      string
	Mounts     string
	Pid        int
	ExitedAt   int64
	ExitCode   int
	Created    int
	StartedAt  int
	IsInfra    bool
	Exited     bool
	AutoRemove bool
}

var portExp = regexp.MustCompile(`:(\d+)->`)

func (c Container) PortList() []int {
	var (
		parts = strings.Split(c.Ports, ",")
		ls    = make([]int, 0, len(parts))
	)

	for _, p := range parts {
		pp := portExp.FindStringSubmatch(p)
		if len(pp) != 2 {
			continue
		}

		i, err := strconv.Atoi(pp[1])
		if err != nil {
			continue
		}

		ls = append(ls, i)
	}

	return ls
}
