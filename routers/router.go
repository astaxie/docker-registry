package routers

import (
  "github.com/astaxie/beego"
  "github.com/dockboard/docker-registry/controllers"
)

func init() {
  beego.Router("/", &controllers.MainController{})

  //Ping Router
  beego.Router("/_ping", &controllers.PingController{}, "get:GetPing")
  beego.Router("/_ping/", &controllers.PingController{}, "get:GetPing")
  beego.Router("/v1/_ping", &controllers.PingController{}, "get:GetPing")
  beego.Router("/v1/_ping/", &controllers.PingController{}, "get:GetPing")

  //User Router
  beego.Router("/v1/users/", &controllers.UsersController{}, "get:GetUsers")
  beego.Router("/v1/users/", &controllers.UsersController{}, "post:PostUsers")
  beego.Router("/v1/users", &controllers.UsersController{}, "get:GetUsers")
  beego.Router("/v1/users", &controllers.UsersController{}, "post:PostUsers")

  //Push Router
  beego.Router("/v1/repositories/:namespace/:repo_name/", &controllers.RepositoryController{}, "put:PutNamespaceRepository")
  beego.Router("/v1/repositories/:namespace/:repo_name", &controllers.RepositoryController{}, "put:PutNamespaceRepository")

  beego.Router("/v1/repositories/:namespace/:repository/tags/:tag", &controllers.RepositoryController{}, "put:PutNamespaceTag")
  beego.Router("/v1/repositories/:repository/tags/:tag", &controllers.RepositoryController{}, "put:PutTag")
  beego.Router("/v1/repositories/:namespace/:repository/images", &controllers.RepositoryController{}, "put:PutNamespaceImages")
  beego.Router("/v1/repositories/:repository/images", &controllers.RepositoryController{}, "put:PutImages")

  beego.Router("/v1/repositories/:repo_name/", &controllers.RepositoryController{}, "put:PutRepo")
  beego.Router("/v1/repositories/:repo_name", &controllers.RepositoryController{}, "put:PutRepo")

  beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "get:GetImageJson")
  beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "put:PutImageJson")

  beego.Router("/v1/images/:image_id/layer", &controllers.ImageController{}, "put:PutImageIdLayer")

  beego.Router("/v1/images/:image_id/checksum", &controllers.ImageController{}, "put:PutChecksum")

  beego.Router("/_status", &controllers.StatusController{})
  beego.Router("/v1/_status", &controllers.StatusController{})
}
