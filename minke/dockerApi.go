package minke

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io"
	"strings"
	"time"
)

func (m *Minke) collectContainer() {
	containers, _ := m.client.ContainerList(context.Background(), container.ListOptions{All: false})
	m.containers = nil
	for _, v := range containers {
		var containerNames []string
		for _, v1 := range v.Names {
			containerNames = append(containerNames, strings.TrimLeft(v1, "/"))
		}
		if m.matchContainerName(containerNames) {
			m.containers = append(m.containers, v)
			continue
		}
		if m.matchImageName(v.Image) {
			m.containers = append(m.containers, v)
			continue
		}
	}
}
func (m *Minke) snapshot() []ContainerFrame {
	var containerFrame []ContainerFrame
	lastContainer := m.containers
	m.collectContainer()
	for _, lc := range lastContainer {
		var mathFlag bool
		for _, c := range m.containers {
			if c.ID == lc.ID {
				containerFrame = append(containerFrame, ContainerFrame{
					ContainerInfo:   c,
					Status:          CONSTANT,
					Logs:            m.getContainerLog(c.ID),
					PerformanceData: PerformanceData{},
				})
				mathFlag = true
				break
			}
		}
		if !mathFlag {
			containerFrame = append(containerFrame, ContainerFrame{
				ContainerInfo:   lc,
				Status:          REDUCE,
				Logs:            m.getContainerLog(lc.ID),
				PerformanceData: PerformanceData{},
			})
		}

	}
	for _, c := range m.containers {
		var mathFlag bool
		for _, lc := range lastContainer {
			if lc.ID == c.ID {
				mathFlag = true
				break
			}
		}
		if !m.firstStartLock {
			containerFrame = append(containerFrame, ContainerFrame{
				ContainerInfo:   c,
				Status:          CONSTANT,
				Logs:            m.getContainerLog(c.ID),
				PerformanceData: PerformanceData{},
			})
			continue
		}
		if !mathFlag {
			containerFrame = append(containerFrame, ContainerFrame{
				ContainerInfo:   c,
				Status:          ADD,
				Logs:            m.getContainerLog(c.ID),
				PerformanceData: PerformanceData{},
			})
		}
	}
	m.firstStartLock = true
	m.snapShotTimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	return containerFrame
}

func (m *Minke) getContainerLog(containerID string) string {
	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Since:      m.snapShotTimeStamp,
		Tail:       m.logLine,
	}
	logs, err := m.client.ContainerLogs(context.Background(), containerID, logOptions)
	if err != nil {
		return ""
	}
	var logBuffer bytes.Buffer
	_, err = io.Copy(&logBuffer, logs)
	if err != nil {
		return ""
	}
	return logBuffer.String()
}
