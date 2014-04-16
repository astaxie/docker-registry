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

// GET /_ping or /v1/_ping
// API Spec GET /_ping http://docs.docker.io/en/latest/reference/api/registry_api/#status
// Every 'docker pull' or 'docker push' command will access /v1/_ping before access other URLs.
// The docker client call at docker/registry/registry.go@pingRegistryEndpoint function.
// if 'X-Docker-Registry-Version' is blank, the docker client function will return error.
// if 'X-Docker-Registry-Standalone' is blank, the docker client will assume it's a earlier registry version.
// So 'X-Docker-Registry-Standalone' must be 'true' or '1'.
func (this *PingController) Get() {
	pingResult := PingResult{Result: true}
	this.Data["json"] = &pingResult
	this.ServeJson()
}
