package test

import (
	"Thor/bootstrap"
	"Thor/ctx"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestMysqlConnection(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()

	// 测试数据库连接是否成功
	db := ctx.Beans.GetByName("MysqlConnection").(*gorm.DB)
	session, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	_ = session.Ping()
}
