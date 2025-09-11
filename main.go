//go:generate statik -src=./resources
//go:generate go fmt statik/statik.go

package main

import (
	"Thor/bootstrap"
	"Thor/ctx"
	"Thor/src/handler/job"
	"Thor/src/models"
	"Thor/src/services"
	_ "Thor/statik"
	"Thor/utils"
	"fmt"
	"github.com/samber/lo"
	"time"
)

func main() {
	// step1 初始化配置
	bootstrap.Initialize()
	defer bootstrap.Close()
	ctx.Logger.Info("bootstrap init success!")
	// step 消息队列
	//go rabbitmq.Producer()
	//go rabbitmq.Consumer()
	// step 定时触发器
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func(t *time.Ticker) {
		var err error
		var jobService *services.JobServiceImpl
		if jobService, err = utils.GetBean[services.JobServiceImpl]("JobService"); err != nil {
			return
		}
		var scheduler *job.Scheduler
		if scheduler, err = utils.GetBean[job.Scheduler]("JobScheduler"); err != nil {
			return
		}

		for {
			<-t.C
			fmt.Println("Ticker:", time.Now().Format("2006-01-02 15:04:05"))
			query := models.JobQuery{Name: "task_input_job", PageSize: 1}
			jobs, err := jobService.Query(&query)
			if err != nil {
				_ = fmt.Errorf("query jobs error: " + err.Error())
				continue
			}
			first, b := lo.First(jobs)
			if !b {
				continue
			}
			if err = scheduler.Start(&first); err != nil {
				_ = fmt.Errorf("query jobs error: " + err.Error())
				continue
			}
		}
	}(ticker)
	// step 启动服务器
	bootstrap.Run()
}
