package controllers

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "regexp"
  "strconv"

  "github.com/astaxie/beego"
  "github.com/dockboard/docker-registry/models"
  "github.com/dockboard/docker-registry/utils"
)

type ImageController struct {
  beego.Controller
}

func (this *ImageController) Prepare() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "Version"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "Config"))
}

//在 Push 的流程中，docker 客户端会先调用 GET /v1/images/:image_id/json 向服务器检查是否已经存在 JSON 信息。
//如果存在了 JSON 信息，docker 客户端就认为是已经存在了 layer 数据，不再向服务器 PUT layer 的 JSON 信息和文件了。
//如果不存在 JSON 信息，docker 客户端会先后执行 PUT /v1/images/:image_id/json 和 PUT /v1/images/:image_id/layer 。
func (this *ImageController) GetImageJson() {

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization:" + this.Ctx.Input.Header("Authorization"))

  //判断用户的token是否可以操作
  r, _ := regexp.Compile(`Token signature=([[:alnum:]]+),repository="([[:alnum:]]+)/([[:alnum:]]+)",access=write`)
  authorizations := r.FindStringSubmatch(this.Ctx.Input.Header("Authorization"))
  beego.Trace("Token: " + authorizations[0])
  token, _, username, _ := authorizations[0], authorizations[1], authorizations[2], authorizations[3]

  user := &models.User{Username: username, Token: token}
  has, err := models.Engine.Get(user)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
    return
  }

  //查找是否有指定的ImageID对应的JSON文件
  image := &models.Image{ImageId: string(this.Ctx.Input.Param(":image_id"))}
  has, err = models.Engine.Get(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("Search image error."))
    return
  }

  if has == true {
    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
    this.Ctx.Output.Context.Output.Body([]byte(image.JSON))
    return
  } else {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("Image not found."))
    return
  }
}

//向数据库写入 Layer 的 JSON 数据
//TODO: 检查 JSON 是否合法
func (this *ImageController) PutImageJson() {

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())

  //判断用户的token是否可以操作
  beego.Trace("Authorization: " + this.Ctx.Input.Header("Authorization"))

  username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))

  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  beego.Trace("[Username & Password] " + username + " -> " + passwd)

  user := &models.User{Username: username, Password: passwd}
  has, err := models.Engine.Get(user)

  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  if user.Actived == false {
    this.Ctx.Output.Context.Output.SetStatus(403)
    this.Ctx.Output.Context.Output.Body([]byte("User is not actived."))
    return
  }

  //判断是否存在 image 的数据, 新建或更改 JSON 数据
  beego.Trace("[Image ID] " + this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: this.Ctx.Input.Param(":image_id")}

  has, err = models.Engine.Get(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Select image record error.\"}"))
    return
  }

  beego.Trace("[Has JSON] " + strconv.FormatBool(has))
  beego.Trace("[JSON] " + string(this.Ctx.Input.CopyBody()))
  image.JSON = string(this.Ctx.Input.CopyBody())

  if has == true {
    _, err = models.Engine.Id(image.Id).Cols("JSON").Update(image)
    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(400)
      this.Ctx.Output.Context.Output.Body([]byte("\"Update the image JSON data error.\""))
      return
    }
  }

  _, err = models.Engine.Insert(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Create the image record error.\""))
    return
  }

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))

}

