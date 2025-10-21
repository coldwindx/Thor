package initializer

import (
	"Thor/bootstrap"
	"Thor/bootstrap/inject"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func init() {
	initializer := &ViperInitializer{name: "ViperInitializer", order: 1}
	bootstrap.Beans.Provide(&inject.Object{Name: initializer.GetName(), Value: initializer, Completed: true})
}

type ViperInitializer struct {
	name  string
	order int
}

func (t *ViperInitializer) GetName() string {
	return t.name
}
func (t *ViperInitializer) GetOrder() int {
	return t.order
}
func (*ViperInitializer) Initialize() {
	// step1 设置配置文件路径
	config := os.Getenv("VIPER_CONFIG")
	if config == "" {
		panic(fmt.Errorf("load config from VIPER_CONFIG failed"))
	}
	// step2 初始化viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}
	// step3 加载配置信息
	if err := v.Unmarshal(&bootstrap.Config); err != nil {
		fmt.Println(err)
	}
	// step4 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed: ", in.Name)
		if err := v.Unmarshal(&bootstrap.Config); err != nil {
			fmt.Println(err)
		}
	})
}

func (*ViperInitializer) Close() {

}
