package minke

import "github.com/docker/docker/client"

type Monitor interface {
	Monitoring(client *client.Client, containerFrame ContainerFrame) error
}

type BaseMonitor struct {
	monitorFunc func(client *client.Client, containerFrame ContainerFrame) error
}

func (m *BaseMonitor) Monitoring(client *client.Client, containerFrame ContainerFrame) error {
	return m.monitorFunc(client, containerFrame)
}

func NewBaseMonitor(monitorFunc func(client *client.Client, containerFrame ContainerFrame) error) *BaseMonitor {
	return &BaseMonitor{
		monitorFunc: monitorFunc,
	}
}
