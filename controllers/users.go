package controllers

import (
	"github.com/astaxie/beego"
)

type UsersController struct {
	beego.Controller
}

func (this *UsersController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", beego.AppConfig.String("Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", beego.AppConfig.String("Standalone"))
}

// http://docs.docker.io/en/latest/reference/api/index_api/#user
// GET /users
// GET /users/
// If you want to check your login, you can try this endpoint
// Example Request:
//    GET /v1/users HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Authorization: Basic akmklmasadalkm==
// Example Response:
//    HTTP/1.1 200 OK
//    Vary: Accept
//    Content-Type: application/json
//    OK
// Status Codes: 
//    200 – no error
//    401 – Unauthorized
//    403 – Account is not Active
func (this *UsersController) GETUsers() {
	this.Ctx.Output.Body([]byte("\"OK\""))
}

// http://docs.docker.io/en/latest/reference/api/index_api/#user
// POST /users
// POST /users/
// Registering a new account.
// Example request:
//    POST /v1/users HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    {
//      "email": "sam@dotcloud.com",
//      "password": "toto42",
//      "username": "foobar"
//    }
// JSON Parameters:  
//    email – valid email address, that needs to be confirmed
//    username – min 4 character, max 30 characters, must match the regular expression [a-z0-9_].
//    password – min 5 characters
// Example Response:
//    HTTP/1.1 201 OK
//    Vary: Accept
//    Content-Type: application/json
//    "User Created"
// Status Codes: 
//    201 – User Created
//    400 – Errors (invalid json, missing or invalid fields, etc)
func (this *UsersController) POSTUsers() {

}

// http://docs.docker.io/en/latest/reference/api/index_api/#user
// PUT /v1/users/(username)/
// Change a password or email address for given user. If you pass in an email, it will add it to your account, it will not remove the old one. Passwords will be updated.
// It is up to the client to verify that that password that is sent is the one that they want. Common approach is to have them type it twice.
// Example Request:
//    PUT /v1/users/fakeuser/ HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Authorization: Basic akmklmasadalkm==
//    {
//      "email": "sam@dotcloud.com",
//      "password": "toto42"
//    }
// Parameters: 
//    username – username for the person you want to update
// Example Response:
//    HTTP/1.1 204
//    Vary: Accept
//    Content-Type: application/json
//    ""
// Status Codes: 
//    204 – User Updated
//    400 – Errors (invalid json, missing or invalid fields, etc)
//    401 – Unauthorized
//    403 – Account is not Active
//    404 – User not found
func (this *UsersController) PUTUsers() {

}
