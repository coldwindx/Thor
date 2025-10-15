package bootstrap

import (
	"Thor/config"
	"Thor/ctx"
	"Thor/utils/inject"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"
)

var Manager = make(map[string]Initializer)

func init() {
	ctx.Beans.Provide(&inject.Object{Name: "bootstrap.Factory", Value: &Factory{}, Private: true})
}

type Initializer interface {
	GetName() string
	GetOrder() int
	Initialize()
	Close()
}

type Factory struct {
	Initializers map[string]Initializer `inject:""`
}

func Initialize() {
	// step 初始化Bean
	if err := ctx.Beans.Populate(); err != nil {
		panic("初始化Bean失败: " + err.Error())
	}
	factory := ctx.Beans.GetByName("bootstrap.Factory").(*Factory)
	// step 初始化组件
	var instances []Initializer
	for _, instance := range factory.Initializers {
		instances = append(instances, instance)
	}
	sort.Slice(instances, func(l, r int) bool {
		return instances[l].GetOrder() < instances[r].GetOrder()
	})

	for _, ins := range instances {
		ins.Initialize()
	}
}

func Close() {
	factory := ctx.Beans.GetByName("bootstrap.Factory").(*Factory)
	// step 初始化组件
	for _, instance := range factory.Initializers {
		instance.Close()
	}
}

func Run() {
	r := ctx.Router
	srv := &http.Server{
		Addr:    config.Config.Application.Host + ":" + strconv.Itoa(config.Config.Application.Port),
		Handler: r,
	}
	go func() {
		err := srv.ListenAndServe()
		if nil != err && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// step 初始化Bean
	fmt.Println("init bean")
	err := ctx.Beans.Populate()
	if err != nil {
		panic("初始化Bean失败: " + err.Error())
	}
	// step 初始化路由
	fmt.Println("init route")
	for _, route := range ctx.Routes {
		route(ctx.Router)
	}
	// step 等待信号，优雅关闭
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(c); nil != err {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
