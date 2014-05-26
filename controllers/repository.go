/*
Docker Push & Pull
执行 docker push 命令流程：
    1. docker 向 registry 服务器注册 repository： PUT /v1/repositories/<username>/<repository> -> PUTRepository()
    2. 参数是 JSON 格式的 <repository> 所有 image 的 id 列表，按照 image 的构建顺序排列。
    3. 根据 <repository> 的 <tags> 进行循环：
       3.1 获取 <image> 的 JSON 文件：GET /v1/images/<image_id>/json -> image.go#GETJSON()
       3.2 如果没有此文件或内容返回 404 。
       3.3 docker push 认为服务器没有 image 对应的文件，向服务器上传 image 相关文件。
           3.3.1 写入 <image> 的 JSON 文件：PUT /v1/images/<image_id>/json -> image.go#PUTJSON()
           3.3.2 写入 <image> 的 layer 文件：PUT /v1/images/<image_id>/layer -> image.go#PUTLayer()
           3.3.3 写入 <image> 的 checksum 信息：PUT /v1/images/<image_id>/checksum -> image.go#PUTChecksum()
       3.4 上传完此 tag 的所有 image 后，向服务器写入 tag 信息：PUT /v1/repositories/(namespace)/(repository)/tags/(tag) -> PUTTag()
    4. 所有 tags 的 image 上传完成后，向服务器发送所有 images 的校验信息，PUT /v1/repositories/(namespace)/(repo_name)/images -> PUTRepositoryImages()
*/
package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/models"
	"github.com/dockboard/docker-registry/utils"
)

type RepositoryController struct {
	beego.Controller
}

func (this *RepositoryController) Prepare() {
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "Version"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "Config"))
}

func (this *RepositoryController) PutRepository() {
	//Decode 后根据数据库判断用户是否存在和是否激活。
	beego.Trace("Authorization" + this.Ctx.Input.Header("Authorization"))
	username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))
	beego.Trace("Decode Basic Auth: " + username + " " + passwd)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(401)
		this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
		return
	}

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

	beego.Trace("User:" + user.Username)

	//获取namespace/repository
	namespace := string(this.Ctx.Input.Param(":namespace"))
	repository := string(this.Ctx.Input.Param(":repo_name"))

	//判断用户的username和namespace是否相同
	if username != namespace {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"username != namespace\""))
		return
	}

	//创建token并保存
	//需要加密的字符串为 UserName+UserPassword+时间戳
	md5String := fmt.Sprintf("%v%v%v", username, passwd, string(time.Now().Unix()))
	h := md5.New()
	h.Write([]byte(md5String))
	signature := hex.EncodeToString(h.Sum(nil))
	token := fmt.Sprintf("Token signature=%v,repository=\"%v/%v\",access=write", signature, namespace, repository)

	beego.Trace("Token:" + token)

	//保存Token
	user.Token = token
	_, err = models.Engine.Id(user.Id).Cols("Token").Update(user)

	if err != nil {
		beego.Trace(err)
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Update token error.\""))
		return
	}

	//创建或更新 Repository 数据
	//也可以采用 ioutil.ReadAll(this.Ctx.Request.Body) 的方式读取 body 数据
	beego.Trace("Request Body: " + string(this.Ctx.Input.CopyBody()))

	repo := &models.Repository{Namespace: namespace, Repository: repository}
	has, err = models.Engine.Get(repo)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Search repository error.\""))
		return
	}

	repo.JSON = string(this.Ctx.Input.CopyBody())

	if has == true {
		_, err := models.Engine.Id(repo.Id).Cols("JSON").Update(repo)
		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Update the repository JSON data error.\""))
			return
		}
	} else {
		_, err := models.Engine.Insert(repo)
		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Create the repository record error: \"" + err.Error()))
			return
		}
	}

	//操作正常的输出
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Token", token)
	this.Ctx.Output.Context.ResponseWriter.Header().Set("WWW-Authenticate", token)
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Endpoints", utils.Cfg.MustValue("docker", "Endpoints"))

	this.Ctx.Output.Context.Output.SetStatus(200)
	this.Ctx.Output.Context.Output.Body([]byte("\"\""))
}

