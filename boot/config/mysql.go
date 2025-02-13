package config

import (
	_ "github.com/go-sql-driver/mysql"
	"go-tasks/boot/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitMysql() {
	mysqlConfig := mysql.Config{
		DriverName:                "mysql",
		DSN:                       getDSN(), // DSN data source name
		DefaultStringSize:         191,      // string 类型字段的默认长度
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		return
	} else {
		mysqlConf := &MasterConfigure.DB

		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConnections)
		sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConnections)
		sqlDB.SetConnMaxIdleTime(time.Second * 300)

		global.DB = db
	}
}

func getDSN() string {
	mysqlConf := &MasterConfigure.DB
	return mysqlConf.Username + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + mysqlConf.Port + ")/" + mysqlConf.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
