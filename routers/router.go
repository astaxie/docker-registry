package routers

import (
	"github.com/astaxie/beego"
	"github.com/docker-registry/controllers"
)

//
// Registry APIs http://docs.docker.io/en/latest/reference/api/registry_api/
// Index APIs http://docs.docker.io/en/latest/reference/api/index_api/
//
func init() {
	// Index
	beego.Router("/", &controllers.MainController{})
  // http://docs.docker.io/en/latest/reference/api/registry_api/#images
  // Documented and implemented in docker-registry
  beego.Router("/v1/images/:image_id/layer", &controllers.ImageController{}, "get:GETLayer")
  beego.Router("/v1/images/:image_id/layer", &controllers.ImageController{}, "put:PUTLayer")
  beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "get:GETJSON")
  beego.Router("/v1/images/:image_id/json", &controllers.ImageController{}, "put:PUTJSON")
  beego.Router("/v1/images/:image_id/ancestry", &controllers.ImageController{}, "get:GETAncestry")
  // Undocumented but implemented in docker-registry and docker-index
  beego.Router("/v1/images/:image_id/checksum", &controllers.ImageController{}, "put:PUTChecksum")  
  beego.Router("/v1/images/:image_id/files", &controllers.ImageController{}, "get:GETFiles")
  beego.Router("/v1/images/:image_id/diff", &controllers.ImageController{}, "get:GETDiff")
  // Private images and undocumented
  beego.Router("/v1/private_images/:image_id/layer", &controllers.ImageController{}, "get:GETPrivateLayer")
  beego.Router("/v1/private_images/:image_id/json", &controllers.ImageController{}, "get:GETPrivateJSON")
  beego.Router("/v1/private_images/:image_id/files", &controllers.ImageController{}, "get:GETPrivateFiles")
  // http://docs.docker.io/en/latest/reference/api/registry_api/#tags
  // Documented and implemented in docker-registry
  beego.Router("/v1/repositories/:username/:repository/tags", &controllers.RepositoryController{}, "get:GETTags")
  beego.Router("/v1/repositories/:username/:repository/tags/:tag", &controllers.RepositoryController{}, "get:GETTag")
  beego.Router("/v1/repositories/:username/:repository/tags/:tag", &controllers.RepositoryController{}, "delete:DELETETag")
  beego.Router("/v1/repositories/:username/:repository/tags/:tag", &controllers.RepositoryController{}, "put:PUTTag")
  // http://docs.docker.io/en/latest/reference/api/registry_api/#repositories
  // Documented and implemented in docker-registry
  beego.Router("/v1/repositories/:username/:repository/", &controllers.RepositoryController{}, "delete:DELETERepository")
  // Documented and implemented in docker-index
  beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "delete:DELETERepositoryImages")
  beego.Router("/v1/repositories/:username/:repository/properties", &controllers.RepositoryController{}, "put:PUTProperties")
  beego.Router("/v1/repositories/:username/:repository/properties", &controllers.RepositoryController{}, "get:GETProperties")
  beego.Router("/v1/repositories/:username/:repository/json", &controllers.RepositoryController{}, "get:GETRepositoryJSON")
  // http://docs.docker.io/en/latest/reference/api/index_api/#repository
  // Documented and implemented in docker-index
  beego.Router("/v1/repositories/:username/:repository", &controllers.RepositoryController{}, "put:PUTRepository")
  // http://docs.docker.io/en/latest/reference/api/index_api/#repository-images
  // Documented and implemented in docker-index
  beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "put:PUTRepositoryImages")
  beego.Router("/v1/repositories/:username/:repository/images", &controllers.RepositoryController{}, "get:GETRepositoryImages")
  // http://docs.docker.io/en/latest/reference/api/index_api/#repository-authorization
  // Documented and implemented in docker-index
  beego.Router("/v1/repositories/:username/:repository/auth", &controllers.RepositoryController{}, "put:PUTRepositoryAuth")
  // Undocumented but implemented in docker-registry and docker-index
  beego.Router("/v1/repositories/:username/:repository/tags/:tag/json", &controllers.RepositoryController{}, "get:GETTagJSON")  
  beego.Router("/v1/repositories/:username/:repository/tags", &controllers.RepositoryController{}, "delete:DELETERepositoryTags")
  // http://docs.docker.io/en/latest/reference/api/registry_api/#status
  // Documented and implemented in docker-registry
  beego.Router("/_ping", &controllers.PingController{})
  beego.Router("/v1/_ping", &controllers.PingController{})
  // Undocumented but implemented in docker-registry
  beego.Router("/_status", &controllers.StatusController{})
  beego.Router("/v1/_status", &controllers.StatusController{})
  // http://docs.docker.io/en/latest/reference/api/index_api/#user-login
  // Documented and implemented in docker-index
  beego.Router("/v1/users", &controllers.UsersController{}, "get:GETUsers")
  beego.Router("/v1/users", &controllers.UsersController{}, "post:POSTUsers")
  beego.Router("/v1/users/", &controllers.UsersController{}, "get:GETUsers")
  beego.Router("/v1/users/", &controllers.UsersController{}, "post:POSTUsers")   
  beego.Router("/v1/users/:username", &controllers.UsersController{}, "put:PUTUsers")
  // http://docs.docker.io/en/latest/reference/api/index_api/#search
  // Documented and implemented in docker-index  
	beego.Router("/v1/search", &controllers.SearchController{})	
}
