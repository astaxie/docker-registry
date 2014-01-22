package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Website"] = "Docker Registry With Beego Framework"
	this.Data["Email"] = "genedna@gmail.com"
	this.TplNames = "index.tpl"
}
