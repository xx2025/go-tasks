package scheduler

import (
	"fmt"
	"go-tasks/master/model"
)

func LoadProcess(nodeId int) {
	processModel := model.NewProcessModel()
	list := processModel.GetAllProcess(nodeId)
	if len(list) > 0 {
		nodeModel := model.NewNodeModel()
		NodeMap := nodeModel.GetNodeWithHost()

		for _, v := range list {
			host, ok := NodeMap[v.NodeId]
			if !ok {
				continue
			}
			_ = ProcessStart(v.Id, v.MaxRetries, v.Name, host)
		}
	}

	//每5分钟检查进程
	spec := "*/5 * * * *"
	_, err := TaskSchedule.Cron.AddFunc(spec, ProcessCheck)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ProcessCheck() {
	processModel := model.NewProcessModel()
	list := processModel.GetAllProcess(0)
	if len(list) <= 0 {
		return
	}

	nodeModel := model.NewNodeModel()
	NodeMap := nodeModel.GetNodeWithHost()

	for _, v := range list {
		host, ok := NodeMap[v.NodeId]
		if !ok {
			continue
		}
		_ = ProcessStart(v.Id, v.MaxRetries, v.Name, host)
	}
}
