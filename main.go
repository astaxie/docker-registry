package main

import (
	"github.com/astaxie/beego"
	"github.com/docker-registry/models"
	_ "github.com/docker-registry/routers"
	"github.com/docker-registry/utils"
)

func main() {
	utils.LoadConfig("conf/app.conf")

	models.InitDb()
	beego.Run()
}