func (this *RepositoryController) PutTag() {

	beego.Trace("Namespace: " + this.Ctx.Input.Param(":namespace"))
	beego.Trace("Repository: " + this.Ctx.Input.Param(":repository"))
	beego.Trace("Tag: " + this.Ctx.Input.Param(":tag"))
	beego.Trace("User-Agent: " + this.Ctx.Input.Header("User-Agent"))

	repository := &models.Repository{Namespace: this.Ctx.Input.Param(":namespace"), Repository: this.Ctx.Input.Param(":repository")}
	has, err := models.Engine.Get(repository)

	if has == false || err != nil {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Unknow namespace and repository.\""))
		return
	}

	tag := &models.Tag{Name: this.Ctx.Input.Param(":tag"), Repository: repository.Id}
	has, err = models.Engine.Get(tag)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Search tag encounter error.\""))
		return
	}

	tag.JSON = this.Ctx.Input.Header("User-Agent")

	r, _ := regexp.Compile(`"([[:alnum:]]+)"`)
	imageIds := r.FindStringSubmatch(string(this.Ctx.Input.CopyBody()))

	tag.ImageId = imageIds[1]

	if has == true {
		_, err := models.Engine.Id(tag.Id).Update(tag)
		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Update the tag data error.\""))
			return
		}
	} else {
		_, err := models.Engine.Insert(tag)
		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Create the tag record error.\""))
			return
		}
	}

	//操作正常的输出
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")

	this.Ctx.Output.Context.Output.SetStatus(200)
	this.Ctx.Output.Context.Output.Body([]byte("\"\""))
}

//docker client 没有上传完整的 checksum ，如何进行完整性检查？
func (this *RepositoryController) PutRepositoryImages() {
	//操作正常的输出
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")

	this.Ctx.Output.Context.Output.SetStatus(204)
	this.Ctx.Output.Context.Output.Body([]byte("\"\""))
}

func (this *RepositoryController) GetRepositoryImages() {
	//Decode 后根据数据库判断用户是否存在和是否激活。
	beego.Trace("Authorization" + this.Ctx.Input.Header("Authorization"))
	username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))
	beego.Trace("Decode Basic Auth: " + username + " " + passwd)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(401)
		this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
		return
	}

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

	beego.Trace("User:" + user.Username)

	//获取namespace/repository
	namespace := string(this.Ctx.Input.Param(":namespace"))
	repository := string(this.Ctx.Input.Param(":repo_name"))

	//创建token并保存
	//需要加密的字符串为 UserName+UserPassword+时间戳
	md5String := fmt.Sprintf("%v%v%v", username, passwd, string(time.Now().Unix()))
	h := md5.New()
	h.Write([]byte(md5String))
	signature := hex.EncodeToString(h.Sum(nil))
	token := fmt.Sprintf("Token signature=%v,repository=\"%v/%v\",access=write", signature, namespace, repository)

	beego.Trace("Token:" + token)

	//保存Token
	user.Token = token
	_, err = models.Engine.Id(user.Id).Cols("Token").Update(user)

	if err != nil {
		beego.Trace(err)
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Update token error.\""))
		return
	}

	//查询 Repository 数据
	repo := &models.Repository{Namespace: namespace, Repository: repository}
	has, err = models.Engine.Get(repo)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Search repository error.\""))
		return
	}

	if has == false {
		this.Ctx.Output.Context.Output.SetStatus(404)
		this.Ctx.Output.Context.Output.Body([]byte("\"Cloud not found repository.\""))
		return
	} else {
		//存在 Repository 数据，查询所有的 Tag 数据。
		tags := make([]models.Tag, 0)
		err := models.Engine.Where("repository_id= ?", repo.Id).Find(&tags)

		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Search repository tag error.\""))
			return
		}

		if len(tags) == 0 {
			this.Ctx.Output.Context.Output.SetStatus(404)
			this.Ctx.Output.Context.Output.Body([]byte("\"Cloud not found any tag.\""))
			return
		}

		//根据 Tag 的 Image ID 值查询 ParentJSON 数据，然后同一在一个数组里面去重。
		var images []string
		for _, tag := range tags {
			image := &models.Image{ImageId: tag.ImageId}
			has, err := models.Engine.Get(image)

			if has == false || err != nil {
				this.Ctx.Output.Context.Output.SetStatus(400)
				this.Ctx.Output.Context.Output.Body([]byte("\"Search image error.\""))
				return
			}

			if has == true {
				var parents []string

				beego.Trace(string(image.Id) + ":\n" + image.ParentJSON)

				if err := json.Unmarshal([]byte(image.ParentJSON), &parents); err != nil {
					this.Ctx.Output.Context.Output.SetStatus(400)
					this.Ctx.Output.Context.Output.Body([]byte("\"Decode the parent image json data encouter error.\""))
					return
				}
				images = append(parents, images...)
			}
		}

		utils.RemoveDuplicateString(&images)

		//转换为 map 的对象返回
		var results []map[string]string
		for _, value := range images {
			result := make(map[string]string)
			result["id"] = value
			results = append(results, result)
		}

		imageIds, _ := json.Marshal(results)

		beego.Trace("Image ID:" + string(imageIds))

		//操作正常的输出
		this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
		this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Token", token)
		this.Ctx.Output.Context.ResponseWriter.Header().Set("WWW-Authenticate", token)
		this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Endpoints", utils.Cfg.MustValue("docker", "Endpoints"))

		this.Ctx.Output.Context.Output.SetStatus(200)
		this.Ctx.Output.Context.Output.Body(imageIds)

	}
}

