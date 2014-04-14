package controllers

import (
	"github.com/astaxie/beego"
)

type SearchController struct {
	beego.Controller
}

func (this *SearchController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

func (this *SearchController) GET() {

}
