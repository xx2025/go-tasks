package executor

import (
	"context"
	"fmt"
	"go-tasks/grpc/service"
	"go-tasks/internal/rpc"
	"go-tasks/utils"
)

func NodeRegister(name, host string) {
	f := func(ctx context.Context, client service.MasterClient) (interface{}, error) {
		return client.NodeRegister(ctx, &service.NodeRegisterRequest{
			Name: name,
			Host: host,
		})
	}
	res, err := rpc.ExecuteGRPCRequest(f)

	if err != nil {
		utils.Logger().Error("NodeRegister: " + err.Error())
	} else {
		fmt.Println("节点注册: ", res)
	}
}
