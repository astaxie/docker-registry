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

func (this *RepositoryController) PUTProperties() {

}

func (this *RepositoryController) GETProperties() {

}

func (this *RepositoryController) GETTags() {

}

func (this *RepositoryController) GETTag() {

}

func (this *RepositoryController) GETRepositoyJSON() {

}

func (this *RepositoryController) GETTagJSON() {

}

func (this *RepositoryController) PUTTag() {

}

func (this *RepositoryController) DELETETag() {

}

func (this *RepositoryController) DELETERepository() {

}

func (this *RepositoryController) DELETERepositoryTags() {

}
