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

  beego.Trace(this.Ctx.Request.Method + " -> " + this.Ctx.Request.URL.String())
  beego.Trace("Authorization:" + this.Ctx.Input.Header("Authorization"))

  r, _ := regexp.Compile(`Token signature=([[:alnum:]]+),repository="([[:alnum:]]+)/([[:alnum:]]+)",access=([[:alnum:]]+)`)
  authorizations := r.FindStringSubmatch(this.Ctx.Input.Header("Authorization"))

  if len(authorizations) == 5 {
    token, _, username, _, access := authorizations[0], authorizations[1], authorizations[2], authorizations[3], authorizations[4]

    user := &models.User{Username: username, Token: token}
    has, err := models.Engine.Get(user)

    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(401)
      this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
      this.StopRun()
    }

    if has == false || user.Actived == false {
      this.Ctx.Output.Context.Output.SetStatus(403)
      this.Ctx.Output.Context.Output.Body([]byte("User is not exist or not actived."))
      this.StopRun()
    }

    this.Data["user"] = user
    this.Data["access"] = access
  }

  if len(authorizations) == 0 {
    //判断用户的 Authorization 是否可以操作
    username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))

    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(401)
      this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
      this.StopRun()
    }

    beego.Trace("[Username & Password] " + username + " -> " + passwd)

    user := &models.User{Username: username, Password: passwd}
    has, err := models.Engine.Get(user)

    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(401)
      this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
      this.StopRun()
    }

    if has == false || user.Actived == false {
      this.Ctx.Output.Context.Output.SetStatus(403)
      this.Ctx.Output.Context.Output.Body([]byte("User is not exist or not actived."))
      this.StopRun()
    }

    this.Data["user"] = user
  }
}

//在 Push 的流程中，docker 客户端会先调用 GET /v1/images/:image_id/json 向服务器检查是否已经存在 JSON 信息。
//如果存在了 JSON 信息，docker 客户端就认为是已经存在了 layer 数据，不再向服务器 PUT layer 的 JSON 信息和文件了。
//如果不存在 JSON 信息，docker 客户端会先后执行 PUT /v1/images/:image_id/json 和 PUT /v1/images/:image_id/layer 。
func (this *ImageController) GetImageJSON() {

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId, Uploaded: true, CheckSumed: true}
  has, err := models.Engine.Get(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    this.StopRun()
  }

  beego.Trace("[Image Has] " + strconv.FormatBool(has))

  if has {
    beego.Trace("[Image JSON] " + string(image.JSON))
    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
    this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Checksum", image.Checksum)
    this.Ctx.Output.Context.Output.SetStatus(200)
    this.Ctx.Output.Context.Output.Body([]byte(image.JSON))
    this.StopRun()
  } else {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("\"No image json.\""))
    this.StopRun()
  }
}

//向数据库写入 Layer 的 JSON 数据
//TODO: 检查 JSON 是否合法
func (this *ImageController) PutImageJson() {

  //判断是否存在 image 的数据, 新建或更改 JSON 数据
  beego.Trace("[Image ID] " + this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: this.Ctx.Input.Param(":image_id")}

  has, err := models.Engine.Get(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Select image record error.\"}"))
    this.StopRun()
  }

  beego.Trace("[Has JSON] " + strconv.FormatBool(has))
  beego.Trace("[JSON] " + string(this.Ctx.Input.CopyBody()))
  image.JSON = string(this.Ctx.Input.CopyBody())

  if has == true {
    _, err = models.Engine.Id(image.Id).Cols("JSON").Update(image)
    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(400)
      this.Ctx.Output.Context.Output.Body([]byte("\"Update the image JSON data error.\""))
      this.StopRun()
    }
  }

  _, err = models.Engine.Insert(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Create the image record error.\""))
    this.StopRun()
  }

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))

}

