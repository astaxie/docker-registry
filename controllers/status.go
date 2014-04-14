package controllers

import (
	"github.com/astaxie/beego"
)

type StatusController struct {
	beego.Controller
}

func (this *StatusController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

func (this *StatusController) GET() {

}
