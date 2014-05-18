package controllers

import (
  "github.com/astaxie/beego"
)

type PingController struct {
  beego.Controller
}

type PingResult struct {
  Result bool
}

func (this *PingController) Prepare() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

func (this *PingController) Get() {
  pingResult := PingResult{Result: true}
  this.Data["json"] = &pingResult
  this.ServeJson()
}
