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

// GET /_ping or /v1/_ping
// API Spec GET /_ping http://docs.docker.io/en/latest/reference/api/registry_api
// Section 2.4 Status
func (this *PingController) Get() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", "0.6.5")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", "true")
  pingResult := PingResult{Result: true}
  this.Data["json"] = &pingResult
  this.ServeJson()
}
