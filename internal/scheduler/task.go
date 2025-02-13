package scheduler

import (
	"errors"
	"github.com/robfig/cron/v3"
	"go-tasks/grpc/service"
	"go-tasks/master/model"
	"go-tasks/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var TaskSchedule *TaskScheduler

type Task struct {
	Id       int
	Command  string
	Spec     string
	NodeHost string
}

func LoadTask() {
	taskModel := model.NewTaskModel()
	list := taskModel.GetListWithSchedule()
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
		_ = TaskSchedule.AddTask(v.Id, v.Spec, v.Name, host, int64(v.IsSingle))
	}
}

type TaskScheduler struct {
	Cron *cron.Cron
}

func NewTaskScheduler() {

	TaskSchedule = &TaskScheduler{
		Cron: cron.New(
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
				cron.SkipIfStillRunning(cron.DefaultLogger),
			),
		),
	}
	TaskSchedule.Cron.Start()
	LoadTask()

}

func (ts *TaskScheduler) UpdateTask(id int, spec, command, nodeHost string, isSingle int64) error {
	entryID, err := GetTaskCron(id)
	if err == nil {
		ts.Cron.Remove(entryID)
		DelTaskCron(id)
	}
	return ts.AddTask(id, spec, command, nodeHost, isSingle)
}

func (ts *TaskScheduler) AddTask(id int, spec, command, nodeHost string, isSingle int64) error {
	entryID, err := GetTaskCron(id)
	if err == nil {
		return errors.New("该任务已在cron调度中")
	}
	entryID, err = ts.Cron.AddFunc(spec, func() {
		_ = ts.ExecTask(id, command, nodeHost, isSingle)
	})
	if err != nil {
		utils.Logger().Error("Add Task failed on " + err.Error())
		return err
	}
	SetTaskCron(id, entryID)
	return nil
}

// 执行任务
func (ts *TaskScheduler) ExecTask(id int, command, nodeHost string, isSingle int64) error {

	taskLogsModel := model.NewTaskLogsModel()
	logId, err := taskLogsModel.NewLog(id)
	if logId <= 0 {
		utils.Logger().Error("任务调度失败：" + err.Error())
		return err
	}

	conn, err := GrpcClient(nodeHost)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	if err != nil {
		utils.Logger().Error("连接到node节点失败：" + err.Error())
		return err
	}
	client := service.NewNodeClient(conn)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err = client.TaskExec(ctx, &service.TaskExecRequest{
		TaskId:   int64(id),
		LogId:    int64(logId),
		Command:  command,
		IsSingle: isSingle,
	})

	return err
}

// 移除任务
func (ts *TaskScheduler) RemoveTask(id int) {
	entryID, err := GetTaskCron(id)
	if err == nil {
		ts.Cron.Remove(entryID)
		DelTaskCron(id)
	}
}

// 验证是否是合法的cron表达式
func IsCronFormat(spec string) bool {
	_, err := cron.ParseStandard(spec)
	return err == nil
}