//向本地硬盘写入 Layer 的文件
func (this *ImageController) PutImageIdLayer() {
  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())

  //判断用户的token是否可以操作
  beego.Trace("Authorization: " + this.Ctx.Input.Header("Authorization"))

  username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))

  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  beego.Trace("[Username & Password] " + username + " -> " + passwd)

  user := &models.User{Username: username, Password: passwd}
  has, err := models.Engine.Get(user)

  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  if user.Actived == false {
    this.Ctx.Output.Context.Output.SetStatus(403)
    this.Ctx.Output.Context.Output.Body([]byte("User is not actived."))
    return
  }

  //查询是否存在 image 的数据库记录
  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err = models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("\"Image not found\""))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }

  //处理 Layer 文件保存的目录
  basePath := utils.Cfg.MustValue("docker", "BasePath")
  repositoryPath := fmt.Sprintf("%v/images/%v", basePath, imageId)
  layerfile := fmt.Sprintf("%v/images/%v/layer", basePath, imageId)

  if !utils.IsDirExists(repositoryPath) {
    os.MkdirAll(repositoryPath, os.ModePerm)
  }

  //如果存在了文件就移除文件
  if _, err := os.Stat(layerfile); err == nil {
    os.Remove(layerfile)
  }

  //写入 Layer 文件
  data, _ := ioutil.ReadAll(this.Ctx.Request.Body)
  err = ioutil.WriteFile(layerfile, data, 0777)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Write image error.\""))
    return
  }

  //更新Image记录
  image.Uploaded = true
  _, err = models.Engine.Id(image.Id).Update(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Update the image upload status error.\"}"))
    return
  }

  //成功则返回 200
  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))
}

func (this *ImageController) PutChecksum() {

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization: " + this.Ctx.Input.Header("Authorization"))

  //判断用户的token是否可以操作
  //显示两个 Token 系 docker client 的 Bug 。

  username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))

  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  beego.Trace("[Username & Password] " + username + " -> " + passwd)

  user := &models.User{Username: username, Password: passwd}
  has, err := models.Engine.Get(user)

  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  if user.Actived == false {
    this.Ctx.Output.Context.Output.SetStatus(403)
    this.Ctx.Output.Context.Output.Body([]byte("User is not actived."))
    return
  }

  beego.Trace("Cookie: " + this.Ctx.Input.Header("Cookie"))
  beego.Trace("X-Docker-Checksum: " + this.Ctx.Input.Header("X-Docker-Checksum"))
  beego.Trace("X-Docker-Checksum-Payload: " + this.Ctx.Input.Header("X-Docker-Checksum-Payload"))

  //将 checksum 的值保存到数据库
  //Cookie: session=sFu7ZQLtC0EJPjH693JqWp61jL4=?checksum=KGxwMApTJ3NoYTI1NjplZTQwY2U4NGU2ZTA4NmIyM2E3ZDg0YzhkZTM0ZWU0YjcyYzgyZGEwMzI3ZmVlODVkZjkzY2M4NDRhMmM5ZmMzJwpwMQphLg==
  //X-Docker-Checksum: tarsum+sha256:6eb9bea3d03c72ec2f652869475e21bc11c0409d412c22ea5c44f371d02dda0b
  //X-Docker-Checksum-Payload: sha256:ee40ce84e6e086b23a7d84c8de34ee4b72c82da0327fee85df93cc844a2c9fc3

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err = models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("\"Image not found\""))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }

  image.Checksum = this.Ctx.Input.Header("X-Docker-Checksum")
  image.Payload = this.Ctx.Input.Header("X-Docker-Checksum-Payload")
  image.CheckSumed = true

  //计算这个Layer的父子结构
  var imageJSON map[string]interface{}
  if err := json.Unmarshal([]byte(image.JSON), &imageJSON); err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Decode the image json data encouter error.\""))
    return
  }

  var parents []string

  //存在 parent 的 ID
  if value, has := imageJSON["parent"]; has {
    parentImage := &models.Image{ImageId: value.(string)}
    has, err := models.Engine.Get(parentImage)
    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(400)
      this.Ctx.Output.Context.Output.Body([]byte("\"Check the parent image error.\""))
      return
    }

    if has {
      if err := json.Unmarshal([]byte(parentImage.ParentJSON), &parents); err != nil {
        this.Ctx.Output.Context.Output.SetStatus(400)
        this.Ctx.Output.Context.Output.Body([]byte("\"Decode the parent image json data encouter error.\""))
        return
      }
    }
  }

  var images []string
  images = append(images, imageId)
  parents = append(images, parents...)

  parentJSON, _ := json.Marshal(parents)
  image.ParentJSON = string(parentJSON)

  _, err = models.Engine.Id(image.Id).Update(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Update the image checksum error.\""))
    return
  }

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))

}

