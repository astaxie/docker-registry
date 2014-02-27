package controllers

import (
  "github.com/astaxie/beego"
)

type UsersController struct {
  beego.Controller
}

func (this *UsersController) Prepare() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", "0.6.5")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", "true")
}

func (this *UsersController) Get() {
  this.Ctx.Output.Body([]byte("\"OK\""))
}
