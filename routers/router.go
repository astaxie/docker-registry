package routers

import (
  "github.com/astaxie/beego"
  "github.com/docker-registry/controllers"
)

func init() {
  beego.Router("/", &controllers.MainController{})
}
