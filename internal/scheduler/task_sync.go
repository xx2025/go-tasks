package scheduler

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sync"
)

var TaskCronSync sync.Map

func SetTaskCron(taskId int, entryID cron.EntryID) {
	TaskCronSync.Store(taskId, entryID)
}

func GetTaskCron(taskId int) (cron.EntryID, error) {

	if value, ok := TaskCronSync.Load(taskId); ok {
		return value.(cron.EntryID), nil
	}

	return 0, errors.New("taskCron not exist")
}

func DelTaskCron(taskId int) {
	TaskCronSync.Delete(taskId)
}
