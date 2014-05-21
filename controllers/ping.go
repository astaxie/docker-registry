package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/utils"
)

type PingController struct {
	beego.Controller
}

type PingResult struct {
	Result bool
}

func (this *PingController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "Config"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", utils.Cfg.MustValue("docker", "Standalone"))
}

func (this *PingController) GetPing() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	pingResult := PingResult{Result: true}
	this.Data["json"] = &pingResult
	this.ServeJson()
}
