package minke

import "github.com/docker/docker/api/types"

type ContainerStatus int

const (
	ADD ContainerStatus = iota
	REDUCE
	CONSTANT
)

type PerformanceData struct {
	CPU   string
	MEM   float64
	LIMIT float64
}

type ContainerFrame struct {
	ContainerInfo types.Container
	Status        ContainerStatus
	Logs          string
	PerformanceData
}
