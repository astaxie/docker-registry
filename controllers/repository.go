package controllers

import (
	"github.com/astaxie/beego"
)

type RepositoryController struct {
	beego.Controller
}

func (this *RepositoryController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

func (this *RepositoryController) PUTRepository() {

}

func (this *RepositoryController) PUTRepositoryImages() {

}

func (this *RepositoryController) GETRepositoryImages() {

}

func (this *RepositoryController) DELETERepositoryImages() {

}

func (this *RepositoryController) PUTRepositoryAuth() {

}
