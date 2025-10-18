package mapper

import (
    "Thor/bootstrap"
    "Thor/src/mapper"
    "Thor/src/models/do"
    "context"
    "errors"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestMapperProxy(t *testing.T) {
    // 指定环境变量
    _ = os.Setenv("VIPER_CONFIG", "../../resources/application.yaml")
    bootstrap.Initialize()
    defer bootstrap.Close()
    bootstrap.Beans.Populate()

    jobMapper := bootstrap.Beans.GetByName("JobMapper").(*mapper.JobMapper)
    assert.Equal(t, "JobMapper.Test()", jobMapper.Test(context.Background()))
}

func TestNoTransaction(t *testing.T) {
    // 指定环境变量
    _ = os.Setenv("VIPER_CONFIG", "../../resources/application.yaml")
    bootstrap.Initialize()
    defer bootstrap.Close()
    bootstrap.Beans.Populate()

    flag := false

    client := bootstrap.Beans.GetByName("DBClient").(*bootstrap.DBClient)
    jobMapper := bootstrap.Beans.GetByName("JobMapper").(*mapper.JobMapper)

    err := client.Transaction(context.Background(), func(ctx context.Context) error {
        job1 := &do.Job{Name: "test"}
        _, _ = jobMapper.Insert(ctx, job1)
        // 触发error
        if flag {
            return errors.New("transactional error to rollback")
        }
        job2 := &do.Job{Name: "test"}
        _, _ = jobMapper.Insert(ctx, job2)
        return nil
    })
    assert.Nil(t, err)
}

func TestTransaction(t *testing.T) {
    // 指定环境变量
    _ = os.Setenv("VIPER_CONFIG", "../../resources/application.yaml")
    bootstrap.Initialize()
    defer bootstrap.Close()
    bootstrap.Beans.Populate()

    flag := true

    client := bootstrap.Beans.GetByName("DBClient").(*bootstrap.DBClient)
    jobMapper := bootstrap.Beans.GetByName("JobMapper").(*mapper.JobMapper)

    err := client.Transaction(context.Background(), func(ctx context.Context) error {
        job1 := &do.Job{Name: "test"}
        _, _ = jobMapper.Insert(ctx, job1)
        // 触发error
        if flag {
            return errors.New("transactional error to rollback")
        }
        job2 := &do.Job{Name: "test"}
        _, _ = jobMapper.Insert(ctx, job2)
        return nil
    })
    assert.Equal(t, "transactional error to rollback", err.Error())
}
