package routers

import (
  "github.com/astaxie/beego"
  "github.com/docker-registry/controllers"
)

func init() {
  beego.Router("/", &controllers.MainController{})
  //Status
  beego.Router("/_ping", &controllers.PingController{})
  beego.Router("/v1/_ping", &controllers.PingController{})
  //Users
  beego.Router("/users", &controllers.UsersController{})
}
