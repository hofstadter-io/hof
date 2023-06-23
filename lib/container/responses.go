package container

import "time"

type Image struct {
	ID          string    `json:"Id"`
	ParentID    string    `json:"ParentId"`
	RepoTags    any       `json:"RepoTags"`
	RepoDigests []string  `json:"RepoDigests"`
	Size        int       `json:"Size"`
	SharedSize  int       `json:"SharedSize"`
	VirtualSize int       `json:"VirtualSize"`
	Labels      any       `json:"Labels"`
	Containers  int       `json:"Containers"`
	Names       []string  `json:"Names"`
	Digest      string    `json:"Digest"`
	History     []string  `json:"History"`
	Created     int       `json:"Created"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type Container struct {
	AutoRemove bool              `json:"AutoRemove"`
	Command    []string          `json:"Command"`
	CreatedAt  string            `json:"CreatedAt"`
	Exited     bool              `json:"Exited"`
	ExitedAt   int64             `json:"ExitedAt"`
	ExitCode   int               `json:"ExitCode"`
	ID         string            `json:"Id"`
	Image      string            `json:"Image"`
	ImageID    string            `json:"ImageID"`
	IsInfra    bool              `json:"IsInfra"`
	Labels     map[string]string `json:"Labels"`
	Mounts     []any             `json:"Mounts"`
	Names      []string          `json:"Names"`
	Namespaces any               `json:"Namespaces"`
	Networks   []any             `json:"Networks"`
	Pid        int               `json:"Pid"`
	Pod        string            `json:"Pod"`
	PodName    string            `json:"PodName"`
	Ports      any               `json:"Ports"`
	Size       any               `json:"Size"`
	StartedAt  int               `json:"StartedAt"`
	State      string            `json:"State"`
	Status     string            `json:"Status"`
	Created    int               `json:"Created"`
}