func (this *ImageController) GetImageAncestry() {
  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  //判断用户的token是否可以操作
  //Token 的样式类似于：Token Token signature=3d490a413351b26419beebf71b120759,repository="genedna/registry",access=read
  //显示两个 Token 系 docker client 的 Bug 。
  beego.Trace("Authorization: " + this.Ctx.Input.Header("Authorization"))
  // r, _ := regexp.Compile(`Token signature=([[:alnum:]]+),repository="([[:alnum:]]+)/([[:alnum:]]+)",access=read`)
  // authorizations := r.FindStringSubmatch(this.Ctx.Input.Header("Authorization"))
  // beego.Trace("Token: " + authorizations[0])
  // token, _, username, _ := authorizations[0], authorizations[1], authorizations[2], authorizations[3]

  // user := &models.User{Username: username, Token: token}
  // has, err := models.Engine.Get(user)
  // if has == false || err != nil {
  // 	this.Ctx.Output.Context.Output.SetStatus(401)
  // 	this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
  // 	return
  // }

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err := models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    return
  }

  if has {
    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
    this.Ctx.Output.Context.Output.SetStatus(200)
    this.Ctx.Output.Context.Output.Body([]byte(image.ParentJSON))
  }
}

func (this *ImageController) GetImageJSON() {
  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization:" + this.Ctx.Input.Header("Authorization"))

  //判断用户的 Authorization 是否可以操作
  username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))

  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  beego.Trace("[Username & Password] " + username + " -> " + passwd)

  user := &models.User{Username: username, Password: passwd}
  has, err := models.Engine.Get(user)

  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
    return
  }

  if user.Actived == false {
    this.Ctx.Output.Context.Output.SetStatus(403)
    this.Ctx.Output.Context.Output.Body([]byte("User is not actived."))
    return
  }

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err = models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    return
  }

  if has {
    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
    this.Ctx.Output.Context.Output.SetStatus(200)
    this.Ctx.Output.Context.Output.Body([]byte(image.JSON))
  }
}

func (this *ImageController) GetImageLayer() {
  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  //判断用户的token是否可以操作
  //Token 的样式类似于：Token Token signature=3d490a413351b26419beebf71b120759,repository="genedna/registry",access=write
  //显示两个 Token 系 docker client 的 Bug 。
  beego.Trace("Authorization: " + this.Ctx.Input.Header("Authorization"))
  // r, _ := regexp.Compile(`Token signature=([[:alnum:]]+),repository="([[:alnum:]]+)/([[:alnum:]]+)",access=write`)
  // authorizations := r.FindStringSubmatch(this.Ctx.Input.Header("Authorization"))
  // beego.Trace("Token: " + authorizations[0])
  // token, _, username, _ := authorizations[0], authorizations[1], authorizations[2], authorizations[3]

  // user := &models.User{Username: username, Token: token}
  // has, err := models.Engine.Get(user)
  // if has == false || err != nil {
  // 	this.Ctx.Output.Context.Output.SetStatus(401)
  // 	this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
  // 	return
  // }

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err := models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    return
  }

  if has == false {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("\"Could not find image record.\""))
  } else {
    //处理 Layer 文件保存的目录
    basePath := utils.Cfg.MustValue("docker", "BasePath")
    layerfile := fmt.Sprintf("%v/images/%v/layer", basePath, imageId)

    if _, err := os.Stat(layerfile); err != nil {
      this.Ctx.Output.Context.Output.SetStatus(404)
      this.Ctx.Output.Context.Output.Body([]byte("\"Could not find image file.\""))
    }

    file, err := ioutil.ReadFile(layerfile)
    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(400)
      this.Ctx.Output.Context.Output.Body([]byte("\"Load layer file error.\""))
    }

    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
    this.Ctx.Output.Context.Output.SetStatus(200)
    this.Ctx.Output.Context.Output.Body(file)

  }
}
