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

func (this *ImageController) GETPrivateLayer() {

}

func (this *ImageController) GETLayer() {

}

func (this *ImageController) PUTLayer() {

}

func (this *ImageController) PUTChecksum() {

}

func (this *ImageController) GETPrivateJSON() {

}

func (this *ImageController) GETJSON() {

}

func (this *ImageController) GETAncestry() {

}

func (this *ImageController) PUTJSON() {

}

func (this *ImageController) GETPrivateFiles() {

}

func (this *ImageController) GETFiles() {
	
}