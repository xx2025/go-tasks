package executor

import (
	"context"
	"errors"
	"go-tasks/grpc/service"
	"go-tasks/internal/rpc"
	"go-tasks/utils"
	"os/exec"
	"runtime"
	"strings"
)

func taskWorking(taskId, logId int64, command string) {
	var cmd *exec.Cmd
	var err error

	defer func() {
		DelTaskCmd(taskId)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Process.Release()
		}
	}()

	if runtime.GOOS != "linux" {
		commandsArr := strings.Split(command, " ")
		exe := commandsArr[0]
		params := commandsArr[1:]
		cmd = exec.Command(exe, params...)

	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	if err = cmd.Start(); err != nil {
		//任务启动失败
		taskResultNotice(taskId, logId, -1, "任务启动失败: "+err.Error())
		return
	}
	SetTaskCmd(taskId, cmd)
	err = cmd.Wait()
	if err != nil {
		taskResultNotice(taskId, logId, -1, err.Error())
		return
	}

	//通知调度中心执行结果
	taskResultNotice(taskId, logId, 1, "执行成功")
}

func taskStop(taskId int64) error {
	cmd, err := GetTaskCmd(taskId)
	if err != nil {
		return errors.New("The task is not running")
	}
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Process.Release()
	}
	return nil

}

// 通知调度中心执行结果
func taskResultNotice(taskId, logId int64, status int64, message string) {

	f := func(ctx context.Context, client service.MasterClient) (interface{}, error) {
		return client.TaskExecRes(ctx, &service.TaskExecResRequest{
			TaskId:     taskId,
			LogId:      logId,
			TaskStatus: status,
			Message:    message,
		})
	}
	_, err := rpc.ExecuteGRPCRequest(f)

	if err != nil {
		utils.Logger().Error("taskResultNotice: " + err.Error())
	}
}
