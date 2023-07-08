package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type imageResponse struct {
	Images []Image
}

func (i *imageResponse) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return io.EOF
	}

	var err error

	switch b[0] {
	case '{':
		var img Image
		err = json.Unmarshal(b, &img)
		i.Images = append(i.Images, img)
	case '[':
		err = json.Unmarshal(b, &i.Images)
	default:
		return errors.New("invalid JSON input")
	}

	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
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
