package job

import "Thor/src/models"

type IJobExecutor interface {
	GetName() string
	Unlock(job *models.Job, args *map[string]string, options *map[string]string) error
	Start(job *models.Job) error
	Callback(body string) (any, error)
}
