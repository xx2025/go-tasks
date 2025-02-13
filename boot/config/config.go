package config

import (
	"github.com/BurntSushi/toml"
	"go-tasks/boot/global"
	"log"
	"os"
)

var (
	MasterConfigure *MasterConfig
	NodeConfigure   *NodeConfig
)

type MasterConfig struct {
	Mode         string `toml:"mode"`
	Name         string `toml:"name"`
	Host         string `toml:"host"`
	GrpcPort     int    `toml:"grpc-port"`
	HttpPort     int    `toml:"http-port"`
	DB           Mysql  `toml:"db"`
	JwtSecretKey string `toml:"jwt-secret-key"`
}

func InitMasterConfig() {
	configFilePath := global.WorkerDir + "../master.config.toml"

	// 确保文件存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %v", err)
	}

	MasterConfigure = &MasterConfig{}
	if _, err := toml.DecodeFile(configFilePath, &MasterConfigure); err != nil {
		log.Fatalf("Failed to parse configuration file: %v", err)
	}
}

type NodeConfig struct {
	Name           string `toml:"name"`
	GrpcPort       int    `toml:"grpc-port"`
	Host           string `toml:"host"`
	MasterGrpcHost string `toml:"master-grpc-host"`
}

func InitNodeConfig() {
	configFilePath := global.WorkerDir + "/../node.config.toml"

	// 确保文件存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %v", err)
	}

	if _, err := toml.DecodeFile(configFilePath, &NodeConfigure); err != nil {
		log.Fatalf("Failed to parse configuration file: %v", err)
	}
}

type Mysql struct {
	Host               string `toml:"host"`
	Port               string `toml:"port"`
	DbName             string `toml:"db-name"`
	Username           string `toml:"username"`
	Password           string `toml:"password"`
	MaxIdleConnections int    `toml:"max-idle-connections"`
	MaxOpenConnections int    `toml:"max-open-connections"`
}
