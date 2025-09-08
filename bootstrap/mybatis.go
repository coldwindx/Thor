package bootstrap

import (
	"Thor/config"
	"Thor/ctx"
	_ "Thor/src/mapper"
	"fmt"
	"github.com/rakyll/statik/fs"
	"github.com/zhuxiujia/GoMybatis"
	"io"
	"strconv"
)

func init() {
	v := &MybatisInitializer{name: "mybatis", order: 10}
	Manager[v.name] = v
}

type MybatisInitializer struct {
	name  string
	order int
}

func (ts *MybatisInitializer) GetName() string {
	return ts.name
}

func (ts *MybatisInitializer) GetOrder() int {
	return ts.order
}

func (ts *MybatisInitializer) Initialize() {
	var err error
	ctx.Statik, err = fs.New()
	if err != nil {
		panic("err: " + err.Error())
	}

	dbConfig := config.Config.Database
	if len(dbConfig.Database) == 0 {
		return
	}
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	ctx.MybatisEngine = GoMybatis.GoMybatisEngine{}.New()
	ctx.DefaultSqlDB, err = ctx.MybatisEngine.Open("mysql", dsn)
	if err != nil {
		_ = fmt.Errorf("数据库链接失败, err: %v\n", err)
		return
	}
	ctx.DefaultSqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	ctx.DefaultSqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// step 绑定xml与mapper
	for _, bind := range ctx.MybatisMapperBinds {
		// 加载xml配置文件
		xml, err := ctx.Statik.Open(bind.XmlFile)
		if err != nil {
			panic("从statik加载xml配置失败")
		}
		// 记载xml内容到内存
		all, err := io.ReadAll(xml)
		if err != nil {
			panic("加载mybatis实现逻辑失败")
		}

		ctx.MybatisEngine.WriteMapperPtr(bind.Mapper, all)
	}
}

func (ts *MybatisInitializer) Close() {

}
