package scheduler

import (
	"context"
	"errors"
	"go-tasks/grpc/service"
	"go-tasks/internal/rpc"
	"go-tasks/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcClient(clientHost string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(clientHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.Logger().Error("Unable to connect to grpc node：" + clientHost)
		return conn, errors.New("Unable to connect to grpc node")
	}
	return conn, nil
}

func NodePing(clientHost string) error {
	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.Ping(ctx, &service.PingRequest{Message: "ping"})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return err
	}
	res2 := res.(*service.PingResponse)
	if res2.GetCode() != 200 {
		return errors.New("节点异常")
	}
	return nil
}

func GetTaskPID(taskId int, clientHost string) (int64, error) {

	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.TaskExecPID(ctx, &service.TaskExecPIDRequest{
			TaskId: int64(taskId),
		})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return 0, err
	}
	nodeRes := res.(*service.TaskExecPIDResponse)
	return nodeRes.GetData().Pid, nil
}

func TaskStop(taskId int, clientHost string) error {

	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.TaskExecStop(ctx, &service.TaskExecStopRequest{
			TaskId: int64(taskId),
		})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return err
	}
	nodeRes := res.(*service.TaskExecStopResponse)
	if nodeRes.GetCode() == 200 {
		return nil
	}
	return errors.New(nodeRes.GetMessage())
}

func ProcessStart(processId int, retryCount int, command string, clientHost string) error {
	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.ProcessStart(ctx, &service.ProcessStartRequest{
			ProcessId:      int64(processId),
			ProcessCommand: command,
			MaxRetries:     uint64(retryCount),
		})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return err
	}
	nodeRes := res.(*service.ProcessStartResponse)
	if nodeRes.GetCode() == 200 {
		return nil
	}
	return errors.New(nodeRes.GetMessage())
}

func ProcessStop(processId int, clientHost string) error {
	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.ProcessStop(ctx, &service.ProcessStopRequest{
			ProcessId: int64(processId),
		})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return err
	}
	nodeRes := res.(*service.ProcessStopResponse)
	if nodeRes.GetCode() == 200 {
		return nil
	}
	return errors.New(nodeRes.GetMessage())
}

func GetProcessPID(processId int, clientHost string) (int64, error) {

	f := func(ctx context.Context, client service.NodeClient) (interface{}, error) {
		return client.ProcessInfo(ctx, &service.ProcessInfoRequest{
			ProcessId: int64(processId),
		})
	}
	res, err := rpc.ScheduleGRPCRequest(clientHost, f)
	if err != nil {
		utils.Logger().Error(err.Error())
		return 0, err
	}
	nodeRes := res.(*service.ProcessInfoResponse)
	if nodeRes.GetCode() == 200 {
		return nodeRes.GetData().Pid, nil
	}
	return 0, errors.New(nodeRes.GetMessage())
}
