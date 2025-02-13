package executor

import (
	"context"
	"fmt"
	"go-tasks/boot/config"
	"go-tasks/grpc/service"
	"go-tasks/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	"time"
)

type NodeServer struct {
	service.UnimplementedNodeServer
}

func NewGrpcService() {

	port := config.NodeConfigure.GrpcPort

	listen, err := net.Listen("tcp", ":"+utils.IntToStr(port))
	if err != nil {
		log.Fatal("rpc server start error", err)
	}

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务
	service.RegisterNodeServer(server, &NodeServer{})

	fmt.Println("gRPC server is listening on port " + utils.IntToStr(port) + "...")

	// 启动服务器
	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatal("failed to serve")
		}
	}()
	time.Sleep(3 * time.Second)
	host := strings.TrimSpace(config.NodeConfigure.Host)
	if len(host) > 0 {
		NodeRegister(config.NodeConfigure.Name, host)
	}

}

func (n *NodeServer) Ping(ctx context.Context, req *service.PingRequest) (*service.PingResponse, error) {
	return &service.PingResponse{
		Code:    200,
		Message: "pong",
	}, nil
}

func (n *NodeServer) TaskExec(ctx context.Context, req *service.TaskExecRequest) (*service.TaskExecResponse, error) {
	command := req.GetCommand()
	taskId := req.GetTaskId()
	logId := req.GetLogId()
	isSingle := req.GetIsSingle()

	_, err := GetTaskCmd(taskId)
	if err == nil && isSingle == 1 {
		return &service.TaskExecResponse{
			Code:    400,
			Message: "任务正在运行",
		}, nil
	}

	go taskWorking(taskId, logId, command)
	return &service.TaskExecResponse{
		Code:    200,
		Message: "调度成功",
	}, nil
}

func (n *NodeServer) TaskExecPID(ctx context.Context, req *service.TaskExecPIDRequest) (*service.TaskExecPIDResponse, error) {
	var pid int
	taskId := req.GetTaskId()
	cmd, err := GetTaskCmd(taskId)
	if err == nil {
		pid = cmd.Process.Pid
	}

	return &service.TaskExecPIDResponse{
		Code:    200,
		Message: "ok",
		Data: &service.TaskExecPIDResponseData{
			Pid: int64(pid), // 设置 pid
		},
	}, nil
}

func (n *NodeServer) TaskExecStop(ctx context.Context, req *service.TaskExecStopRequest) (*service.TaskExecStopResponse, error) {

	taskId := req.GetTaskId()
	err := taskStop(taskId)
	if err != nil {
		return &service.TaskExecStopResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}
	return &service.TaskExecStopResponse{
		Code:    200,
		Message: "ok",
	}, nil
}

func (n *NodeServer) ProcessStart(ctx context.Context, req *service.ProcessStartRequest) (*service.ProcessStartResponse, error) {
	processId := req.GetProcessId()
	command := req.GetProcessCommand()
	maxRetries := int(req.GetMaxRetries())

	commandsArr := strings.Split(command, " ")
	exe := commandsArr[0]
	params := commandsArr[1:]

	process, err := NewProcess(processId, exe, params, maxRetries)

	if err != nil {
		return &service.ProcessStartResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	go process.Start()
	return &service.ProcessStartResponse{
		Code:    200,
		Message: "ok",
	}, nil
}

func (n *NodeServer) ProcessStop(ctx context.Context, req *service.ProcessStopRequest) (*service.ProcessStopResponse, error) {
	processId := req.GetProcessId()
	process, _ := GetProcessMap(int64(processId))
	if process == nil {
		return &service.ProcessStopResponse{
			Code:    400,
			Message: "该进程没有运行",
		}, nil
	}
	go process.Stop()
	return &service.ProcessStopResponse{
		Code:    200,
		Message: "ok",
	}, nil
}

func (n *NodeServer) ProcessInfo(ctx context.Context, req *service.ProcessInfoRequest) (*service.ProcessInfoResponse, error) {
	processId := req.GetProcessId()

	process, _ := GetProcessMap(int64(processId))
	if process == nil {
		return &service.ProcessInfoResponse{
			Code:    400,
			Message: "该进程没有运行",
		}, nil
	}
	data := &service.ProcessInfoResponseData{
		Pid: int64(process.cmd.Process.Pid),
	}
	return &service.ProcessInfoResponse{
		Code:    200,
		Message: "ok",
		Data:    data,
	}, nil
}
