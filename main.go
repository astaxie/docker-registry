package main

import (
	_ "docker-registry/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
