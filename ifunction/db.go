package ifunction

import (
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"xorm.io/xorm"
)

var (
	// 清空表
	whereSql = "type = ?"
	Engine   *xorm.Engine
)

const (
	Web_Hook = "whk"
	Jen_kins = "jks"
	Ran_Cher = "rcr"
)

type Conf struct {
	Type     string
	Url      string `xorm:"url notnull"`
	Username string `xorm:"username notnull"`
	Password string `xorm:"password notnull"`
}

func CloseDbEngine(engine *xorm.Engine) {
	_ = engine.Close()
}

func InitDb() *xorm.Engine {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("获取用户目录失败: " + err.Error())
	}

	filename := filepath.Join(homeDir, "gojks.db") // 自动处理分隔符

	// 打开 SQLite 数据库连接
	Engine, err = xorm.NewEngine("sqlite", filename)
	if err != nil {
		panic("初始化失败")
	}
	// 同步结构体到数据库表
	_ = Engine.Sync2(new(Conf))
	return Engine
}

func SaveOrUpdate(engine *xorm.Engine, conf *Conf) {
	// 查询是否已存在记录
	existing := new(Conf)
	has, err := engine.Where(whereSql, conf.Type).Get(existing)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	if has {
		// 如果存在，更新记录
		_, err := engine.Table(new(Conf)).
			Where(whereSql, conf.Type).
			Update(conf)
		if err != nil {
			log.Fatalf("更新失败: %v", err)
		}
	} else {
		// 如果不存在，插入新记录
		_, err := engine.Insert(conf)
		if err != nil {
			log.Fatalf("插入失败: %v", err)
		}
	}
}

func DelConf(engine *xorm.Engine, category string) {
	engine.Where(whereSql, category).Delete(new(Conf))
}

func GetConf(engine *xorm.Engine, category string) *Conf {
	c := new(Conf)
	_, _ = engine.Where(whereSql, category).Get(c)
	return c
}

func ListConf(engine *xorm.Engine) {
	var configs []Conf
	_ = engine.Find(&configs)
	if len(configs) != 0 {
		fmt.Printf("%s\t%s\t%s\n", "名称", "地址", "授权")
		for _, config := range configs {
			fmt.Printf("%s\t%s\t%s\n", config.Type, config.Url, getAuth(config))
		}
	}
}

func getAuth(conf Conf) string {
	if len(conf.Username) == 0 {
		return ""
	} else {
		return conf.Password + ":" + conf.Password
	}
}
