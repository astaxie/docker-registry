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

// http://docs.docker.io/en/latest/reference/api/registry_api/#tags
// GET /v1/repositories/(namespace)/(repository)/tags
// get all of the tags for the given repo.
// Example Request:
//    GET /v1/repositories/foo/bar/tags HTTP/1.1
//    Host: registry-1.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    Cookie: (Cookie provided by the Registry)
// Parameters: 
//    namespace – namespace for the repo
//    repository – name for the repo
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    {
//      "latest": "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f",
//      "0.1.1":  "b486531f9a779a0c17e3ed29dae8f12c4f9e89cc6f0bc3c38722009fe6857087"
//    }
// Status Codes: 
//    200 – OK
//    401 – Requires authorization
//    404 – Repository not found
func (this *RepositoryController) GETTags() {

}

// http://docs.docker.io/en/latest/reference/api/registry_api/#tags
// GET /v1/repositories/(namespace)/(repository)/tags/(tag)
// get a tag for the given repo.
// Example Request:
//    GET /v1/repositories/foo/bar/tags/latest HTTP/1.1
//    Host: registry-1.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    Cookie: (Cookie provided by the Registry)
// Parameters: 
//    namespace – namespace for the repo
//    repository – name for the repo
//    tag – name of tag you want to get
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f"
// Status Codes: 
//    200 – OK
//    401 – Requires authorization
//    404 – Tag not found
func (this *RepositoryController) GETTag() {

}

// http://docs.docker.io/en/latest/reference/api/registry_api/#tags
// DELETE /v1/repositories/(namespace)/(repository)/tags/(tag)
// delete the tag for the repo
// Example Request:
//    DELETE /v1/repositories/foo/bar/tags/latest HTTP/1.1
//    Host: registry-1.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Cookie: (Cookie provided by the Registry)
// Parameters: 
//    namespace – namespace for the repo
//    repository – name for the repo
//    tag – name of tag you want to delete
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    ""
// Status Codes: 
//    200 – OK
//    401 – Requires authorization
//    404 – Tag not found
func (this *RepositoryController) DELETETag() {

}

// http://docs.docker.io/en/latest/reference/api/registry_api/#tags
// PUT /v1/repositories/(namespace)/(repository)/tags/(tag)
// put a tag for the given repo.
// Example Request:
//    PUT /v1/repositories/foo/bar/tags/latest HTTP/1.1
//    Host: registry-1.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Cookie: (Cookie provided by the Registry)
//    "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f"
// Parameters: 
//    namespace – namespace for the repo
//    repository – name for the repo
//    tag – name of tag you want to add
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    ""
// Status Codes: 
//    200 – OK
//    400 – Invalid data
//    401 – Requires authorization
//    404 – Image not found
func (this *RepositoryController) PUTTag() {

}

// http://docs.docker.io/en/latest/reference/api/registry_api/#repositories
// DELETE /v1/repositories/(namespace)/(repository)/
// delete a repository
// Example Request:
//    DELETE /v1/repositories/foo/bar/ HTTP/1.1
//    Host: registry-1.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Cookie: (Cookie provided by the Registry)
//    ""
// Parameters: 
//    namespace – namespace for the repo
//    repository – name for the repo
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    X-Docker-Registry-Version: 0.6.0
//    ""
// Status Codes: 
//    200 – OK
//    401 – Requires authorization
//    404 – Repository not found
func (this *RepositoryController) DELETERepositoryImages() {

}

// http://docs.docker.io/en/latest/reference/api/index_api/#repository

func (this *RepositoryController) PUTRepository() {

}

func (this *RepositoryController) PUTRepositoryImages() {

}

func (this *RepositoryController) GETRepositoryImages() {

}



func (this *RepositoryController) PUTRepositoryAuth() {

}

func (this *RepositoryController) PUTProperties() {

}

func (this *RepositoryController) GETProperties() {

}



func (this *RepositoryController) GETRepositoyJSON() {

}

func (this *RepositoryController) GETTagJSON() {

}





func (this *RepositoryController) DELETERepository() {

}

func (this *RepositoryController) DELETERepositoryTags() {

}
