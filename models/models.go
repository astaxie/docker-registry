package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var Engine *xorm.Engine

func setEngine() {
	host := utils.Cfg.MustValue("mysql", "Host")
	name := utils.Cfg.MustValue("mysql", "Name")
	user := utils.Cfg.MustValue("mysql", "User")
	passwd := utils.Cfg.MustValue("mysql", "Passwd")

	var err error
	conn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", user, passwd, host, name)
	beego.Trace("Initialized database connStr ->", conn)

	Engine, err = xorm.NewEngine("mysql", conn)
	if err != nil {
		log.Fatalf("models.init -> fail to conntect database: %v", err)
	}

	Engine.ShowDebug = true
	Engine.ShowErr = true
	Engine.ShowSQL = true

	beego.Trace("Initialized database ->", name)

}

// InitDb initializes the database.
func InitDb() {
	setEngine()
	err := Engine.Sync(new(User), new(Image), new(Repository))
	if err != nil {
		log.Fatalf("models.init -> fail to sync database: %v", err)
	}
}
