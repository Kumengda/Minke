package minke

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
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
					PerformanceData: m.getContainerPerformance(c.ID),
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
				PerformanceData: m.getContainerPerformance(lc.ID),
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
				PerformanceData: m.getContainerPerformance(c.ID),
			})
			continue
		}
		if !mathFlag {
			containerFrame = append(containerFrame, ContainerFrame{
				ContainerInfo:   c,
				Status:          ADD,
				Logs:            m.getContainerLog(c.ID),
				PerformanceData: m.getContainerPerformance(c.ID),
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

func (m *Minke) getContainerPerformance(containerID string) PerformanceData {
	var performanceData PerformanceData
	stats, err := m.client.ContainerStats(context.Background(), containerID, false)
	if err != nil {
		return performanceData
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stats.Body)
	if err != nil {
		return performanceData
	}
	newStr := buf.String()
	cpuDelta := gjson.Get(newStr, "cpu_stats").Get("cpu_usage").Get("total_usage").Float() - gjson.Get(newStr, "precpu_stats").Get("cpu_usage").Get("total_usage").Float()
	systemDelta := gjson.Get(newStr, "cpu_stats").Get("system_cpu_usage").Float() - gjson.Get(newStr, "precpu_stats").Get("system_cpu_usage").Float()
	cpuUseage := cpuDelta / systemDelta * 100 * float64(len(gjson.Get(newStr, "cpu_stats").Get("cpu_usage").Get("percpu_usage").Array()))
	performanceData.CPU = strconv.FormatFloat(cpuUseage, 'f', 2, 64) + "%"
	memoryUsage := gjson.Get(newStr, "memory_stats").Get("usage").Float() / 1024 / 1024
	limit := gjson.Get(newStr, "memory_stats").Get("limit").Float() / 1024 / 1024
	performanceData.MEM = memoryUsage
	performanceData.LIMIT = limit
	return performanceData
}
