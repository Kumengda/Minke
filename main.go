package main

import (
	"fmt"
	"github.com/Kumengda/Minke/minke"
	"regexp"
)

func main() {
	myMink, err := minke.NewMinke()
	if err != nil {
		return
	}

	regex1, _ := regexp.Compile(".*ai.radar.*")
	myMink.SetImageNameFilter(regex1)
	myMink.SetLogLine(1)
	myMink.SetLogMonitor(minke.NewBaseMonitor(func(containerFrame minke.ContainerFrame) error {
		var status string
		switch containerFrame.Status {
		case minke.ADD:
			status = "add"
		case minke.CONSTANT:
			status = "constant"
		case minke.REDUCE:
			status = "reduce"
		}
		fmt.Println("----------------------------------------------------------------------------------")
		fmt.Printf("name:%s id:%s status:%s Logs:%s\n", containerFrame.ContainerInfo.Names[0], containerFrame.ContainerInfo.ID, status, containerFrame.Logs)
		return nil
	}))
	myMink.Monitor(1)
}