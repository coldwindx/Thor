package test

import (
	"Thor/bootstrap"
	"Thor/src/mapper"
	"github.com/stretchr/testify/assert"
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
	db := bootstrap.Beans.GetByName("MysqlConnection").(*gorm.DB)
	session, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	_ = session.Ping()
}

func TestMapperProxy(t *testing.T) {
	// 指定环境变量
	_ = os.Setenv("VIPER_CONFIG", "../resources/application.yaml")
	bootstrap.Initialize()
	defer bootstrap.Close()
	bootstrap.Beans.Populate()

	jobMapper := bootstrap.Beans.GetByName("JobMapper").(*mapper.JobMapper)
	assert.Equal(t, "JobMapper.Test()", jobMapper.Test())
}
