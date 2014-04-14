package controllers

package controllers

import (
	"github.com/astaxie/beego"
)

type ImageController struct {
	beego.Controller
}

func (this *ImageController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}