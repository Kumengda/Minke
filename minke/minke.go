package minke

import (
	"fmt"
	. "github.com/Kumengda/Minke/runtime"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"regexp"
	"strconv"
	"time"
)

type Minke struct {
	client            *client.Client
	imageNameRule     []*regexp.Regexp
	containerNameRule []*regexp.Regexp
	containers        []types.Container
	monitorErrorFunc  func(err error)
	frameInterval     int
	firstStartLock    bool
	snapShotTimeStamp string
	logLine           string
	monitor           Monitor
}

func NewMinke(ops ...client.Opt) (*Minke, error) {
	Init()
	ops = append(ops, client.FromEnv)
	dockerClient, err := client.NewClientWithOpts(ops...)
	if err != nil {
		return nil, err
	}
	dockerClient.NegotiateAPIVersion(context.Background())
	return &Minke{
		client:        dockerClient,
		frameInterval: 1,
	}, nil
}

func (m *Minke) SetLogLine(line int) {
	m.logLine = strconv.Itoa(line)
}

func (m *Minke) SetContainerNameFilter(rules ...*regexp.Regexp) {
	m.containerNameRule = rules
}

func (m *Minke) SetMonitorErrorFunc(monitorErrorFunc func(err error)) {
	m.monitorErrorFunc = monitorErrorFunc
}

func (m *Minke) SetImageNameFilter(rules ...*regexp.Regexp) {
	m.imageNameRule = rules
}

func (m *Minke) SetLogMonitor(monitor Monitor) {
	m.monitor = monitor
}

func (m *Minke) Monitor(frameInterval int) {

	m.frameInterval = frameInterval
	for {
		snapShot := m.snapshot()
		for _, v := range snapShot {
			err := m.monitor.Monitoring(m.client, v)
			if err != nil {
				if m.monitorErrorFunc != nil {
					m.monitorErrorFunc(err)
				} else {
					MainInsp.Print(LEVEL_ERROR, Text(fmt.Sprintf("unhandled exception:%s", err)))
				}
			}
		}
		time.Sleep(time.Duration(m.frameInterval) * time.Second)
	}
}
