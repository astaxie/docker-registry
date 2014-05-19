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
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/auth"
	"github.com/dockboard/docker-registry/models"
	"github.com/dockboard/docker-registry/utils"
	"github.com/nu7hatch/gouuid"
)

type RepositoryController struct {
	beego.Controller
}

func (this *RepositoryController) Prepare() {
}

func (this *RepositoryController) PutNamespaceRepo() {
	fmt.Println("进入PutNamespaceRepo")
	//判断用户是否合法
	authorizationBasic := this.Ctx.Input.Header("Authorization")
	authUsername, authPasswd, authErr := auth.BaseAuth(authorizationBasic)
	if authErr != nil {
		this.Ctx.Output.Context.Output.SetStatus(401)
		this.Ctx.Output.Context.Output.Body([]byte("{\"Unauthorized\"}"))
		fmt.Println("authErr")
		return
	}

	//Content-Type: application/json
	//X-Docker-Registry-Version: 0.6.8
	//X-Docker-Registry-Config: dev
	this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))

	dockerRegistryBasePath := utils.Cfg.MustValue("docker", "DockerRegistryBasePath")
	xDockerEndpoints := utils.Cfg.MustValue("docker", "XDockerEndpoints")
	strNamespace := string(this.Ctx.Input.Param(":namespace"))
	strRepoName := string(this.Ctx.Input.Param(":repo_name"))
	dockerRegistryRepoPath := fmt.Sprintf("%v/repositories/%v/%v", dockerRegistryBasePath, strNamespace, strRepoName)

	if authUsername != strNamespace {
		this.Ctx.Output.Context.Output.SetStatus(401)
		this.Ctx.Output.Context.Output.Body([]byte("{\"username != namespace\"}"))
		fmt.Println(authUsername, ":authUsername != strNamespace:", strNamespace)

		return
	}

	//判断目录是否存在，不存在则创建对应目录
	if !utils.IsDirExists(dockerRegistryRepoPath) {
		os.MkdirAll(dockerRegistryRepoPath, os.ModePerm)
	}
	//创建token并保存
	// 需要加密的字符串为 UserName+UserPassword+时间戳
	md5String := fmt.Sprintf("%v%v%v", authUsername, authPasswd, string(time.Now().Unix()))
	h := md5.New()
	h.Write([]byte(md5String))
	tokenSignature := hex.EncodeToString(h.Sum(nil))
	xDockerToken := fmt.Sprintf("Token signature=%v,repository=\"%v/%v\",access=write",
		tokenSignature, strNamespace, strRepoName)
	mRegistryUser, err := models.GetRegistryUserByUserName(authUsername)
	if err != nil {
		this.Ctx.Output.Context.Output.SetStatus(417)
		this.Ctx.Output.Context.Output.Body([]byte("{\"Expectation failed\"}"))

		return
	}
	mRegistryUser.UserToken = xDockerToken
	models.UpRegistryUser(mRegistryUser)

	//X-Docker-Token: Token signature=OT63AV22Y5CGZV7N,repository="dockerfile/redis",access=write
	//WWW-Authenticate: Token signature=OT63AV22Y5CGZV7N,repository="dockerfile/redis",access=write
	//X-Docker-Endpoints: 192.168.1.132:5000
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Token", xDockerToken)
	this.Ctx.Output.Context.ResponseWriter.Header().Set("WWW-Authenticate", xDockerToken)
	this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Endpoints", xDockerEndpoints)

	this.Ctx.Output.Context.Output.Body([]byte("\"\""))
}

func (this *RepositoryController) PutRepo() {

	//判断目录是否存在，不存在则创建对应目录
	dockerRegistryBasePath := utils.Cfg.MustValue("docker", "DockerRegistryBasePath")
	dockerRegistryRepoPath := fmt.Sprintf("%v/repositories/library/%v", dockerRegistryBasePath, string(this.Ctx.Input.Param(":repo_name")))
	if !utils.IsDirExists(dockerRegistryRepoPath) {
		os.MkdirAll(dockerRegistryRepoPath, os.ModePerm)
	}

	//返回结果处理
	//X-Docker-Token: Token signature=NU66YCHG8FK63I3G,repository="library/redis",access=write
	tokenSignature, err := uuid.NewV5(uuid.NamespaceURL, []byte(this.Ctx.Input.Url()))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	xDockerToken := fmt.Sprintf("Token signature=%v,repository=\"library/%v\",access=write",
		tokenSignature, string(this.Ctx.Input.Param(":repo_name")))
	xDockerEndpoints := utils.Cfg.MustValue("docker", "XDockerEndpoints")

	this.Ctx.Output.ContentType("application/json")
	this.Ctx.Output.Header("X-Docker-Token", xDockerToken)
	this.Ctx.Output.Header("WWW-Authenticate", xDockerToken)
	this.Ctx.Output.Header("X-Docker-Endpoints", xDockerEndpoints)
	this.Ctx.Output.Header("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	this.Ctx.Output.Header("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))

	this.Ctx.Output.Body([]byte("\"\""))

}

func (this *RepositoryController) PutNamespaceTag() {
	//保存Tag信息
	nowPutTag := new(models.RegistryRepositorieTag)

	nowPutTag.RepositorieTag = string(this.Ctx.Input.CopyBody())
	nowPutTag.RepositorieTagName = this.Ctx.Input.Param(":tag")
	nowPutTag.RepositorieTagJson = this.Ctx.Input.Header("User-Agent")
	nowPutTag.RepositorieTagNamespace = this.Ctx.Input.Param(":namespace")
	nowPutTag.RepositorieTagRepository = this.Ctx.Input.Param(":repository")
	models.PutOneTag(nowPutTag)
}

func (this *RepositoryController) PutTag() {
	//保存Tag信息
	nowPutTag := new(models.RegistryRepositorieTag)

	nowPutTag.RepositorieTag = string(this.Ctx.Input.CopyBody())
	nowPutTag.RepositorieTagName = this.Ctx.Input.Param(":tag")
	nowPutTag.RepositorieTagJson = this.Ctx.Input.Header("User-Agent")
	nowPutTag.RepositorieTagNamespace = "library"
	nowPutTag.RepositorieTagRepository = this.Ctx.Input.Param(":repository")
	models.PutOneTag(nowPutTag)

}

func (this *RepositoryController) PutNamespaceImages() {
	//这里应该计算checksum
	//返回204
	this.Ctx.Output.Context.Output.SetStatus(204)

}

func (this *RepositoryController) PutImages() {
	this.Ctx.Output.Context.Output.SetStatus(204)
}

func (this *RepositoryController) GETTags() {

}

func (this *RepositoryController) GETTag() {

}

func (this *RepositoryController) DELETETag() {

}

func (this *RepositoryController) DELETERepositoryImages() {

}

func (this *RepositoryController) PUTRepository() {

}

func (this *RepositoryController) DELETERepository() {

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

func (this *RepositoryController) GETRepositoryJSON() {

}

func (this *RepositoryController) GETTagJSON() {

}

func (this *RepositoryController) DELETERepositoryTags() {

}
