package main

import (
	"go-tasks/boot"
	"go-tasks/internal/scheduler"
)

func main() {
	boot.SetBootMod(boot.MasterMod)
	boot.Initialization()

	boot.StartHttp()

	scheduler.NewGrpcService()   //启动grpc服务
	scheduler.NewTaskScheduler() //启动定时任务
	go scheduler.LoadProcess(0)  //加载进程
	select {}

}
