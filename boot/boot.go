package boot

import (
	"go-tasks/boot/config"
	"go-tasks/boot/global"
	"go-tasks/master/router"
	"go-tasks/utils"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	MasterMod = "master"
	NodeMod   = "node"
)

func SetBootMod(mod string) {
	global.Mod = mod
}

func setResourceDir(dir string) {
	global.ResourceDir = dir
}

func Initialization() {
	//初始化工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Exception in obtaining working directory:  %v", err)
	}
	global.WorkerDir = cwd + "/"

	//初始化日志目录
	global.LogDir = global.WorkerDir + "/../log/"

	//初始化 config
	if global.Mod == MasterMod {
		config.InitMasterConfig()
		config.InitZapLogger()
		config.InitMysql()
		setResourceDir(global.WorkerDir + "../resource/")
	} else {
		config.InitNodeConfig()
		config.InitZapLogger()
	}
}

func StartHttp() {
	go func() {
		Router := router.InitRouter()
		port := config.MasterConfigure.HttpPort
		s := &http.Server{
			Addr:           ":" + utils.IntToStr(port),
			Handler:        Router,
			ReadTimeout:    10 * time.Minute,
			WriteTimeout:   10 * time.Minute,
			MaxHeaderBytes: 1 << 20,
		}
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal("http service start error: " + err.Error())
		}
	}()

}
