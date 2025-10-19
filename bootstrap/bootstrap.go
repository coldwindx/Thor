package bootstrap

import (
	"Thor/utils/inject"
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
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
	Beans.Provide(&inject.Object{Name: "bootstrap.Factory", Value: &Factory{}, Private: true})
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
	// step 按顺序初始化基础组件
	objs := Beans.GetByType(new(Initializer))
	instances := lo.Map(objs, func(obj any, _ int) Initializer { return obj.(Initializer) })
	sort.Slice(instances, func(l, r int) bool { return instances[l].GetOrder() < instances[r].GetOrder() })
	lo.ForEach(instances, func(ins Initializer, _ int) { ins.Initialize() })

	// step 初始化容器管理
	Beans.Populate()
}

func Close() {
	factory := Beans.GetByName("bootstrap.Factory").(*Factory)
	// step 初始化组件
	for _, instance := range factory.Initializers {
		instance.Close()
	}
}

func Run() {
	r := Router
	srv := &http.Server{
		Addr:    Config.Application.Host + ":" + strconv.Itoa(Config.Application.Port),
		Handler: r,
	}
	go func() {
		err := srv.ListenAndServe()
		if nil != err && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// step 初始化路由
	fmt.Println("init route")
	for _, route := range Routes {
		route(Router)
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
