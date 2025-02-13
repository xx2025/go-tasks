package scheduler

import (
	"context"
	"errors"
	"fmt"
	"go-tasks/boot/config"
	"go-tasks/boot/global"
	"go-tasks/grpc/service"
	"go-tasks/master/model"
	"go-tasks/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type MasterServer struct {
	service.UnimplementedMasterServer
}

func NewGrpcService() {

	port := config.MasterConfigure.GrpcPort
	listen, err := net.Listen("tcp", ":"+utils.IntToStr(port))
	if err != nil {
		log.Fatal("rpc server start error: " + err.Error())
	}

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册服务
	service.RegisterMasterServer(server, &MasterServer{})

	fmt.Println("gRPC server is listening on port " + utils.IntToStr(port) + "...")

	go func() {
		// 启动服务器
		if err := server.Serve(listen); err != nil {
			log.Fatal("failed to serve")
		}
	}()
	time.Sleep(3 * time.Second)

}

func (MasterServer) TaskExecRes(ctx context.Context, req *service.TaskExecResRequest) (*service.TaskExecResResponse, error) {
	logId := int(req.GetLogId())
	taskLogsModel := model.NewTaskLogsModel()
	taskLog := taskLogsModel.FindById(logId)
	if taskLog == nil {
		return nil, errors.New("调度记录不存在")
	}
	taskLog.Status = int(req.GetTaskStatus())
	taskLog.Message = req.GetMessage()
	err := global.DB.WithContext(ctx).Save(taskLog).Error
	if err != nil {
		return nil, errors.New("更新结果失败")
	}
	return &service.TaskExecResResponse{
		Code: 200,
	}, nil
}

func (MasterServer) ProcessNotice(ctx context.Context, req *service.ProcessNoticeRequest) (*service.ProcessNoticeResponse, error) {
	processId := int(req.GetProcessId())
	message := req.GetMessage()
	processLogsModel := model.NewProcessLogsModel()
	_, _ = processLogsModel.NewLog(processId, message)

	return &service.ProcessNoticeResponse{
		Code:    200,
		Message: "ok",
	}, nil
}

func (MasterServer) NodeRegister(ctx context.Context, req *service.NodeRegisterRequest) (*service.NodeRegisterResponse, error) {
	name := req.GetName()
	host := req.GetHost()
	nodeModel := model.NewNodeModel()
	node, err := nodeModel.InsertIfNotExsit(ctx, name, host)

	if err == nil && node.Id > 0 && node.Status == 1 {
		go LoadProcess(0)
	}
	return &service.NodeRegisterResponse{
		Code:    200,
		Message: "ok",
	}, nil
}
