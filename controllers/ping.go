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

func (this *PingController) Get() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", "0.6.0")
  pingResult := PingResult{Result: true}
  this.Data["json"] = &pingResult
  this.ServeJson()
}
