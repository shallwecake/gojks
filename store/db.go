package store

import (
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strconv"
	"xorm.io/xorm"
)

var (
	// 清空表
	sql1 = "DROP TABLE active"
	sql2 = "DROP TABLE config"
	sql3 = "DELETE FROM active"
)

type Config struct {
	Id       int64  `xorm:"id pk autoincr"` // 主键，自增
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("获取用户目录失败: " + err.Error())
	}

	filename := filepath.Join(homeDir, "jconfig.db") // 自动处理分隔符

	// 打开 SQLite 数据库连接
	engine, err := xorm.NewEngine("sqlite", filename)
	if err != nil {
		panic("初始化失败")
	}
	// 同步结构体到数据库表
	_ = engine.Sync2(new(Config))
	_ = engine.Sync2(new(Active))
	return engine
}

func Save(engine *xorm.Engine, config *Config) {
	has, _ := engine.Where("url=?", config.Url).Get(config)
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
	ls := Ls(engine)

	if ls == 0 {
		_, _ = engine.Exec(sql1)
		_, _ = engine.Exec(sql2)
	}
}

func Get(engine *xorm.Engine, id int64) *Config {
	var config Config
	_, _ = engine.ID(id).Get(&config)
	return &config
}

func Ls(engine *xorm.Engine) int {
	var configs []Config
	_ = engine.Find(&configs)
	if len(configs) != 0 {
		fmt.Printf("%s\t%s\t%s\n", "序号", "地址", "授权")
		for _, config := range configs {
			fmt.Printf("%d\t%s\t%s:%s\n", config.Id, config.Url, config.Username, config.Password)
		}
	}
	return len(configs)
}

func Use(engine *xorm.Engine, id string) {

	_, _ = engine.Exec(sql3)

	parseInt, _ := strconv.ParseInt(id, 10, 64)
	_, _ = engine.Insert(&Active{Id: parseInt})

}

func UseLs(engine *xorm.Engine) int64 {

	var data Active

	_, _ = engine.Get(&data)

	config := Get(engine, data.Id)

	if len(config.Url) > 0 {
		fmt.Printf("%d\t%s\t%s:%s\n", config.Id, config.Url, config.Username, config.Password)
	}

	return data.Id
}

func GetUse(engine *xorm.Engine) int64 {

	var data Active

	_, _ = engine.Get(&data)

	return data.Id
}
