package controllers

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "strings"

  "github.com/astaxie/beego"
  "github.com/bitly/go-simplejson"
  "github.com/dockboard/docker-registry/models"
  "github.com/dockboard/docker-registry/utils"
)

type ImageController struct {
  beego.Controller
}

func (this *ImageController) Prepare() {
}

func (this *ImageController) GetImageJson() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))

  //判断用户的token是否可以操作
  mToken := string(this.Ctx.Input.Header("Authorization"))
  mToken = strings.TrimSpace(mToken)
  mToken = utils.Substr(mToken, 6, len(mToken))

  mUserName := strings.Split(mToken, ",")[1]
  mUserName = strings.Split(mUserName, "=")[1]
  mUserName = utils.Substr(mUserName, 1, len(mUserName)-1)
  mUserName = strings.Split(mUserName, "/")[0]

  mRegistryUser, err := models.GetRegistryUserByToken(mUserName, mToken)

  if mRegistryUser == nil || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
    return
  }
  //查找是否有指定的ImageID对应的JSON文件
  imageId := string(this.Ctx.Input.Param(":image_id"))
  returnImage, err := models.GetImageById(imageId)

  if err != nil || returnImage == nil {
    this.Ctx.Output.Context.Output.SetStatus(404)
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    return
  } else {
    this.Ctx.Output.Context.Output.Body([]byte(returnImage.JSON))
    return
  }
}
func (this *ImageController) PutImageJson() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))
  //判断用户的token是否可以操作
  mToken := string(this.Ctx.Input.Header("Authorization"))
  mToken = strings.TrimSpace(mToken)
  mToken = utils.Substr(mToken, 6, len(mToken))

  mUserName := strings.Split(mToken, ",")[1]
  mUserName = strings.Split(mUserName, "=")[1]
  mUserName = utils.Substr(mUserName, 1, len(mUserName)-1)
  mUserName = strings.Split(mUserName, "/")[0]

  mRegistryUser, err := models.GetRegistryUserByToken(mUserName, mToken)

  if mRegistryUser == nil || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
    return
  }

  putRegistryImage := new(models.Image)

  strImageId := string(this.Ctx.Input.Param(":image_id"))
  putRegistryImage.ImageId = strImageId

  putRegistryImage.JSON = string(this.Ctx.Input.CopyBody())
  putRegistryImage.CheckSumed = false
  putRegistryImage.Uploaded = false

  _, errInsertOneImage := models.InsertOneImage(putRegistryImage)
  if errInsertOneImage != nil {
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
  } else {
    this.Ctx.Output.Context.Output.Body([]byte("true"))
  }
}

func (this *ImageController) PutImageIdLayer() {
  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))
  //判断用户的token是否可以操作
  mToken := string(this.Ctx.Input.Header("Authorization"))
  mToken = strings.TrimSpace(mToken)
  mToken = utils.Substr(mToken, 6, len(mToken))

  mUserName := strings.Split(mToken, ",")[1]
  mUserName = strings.Split(mUserName, "=")[1]
  mUserName = utils.Substr(mUserName, 1, len(mUserName)-1)
  mUserName = strings.Split(mUserName, "/")[0]

  mRegistryUser, err := models.GetRegistryUserByToken(mUserName, mToken)

  if mRegistryUser == nil || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
    return
  }

  //添加指定的ImageID对应的Layer文件
  dockerRegistryBasePath := utils.Cfg.MustValue("docker", "DockerRegistryBasePath")
  strImageId := string(this.Ctx.Input.Param(":image_id"))
  dockerImageJsonSavePath := fmt.Sprintf("%v/images/%v", dockerRegistryBasePath, strImageId)
  //处理目录
  if !utils.IsDirExists(dockerImageJsonSavePath) {
    os.MkdirAll(dockerImageJsonSavePath, os.ModePerm)
  }
  dockerImageJsonSavePath = fmt.Sprintf("%v/images/%v/layer", dockerRegistryBasePath, strImageId)
  if _, err := os.Stat(dockerImageJsonSavePath); err == nil {
    os.Remove(dockerImageJsonSavePath)
  }
  //写入文件
  wfErr := ioutil.WriteFile(dockerImageJsonSavePath, this.Ctx.Input.CopyBody(), 0777)
  //查看文件写入是否成功
  if wfErr != nil {
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }
  //更新Image记录
  imageId := string(this.Ctx.Input.Param(":image_id"))
  returnImage, err := models.GetImageById(imageId)
  if err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }
  returnImage.Uploaded = true
  _, err = models.UpOneImage(returnImage)
  //判断Image记录是否更新成功
  if err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }
  //成功则返回true
  this.Ctx.Output.Context.Output.Body([]byte("true"))
}

