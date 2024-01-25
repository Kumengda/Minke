package minke

import "github.com/docker/docker/api/types"

type ContainerStatus int

const (
	ADD ContainerStatus = iota
	REDUCE
	CONSTANT
)

type PerformanceData struct {
	NetIo      string
	BlockIo    string
	CPU        string
	MEM        string
	LIMIT      string
	MEMPercent string
}

type ContainerFrame struct {
	ContainerInfo types.Container
	Status        ContainerStatus
	Logs          string
	PerformanceData
}
