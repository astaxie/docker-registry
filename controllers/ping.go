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
}

func (this *PingController) Get() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", utils.Cfg.MustValue("docker", "XDockerRegistryStandalone"))
	pingResult := PingResult{Result: true}
	this.Data["json"] = &pingResult
	this.ServeJson()
}

func (this *PingController) GetPing() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", utils.Cfg.MustValue("docker", "XDockerRegistryStandalone"))
	pingResult := PingResult{Result: true}
	this.Data["json"] = &pingResult
	this.ServeJson()

}
