package controllers

import (
	"github.com/astaxie/beego"
)

type UsersController struct {
	beego.Controller
}

func (this *UsersController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

func (this *UsersController) GETUsers() {
	this.Ctx.Output.Body([]byte("\"OK\""))
}

func (this *UsersController) POSTUsers() {

}

func (this *UsersController) PUTUsers() {

}
