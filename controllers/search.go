package controllers

import (
	"github.com/astaxie/beego"
)

type SearchController struct {
	beego.Controller
}

func (this *SearchController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

// http://docs.docker.io/en/latest/reference/api/index_api/#search
// GET /v1/search
// Search the Index given a search term. It accepts GET only.
// Example request:
//    GET /v1/search?q=search_term HTTP/1.1
//    Host: example.com
//    Accept: application/json
// Example response:
//    HTTP/1.1 200 OK
//    Vary: Accept
//    Content-Type: application/json
//    {
//      "query":"search_term",
//      "num_results": 3,
//      "results" : [
//          {"name": "ubuntu", "description": "An ubuntu image..."},
//          {"name": "centos", "description": "A centos image..."},
//          {"name": "fedora", "description": "A fedora image..."}
//      ]
//    }
// Query Parameters:  
//    q – what you want to search for
// Status Codes: 
//    200 – no error
//    500 – server error
func (this *SearchController) GET() {

}
