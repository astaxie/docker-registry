package main

import (
	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/models"
	_ "github.com/dockboard/docker-registry/routers"
	"github.com/dockboard/docker-registry/utils"
)

func main() {
	utils.LoadConfig("conf/app.conf")

	models.InitDb()
	beego.Run()
}