func (this *RepositoryController) GetRepositoryTags() {
	//Decode 后根据数据库判断用户是否存在和是否激活。
	beego.Trace("Authorization" + this.Ctx.Input.Header("Authorization"))
	username, passwd, err := utils.DecodeBasicAuth(this.Ctx.Input.Header("Authorization"))
	beego.Trace("Decode Basic Auth: " + username + " " + passwd)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(401)
		this.Ctx.Output.Context.Output.Body([]byte("\"Unauthorized\""))
		return
	}

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

	beego.Trace("User:" + user.Username)

	//获取namespace/repository
	namespace := string(this.Ctx.Input.Param(":namespace"))
	repository := string(this.Ctx.Input.Param(":repo_name"))

	//查询 Repository 数据
	repo := &models.Repository{Namespace: namespace, Repository: repository}
	has, err = models.Engine.Get(repo)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(400)
		this.Ctx.Output.Context.Output.Body([]byte("\"Search repository error.\""))
		return
	}

	if has == false {
		this.Ctx.Output.Context.Output.SetStatus(404)
		this.Ctx.Output.Context.Output.Body([]byte("\"Cloud not found repository.\""))
		return
	} else {
		//存在 Repository 数据，查询所有的 Tag 数据。
		tags := make([]models.Tag, 0)
		err := models.Engine.Where("repository_id= ?", repo.Id).Find(&tags)

		if err != nil {
			this.Ctx.Output.Context.Output.SetStatus(400)
			this.Ctx.Output.Context.Output.Body([]byte("\"Search repository tag error.\""))
			return
		}

		if len(tags) == 0 {
			this.Ctx.Output.Context.Output.SetStatus(404)
			this.Ctx.Output.Context.Output.Body([]byte("\"Cloud not found any tag.\""))
			return
		}

		results := make([]interface{}, 0)
		for _, v := range tags {
			tag := make(map[string]string, 0)
			tag[v.Name] = v.ImageId
			results = append(results, tag)
		}

		result, _ := json.Marshal(results)

		//操作正常的输出
		this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")

		this.Ctx.Output.Context.Output.SetStatus(200)
		this.Ctx.Output.Context.Output.Body(result)
	}
}
