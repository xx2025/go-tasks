package rpc

import (
	"context"
	"go-tasks/boot/config"
	"go-tasks/grpc/service"
	"go-tasks/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// 调度中心 grpc请求封装
func ScheduleGRPCRequest(
	clientHost string,
	method func(ctx context.Context, client service.NodeClient) (interface{}, error),
) (interface{}, error) {
	pool, err := GetPool(clientHost)
	if err != nil {
		return nil, err
	}
	client, err := pool.get()

	if err != nil {
		return nil, err
	}

	defer pool.put(client)

	conn := service.NewNodeClient(client.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return method(ctx, conn)
}

// 执行器 grpc请求封装
func ExecuteGRPCRequest(
	method func(ctx context.Context, client service.MasterClient) (interface{}, error),
) (interface{}, error) {

	client, err := grpc.NewClient(config.NodeConfigure.MasterGrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		utils.Logger().Error("did not connect: " + err.Error())
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(client)
	conn := service.NewMasterClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return method(ctx, conn)

}