func (this *ImageController) PutChecksum() {

  this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
  this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))
  //判断用户的token是否可以操作
  mToken := string(this.Ctx.Input.Header("Authorization"))
  mToken = strings.TrimSpace(mToken)
  mToken = utils.Substr(mToken, 6, len(mToken))

  mUserName := strings.Split(mToken, ",")[1]
  mUserName = strings.Split(mUserName, "=")[1]
  mUserName = utils.Substr(mUserName, 1, len(mUserName)-1)
  mUserName = strings.Split(mUserName, "/")[0]

  mRegistryUser, err := models.GetRegistryUserByToken(mUserName, mToken)

  if mRegistryUser == nil || err != nil {
    this.Ctx.Output.Context.Output.SetStatus(401)
    this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
    return
  }

  //记录checksum的值
  //将checksum的值保存到数据库
  //returnImage.ImageCheckSumed = 1
  //Cookie: session=sFu7ZQLtC0EJPjH693JqWp61jL4=?checksum=KGxwMApTJ3NoYTI1NjplZTQwY2U4NGU2ZTA4NmIyM2E3ZDg0YzhkZTM0ZWU0YjcyYzgyZGEwMzI3ZmVlODVkZjkzY2M4NDRhMmM5ZmMzJwpwMQphLg==
  //X-Docker-Checksum: tarsum+sha256:6eb9bea3d03c72ec2f652869475e21bc11c0409d412c22ea5c44f371d02dda0b
  //X-Docker-Checksum-Payload: sha256:ee40ce84e6e086b23a7d84c8de34ee4b72c82da0327fee85df93cc844a2c9fc3

  imageId := string(this.Ctx.Input.Param(":image_id"))
  returnImage, err := models.GetImageById(imageId)
  returnImage.Checksum = this.Ctx.Input.Header("X-Docker-Checksum")
  returnImage.Payload = this.Ctx.Input.Header("X-Docker-Checksum-Payload")
  returnImage.CheckSumed = true
  _, err = models.UpOneImage(returnImage)
  if err != nil {
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }

  //计算这个Layer的父子结构
  infoJson, infoJsonErr := simplejson.NewJson([]byte(returnImage.JSON))
  if infoJsonErr != nil {
    panic("json format error")
    this.Ctx.Output.Context.Output.Body([]byte("{\"error\": \"Image not found\"}"))
    this.Ctx.Output.Context.Output.SetStatus(404)
    return
  }
  parentJson := []string{returnImage.ImageId}
  // 检查key是否存在
  parentIdJson, ok := infoJson.CheckGet("parent")
  if ok {
    //取得父类的parent相关json
    parentId, _ := parentIdJson.String()
    parentImage, _ := models.GetImageById(parentId)
    parentImageArrayJSON, _ := simplejson.NewJson([]byte(parentImage.ParentJSON))

    parentImageArray, _ := parentImageArrayJSON.StringArray()
    parentJson = append(parentJson, parentImageArray...)

  }
  jsonByte, jsonByteErr := json.Marshal(parentJson)
  fmt.Println("准备更新ImageParentJSON")

  //更新Image数据
  if jsonByteErr == nil {
    returnImage.ParentJSON = string(jsonByte)
    fmt.Println("ImageParentJson:", returnImage.ParentJSON)
    models.UpOneImage(returnImage)
  }

  this.Ctx.Output.Context.Output.Body([]byte("true"))

}

func (this *ImageController) GETLayer() {

}

func (this *ImageController) PUTLayer() {

}

func (this *ImageController) GETJSON() {

}

func (this *ImageController) PUTJSON() {

}

func (this *ImageController) GETAncestry() {

}

func (this *ImageController) PUTChecksum() {

}

func (this *ImageController) GETFiles() {

}

func (this *ImageController) GETDiff() {

}

func (this *ImageController) GETPrivateLayer() {

}

func (this *ImageController) GETPrivateJSON() {

}

func (this *ImageController) GETPrivateFiles() {

}
