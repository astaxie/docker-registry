// Docker Registry & Login
// 执行 docker login 命令流程：
//     1. docker 向 registry 的服务器进行注册执行：POST /v1/users or /v1/users/ -> POSTUsers()
//     2. 创建用户成功返回 201；提交的格式有误、无效的字段等返回 400；已经存在用户了返回 401。
//     3. docker login 收到 401 的状态后，进行登录：GET /v1/users or /v1/users/ -> GETUsers()
//     4. 在登录时，将用户名和密码进行 SetBasicAuth 处理，放到 HEADER 的 Authorization 中，例如：Authorization: Basic ZnNrOmZzaw==
//     5. registry 收到登录的请求，Decode 请求 HEADER 中 Authorization 的部分进行判断。
//     6. 用户名和密码正确返回 200；用户名密码错误返回 401；账户未激活返回 403 错误；其它错误返回 500。
// 注：
//     Decode HEADER authorization function named decodeAuth in https://github.com/dotcloud/docker/blob/master/registry/auth.go.
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

// http://docs.docker.io/en/latest/reference/api/index_api/#users
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

// http://docs.docker.io/en/latest/reference/api/index_api/#users
// POST /v1/users
// POST /v1/users/
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
//    401 - 已经存在此用户
func (this *UsersController) POSTUsers() {

}

// http://docs.docker.io/en/latest/reference/api/index_api/#users
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
