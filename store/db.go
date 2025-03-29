package store

import (
	"fmt"
	_ "modernc.org/sqlite"
	"strconv"
	"xorm.io/xorm"
)

var (
	filename = "./test.db"
)

type Config struct {
	Id       int64  `xorm:"id pk autoincr"` // 主键，自增
	Env      string `xorm:"env notnull"`
	Url      string `xorm:"url notnull"`
	Username string `xorm:"username notnull"`
	Password string `xorm:"password notnull"`
}

type Active struct {
	//Id  int64  `xorm:"id pk autoincr"` // 主键，自增
	Id int64 `xorm:"id notnull"`
}

func CloseDb(engine *xorm.Engine) {
	_ = engine.Close()
}
func InitDb() *xorm.Engine {
	// 打开 SQLite 数据库连接
	engine, err := xorm.NewEngine("sqlite", "./test.db")
	if err != nil {
		panic("初始化失败")
	}
	// 同步结构体到数据库表
	_ = engine.Sync2(new(Config))
	_ = engine.Sync2(new(Active))
	return engine
}

func Save(engine *xorm.Engine, config *Config) {
	has, _ := engine.Where("env=?", config.Env).Get(config)
	if !has {
		if _, err := engine.Insert(config); err != nil {
			panic("创建失败")
		}
	} else {
		Ls(engine)
	}

}

func Del(engine *xorm.Engine, id string) {
	i, _ := engine.ID(id).Delete(&Config{})
	if i > 0 {
		Ls(engine)
	}
}

func Get(engine *xorm.Engine, id int64) Config {
	var config Config
	_, _ = engine.ID(id).Get(&config)

	return config
}

func Ls(engine *xorm.Engine) {
	var configs []Config
	_ = engine.Find(&configs)
	fmt.Printf("%s\t%s\t%s\t%s\n", "序号", "环境", "地址", "授权")
	for _, config := range configs {
		fmt.Printf("%d\t%s\t%s\t%s:%s\n", config.Id, config.Env, config.Url, config.Username, config.Password)
	}
}

func Use(engine *xorm.Engine, id string) {

	// 清空表
	sql := "DELETE FROM active"
	_, _ = engine.Exec(sql)

	parseInt, _ := strconv.ParseInt(id, 10, 64)
	_, _ = engine.Insert(&Active{Id: parseInt})

}

func UseLs(engine *xorm.Engine) int64 {

	var data Active

	_, _ = engine.Get(&data)

	config := Get(engine, data.Id)

	fmt.Printf("%d\t%s\t%s\t%s:%s\n", config.Id, config.Env, config.Url, config.Username, config.Password)

	return data.Id
}

func GetUse(engine *xorm.Engine) int64 {

	var data Active

	_, _ = engine.Get(&data)

	return data.Id
}
