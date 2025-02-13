package main

import (
	"go-tasks/boot"
	"go-tasks/internal/executor"
)

func main() {
	boot.SetBootMod(boot.NodeMod)
	boot.Initialization()

	executor.NewGrpcService()
	select {}
}
