package routers

import (
  "github.com/astaxie/beego"
  "github.com/dockboard/docker-registry/controllers"
)

func init() {
  beego.Router("/", &controllers.MainController{})
  beego.Router("/_status", &controllers.StatusController{})

  beego.Router("/_ping", &controllers.PingController{}, "get:GetPing")

  ns := beego.NewNamespace("/v1",
    beego.NSRouter("/_ping", &controllers.PingController{}, "get:GetPing"),

    beego.NSRouter("/_status", &controllers.StatusController{}),

    beego.NSRouter("/users", &controllers.UsersController{}, "get:GetUsers"),
    beego.NSRouter("/users", &controllers.UsersController{}, "post:PostUsers"),

    beego.NSNamespace("/repositories",
      beego.NSRouter("/:namespace/:repo_name/tags/:tag", &controllers.RepositoryController{}, "put:PutTag"),
      beego.NSRouter("/:namespace/:repo_name/tags/:tag", &controllers.RepositoryController{}, "put:PutTag"),
      beego.NSRouter(":namespace/:repo_name/images", &controllers.RepositoryController{}, "put:PutRepositoryImages"),
      beego.NSRouter("/:namespace/:repo_name/images", &controllers.RepositoryController{}, "get:GetRepositoryImages"),
      beego.NSRouter("/:namespace/:repo_name/tags", &controllers.RepositoryController{}, "get:GetRepositoryTags"),
      beego.NSRouter("/:namespace/:repo_name/tags", &controllers.RepositoryController{}, "get:GetRepositoryTags"),
      beego.NSRouter("/:namespace/:repo_name", &controllers.RepositoryController{}, "put:PutRepository"),
    ),

    beego.NSNamespace("/images",
      beego.NSRouter("/:image_id/ancestry", &controllers.ImageController{}, "get:GetImageAncestry"),
      beego.NSRouter("/:image_id/json", &controllers.ImageController{}, "get:GetImageJSON"),
      beego.NSRouter("/:image_id/layer", &controllers.ImageController{}, "get:GetImageLayer"),
      beego.NSRouter("/:image_id/json", &controllers.ImageController{}, "put:PutImageJson"),
      beego.NSRouter("/:image_id/layer", &controllers.ImageController{}, "put:PutImageLayer"),
      beego.NSRouter("/:image_id/checksum", &controllers.ImageController{}, "put:PutChecksum"),
    ),
  )

  beego.AddNamespace(ns)
}
