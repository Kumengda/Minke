package minke

type Monitor interface {
	Monitoring(containerFrame ContainerFrame) error
}

type BaseMonitor struct {
	monitorFunc func(containerFrame ContainerFrame) error
}

func (m *BaseMonitor) Monitoring(containerFrame ContainerFrame) error {
	return m.monitorFunc(containerFrame)
}

func NewBaseMonitor(monitorFunc func(containerFrame ContainerFrame) error) *BaseMonitor {
	return &BaseMonitor{
		monitorFunc: monitorFunc,
	}
}
