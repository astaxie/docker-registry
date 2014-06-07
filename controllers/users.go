/*
Docker Registry & Login
执行 docker login 命令流程：
    1. docker 向 registry 的服务器进行注册执行：POST /v1/users or /v1/users/ -> POSTUsers()
    2. 创建用户成功返回 201；提交的格式有误、无效的字段等返回 400；已经存在用户了返回 401。
    3. docker login 收到 401 的状态后，进行登录：GET /v1/users or /v1/users/ -> GETUsers()
    4. 在登录时，将用户名和密码进行 SetBasicAuth 处理，放到 HEADER 的 Authorization 中，例如：Authorization: Basic ZnNrOmZzaw==
    5. registry 收到登录的请求，Decode 请求 HEADER 中 Authorization 的部分进行判断。
    6. 用户名和密码正确返回 200；用户名密码错误返回 401；账户未激活返回 403 错误；其它错误返回 417 (Expectation Failed)
注：
    Decode HEADER authorization function named decodeAuth in https://github.com/dotcloud/docker/blob/master/registry/auth.go.
更新 Docker Registry User 的属性：
    1. 调用 PUT /v1/users/(username)/ 向服务器更新 User 的 Email 和 Password 属性。
    2. 参数包括 User Email 或 User Password，或两者都包括。
    3. 更新成功返回 204；传递的参数不是有效的 JSON 格式等错误返回 400；认证失败返回 401；用户没有激活返回 403；没有用户现实 404。
注：
    HTTP HEADER authorization decode 验证同 docker login 命令。
*/
package controllers

import (
  "github.com/astaxie/beego"
  "github.com/dockboard/docker-registry/models"
  "github.com/dockboard/docker-registry/utils"
)

type UsersController struct {
  beego.Controller
}

func (this *UsersController) Prepare() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "Version"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "Config"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Standalone", utils.Cfg.MustValue("docker", "Standalone"))
}

func (this *UsersController) PostUsers() {

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization:" + this.Ctx.Input.Header("Authorization"))

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(401)
  this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"We are limited beta testing now. Please contact Meaglith Ma <genedna@gmail.com> for beta account.\"}"))
}

func (this *UsersController) GetUsers() {

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization:" + this.Ctx.Input.Header("Authorization"))

  username, password, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401) //根据 Specification ，解码 Basic Authorization 数据失败也认为是 401 错误。
    this.Ctx.Output.Body([]byte("\"Unauthorized\""))
    this.StopRun()
  }

  beego.Trace("[Username & Password] " + username + " -> " + password)

  var has bool
  user := &models.User{Username: username, Password: password}
  has, err = models.Engine.Get(user)

  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Body([]byte("\"Unauthorized\""))
    this.StopRun()
  }

  if has == false {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Body([]byte("\"No User\""))
    this.StopRun()
  }

  if user.Actived == false {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Body([]byte("\"Actived First!\""))
    this.StopRun()
  }

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte("{\"OK\"}"))
}
