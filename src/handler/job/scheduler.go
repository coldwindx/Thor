package job

var JobScheduler = &jobScheduler{ExecutorManager: make(map[string]IJobExecutor)}

type jobScheduler struct {
	ExecutorManager map[string]IJobExecutor
}
