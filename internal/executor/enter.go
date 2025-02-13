package executor

import (
	"errors"
	"os/exec"
	"sync"
)

var TaskCmd sync.Map

func SetTaskCmd(taskId int64, cmd *exec.Cmd) {
	TaskCmd.Store(taskId, cmd)
}

func GetTaskCmd(taskId int64) (*exec.Cmd, error) {

	if value, ok := TaskCmd.Load(taskId); ok {
		return value.(*exec.Cmd), nil
	}

	return nil, errors.New("task not running")
}

func DelTaskCmd(taskId int64) {
	TaskCmd.Delete(taskId)
}