//向本地硬盘写入 Layer 的文件
func (this *ImageController) PutImageLayer() {

  //查询是否存在 image 的数据库记录
  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err := models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("\"Image not found\""))
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.StopRun()
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
    this.StopRun()
  }

  //更新Image记录
  image.Uploaded = true
  _, err = models.Engine.Id(image.Id).Cols("Uploaded").Update(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Update the image upload status error.\"}"))
    this.StopRun()
  }

  //成功则返回 200
  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))
}

func (this *ImageController) PutChecksum() {

  beego.Trace("Cookie: " + this.Ctx.Input.Header("Cookie"))
  beego.Trace("X-Docker-Checksum: " + this.Ctx.Input.Header("X-Docker-Checksum"))
  beego.Trace("X-Docker-Checksum-Payload: " + this.Ctx.Input.Header("X-Docker-Checksum-Payload"))

  //将 checksum 的值保存到数据库
  //X-Docker-Checksum: tarsum+sha256:6eb9bea3d03c72ec2f652869475e21bc11c0409d412c22ea5c44f371d02dda0b
  //X-Docker-Checksum-Payload: sha256:ee40ce84e6e086b23a7d84c8de34ee4b72c82da0327fee85df93cc844a2c9fc3

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId}
  has, err := models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("\"Image not found\""))
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.StopRun()
  }

  //TODO 检查上传的 Layer 文件的 SHA256 和传上来的 Checksum 的值是否一致？
  image.CheckSumed = true
  _, err = models.Engine.Id(image.Id).Cols("CheckSumed").Update(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Update the image checksum error.\""))
    this.StopRun()
  }

  //计算这个Layer的父子结构
  var imageJSON map[string]interface{}
  if err := json.Unmarshal([]byte(image.JSON), &imageJSON); err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Decode the image json data encouter error.\""))
    this.StopRun()
  }

  var parents []string

  //存在 parent 的 ID
  if value, has := imageJSON["parent"]; has {
    parentImage := &models.Image{ImageId: value.(string)}
    has, err := models.Engine.Get(parentImage)
    if err != nil {
      this.Ctx.Output.Context.Output.SetStatus(400)
      this.Ctx.Output.Context.Output.Body([]byte("\"Check the parent image error.\""))
      this.StopRun()
    }

    if has {
      if err := json.Unmarshal([]byte(parentImage.ParentJSON), &parents); err != nil {
        this.Ctx.Output.Context.Output.SetStatus(400)
        this.Ctx.Output.Context.Output.Body([]byte("\"Decode the parent image json data encouter error.\""))
        this.StopRun()
      }
    }
  }

  var images []string
  images = append(images, imageId)
  parents = append(images, parents...)

  parentJSON, _ := json.Marshal(parents)
  image.ParentJSON = string(parentJSON)

  image.Checksum = this.Ctx.Input.Header("X-Docker-Checksum")
  image.Payload = this.Ctx.Input.Header("X-Docker-Checksum-Payload")

  _, err = models.Engine.Update(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Update the image checksum error.\""))
    this.StopRun()
  }

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.Output.SetStatus(200)
  this.Ctx.Output.Context.Output.Body([]byte(""))

}

func (this *ImageController) GetImageAncestry() {

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId, Uploaded: true, CheckSumed: true}
  has, err := models.Engine.Get(image)
  if err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    this.StopRun()
  }

  if has {
    beego.Trace("[Image Ancestry] " + image.ParentJSON)

    this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
    this.Ctx.Output.Context.Output.SetStatus(200)
    this.Ctx.Output.Context.Output.Body([]byte(image.ParentJSON))
  }
}

func (this *ImageController) GetImageLayer() {

  imageId := string(this.Ctx.Input.Param(":image_id"))
  image := &models.Image{ImageId: imageId, Uploaded: true, CheckSumed: true}
  has, err := models.Engine.Get(image)
  if has == false || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(400)
    this.Ctx.Output.Context.Output.Body([]byte("\"Check the image error.\""))
    this.StopRun()
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
