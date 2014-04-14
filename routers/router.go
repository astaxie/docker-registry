package routers

import (
	"github.com/astaxie/beego"
	"github.com/docker-registry/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//Status
	beego.Router("/_ping", &controllers.PingController{})
	beego.Router("/v1/_ping", &controllers.PingController{})
	//Users
	beego.Router("/users", &controllers.UsersController{}, "get:GETUsers")
	beego.Router("/users", &controllers.UsersController{}, "post:POSTUsers")
	beego.Router("/users/:username", &controllers.UsersController{}, "put:PUTUsers")
	//Repository
	beego.Router("/v1/repositories/:username/:repository", &controllers.RepositoryController{}, "put:PUTRepository")
	beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "put:PUTRepositoryImages")
	beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "get:GETRepositoryImages")
	beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "delete:DELETERepositoryImages")
	beego.Router("/v1/repositories/:username/:repository/auth", &controllers.RepositoryController{}, "put:PUTRepositoryAuth")
	//Images
	beego.Router("/v1/private_images/:image_id/layer", &controllers.ImageController{}, "get:GETPrivateLayer")
	beego.Router("/v1/images/:image_id/layer", &controllers.ImageController{}, "get:GETLayer")
	beego.Router("/v1/images/:image_id/layer", &controllers.ImageController{}, "put:PUTLayer")
	beego.Router("/v1/images/:image_id/checksum", &controllers.ImageController{}, "put:PUTChecksum")
	beego.Router("/v1/private_images/:image_id/json", &controllers.ImageController{}, "get:GETPrivateJSON")
	beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "get:GETJSON")
	beego.Router("/v1/images/:image_id/ancestry", &controllers.ImageController{}, "get:GETAncestry")
	beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "put:PUTJSON")
	beego.Router("/v1/private_images/:image_id/files", &controllers.ImageController{}, "get:GETPrivateFiles")
	beego.Router("/v1/images/:image_id/files", &controllers.ImageController{}, "get:GETFiles")
}
