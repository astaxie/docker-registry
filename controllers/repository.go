// Docker Push & Pull
// 执行 docker push 命令流程：
//     1. docker 向 registry 服务器注册 repository： PUT /v1/repositories/<username>/<repository> -> PUTRepository()
//     2. 参数是 JSON 格式的 <repository> 所有 image 的 id 列表，按照 image 的构建顺序排列。
//     3. 根据 <repository> 的 <tags> 进行循环：
//        3.1 获取 <image> 的 JSON 文件：GET /v1/images/<image_id>/json -> image.go#GETJSON()
//        3.2 如果没有此文件或内容返回 404 。
//        3.3 docker push 认为服务器没有 image 对应的文件，向服务器上传 image 相关文件。
//            3.3.1 写入 <image> 的 JSON 文件：PUT /v1/images/<image_id>/json -> image.go#PUTJSON()
//            3.3.2 写入 <image> 的 layer 文件：PUT /v1/images/<image_id>/layer -> image.go#PUTLayer()
//            3.3.3 写入 <image> 的 checksum 信息：PUT /v1/images/<image_id>/checksum -> image.go#PUTChecksum()
//        3.4 上传完此 tag 的所有 image 后，向服务器写入 tag 信息：PUT /v1/repositories/(namespace)/(repository)/tags/(tag) -> PUTTag()
//     4. 所有 tags 的 image 上传完成后，向服务器发送所有 images 的校验信息，PUT /v1/repositories/(namespace)/(repo_name)/images -> PUTRepositoryImages()
package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker-registry/auth"
	"github.com/docker-registry/models"
	"github.com/docker-registry/utils"
	"github.com/nu7hatch/gouuid"
	"os"
	"time"
)

type RepositoryController struct {
	beego.Controller
}

func (this *RepositoryController) Prepare() {
	//TODO Generate token & read endpoints from beego.AppConfig.String("Endpoints")
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Endpoints", "")
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("WWW-Authenticate", "")
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Token", "")
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

	//tokenSignature, err := uuid.NewV5(uuid.NamespaceURL, []byte(this.Ctx.Input.Url()))
	//if err != nil {
	//fmt.Println("error:", err)
	//return
	//}

	//PUT /v1/repositories/dockerfile/redis/ HTTP/1.1
	//Host: 192.168.168.86:5000
	//User-Agent: docker/0.10.0 go/go1.2.1 git-commit/dc9c28f kernel/3.8.0-38-generic os/linux arch/amd64
	//Content-Length: 1298
	//Authorization: Basic Og==
	//Content-Type: application/json
	//X-Docker-Token: true
	//Accept-Encoding: gzip
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	//this.Ctx.Output.Context.ResponseWriter.Header().Set("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))

	//this.Ctx.Output.ContentType("application/json")
	//this.Ctx.Output.Header("X-Docker-Token", xDockerToken)
	//this.Ctx.Output.Header("WWW-Authenticate", xDockerToken)
	//this.Ctx.Output.Header("X-Docker-Endpoints", xDockerEndpoints)
	//this.Ctx.Output.Header("X-Docker-Registry-Version", utils.Cfg.MustValue("docker", "XDockerRegistryVersion"))
	//this.Ctx.Output.Header("X-Docker-Registry-Config", utils.Cfg.MustValue("docker", "XDockerRegistryConfig"))

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
	//这里应该计算checksum
	//返回204
	this.Ctx.Output.Context.Output.SetStatus(204)
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
// Create a user repository with the given namespace and repo_name.
// Example Request:
//    PUT /v1/repositories/foo/bar/ HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Authorization: Basic akmklmasadalkm==
//    X-Docker-Token: true
//    [{"id": "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f"}]
// Parameters:
//    namespace – the namespace for the repo
//    repo_name – the name for the repo
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    WWW-Authenticate: Token signature=123abc,repository="foo/bar",access=write
//    X-Docker-Token: signature=123abc,repository="foo/bar",access=write
//    X-Docker-Endpoints: registry-1.docker.io [, registry-2.docker.io]
//    ""
// Status Codes:
//    200 – Created
//    400 – Errors (invalid json, missing or invalid fields, etc)
//    401 – Unauthorized
//    403 – Account is not Active
func (this *RepositoryController) PUTRepository() {

}

// http://docs.docker.io/en/latest/reference/api/index_api/#repository
// DELETE /v1/repositories/(namespace)/(repo_name)/
// Delete a user repository with the given namespace and repo_name.
// Example Request:
//    DELETE /v1/repositories/foo/bar/ HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Authorization: Basic akmklmasadalkm==
//    X-Docker-Token: true
//    ""
// Parameters:
//    namespace – the namespace for the repo
//    repo_name – the name for the repo
// Example Response:
//    HTTP/1.1 202
//    Vary: Accept
//    Content-Type: application/json
//    WWW-Authenticate: Token signature=123abc,repository="foo/bar",access=delete
//    X-Docker-Token: signature=123abc,repository="foo/bar",access=delete
//    X-Docker-Endpoints: registry-1.docker.io [, registry-2.docker.io]
//    ""
// Status Codes:
//    200 – Deleted
//    202 – Accepted
//    400 – Errors (invalid json, missing or invalid fields, etc)
//    401 – Unauthorized
//    403 – Account is not Active
func (this *RepositoryController) DELETERepository() {

}

// http://docs.docker.io/reference/api/index_api/#repository-images
// PUT /v1/repositories/(namespace)/(repo_name)/images
// Update the images for a user repo.
// Example Request:
//    PUT /v1/repositories/foo/bar/images HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Content-Type: application/json
//    Authorization: Basic akmklmasadalkm==
//    [
//      {
//        "id": "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f",
//        "checksum": "b486531f9a779a0c17e3ed29dae8f12c4f9e89cc6f0bc3c38722009fe6857087"
//      }
//    ]
// Parameters:
//    namespace – the namespace for the repo
//    repo_name – the name for the repo
// Example Response:
//    HTTP/1.1 204
//    Vary: Accept
//    Content-Type: application/json
//    ""
// Status Codes:
//    204 – Created
//    400 – Errors (invalid json, missing or invalid fields, etc)
//    401 – Unauthorized
//    403 – Account is not Active or permission denied
func (this *RepositoryController) PUTRepositoryImages() {

}

// http://docs.docker.io/reference/api/index_api/#repository-images
// GET /v1/repositories/(namespace)/(repo_name)/images
// get the images for a user repo.
// Example Request:
//    GET /v1/repositories/foo/bar/images HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
// Parameters:
//    namespace – the namespace for the repo
//    repo_name – the name for the repo
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//     [{"id": "9e89cc6f0bc3c38722009fe6857087b486531f9a779a0c17e3ed29dae8f12c4f",
//     "checksum": "b486531f9a779a0c17e3ed29dae8f12c4f9e89cc6f0bc3c38722009fe6857087"},
//     {"id": "ertwetewtwe38722009fe6857087b486531f9a779a0c1dfddgfgsdgdsgds",
//     "checksum": "34t23f23fc17e3ed29dae8f12c4f9e89cc6f0bsdfgfsdgdsgdsgerwgew"}]
// Status Codes:
//    200 – OK
//    404 – Not found
func (this *RepositoryController) GETRepositoryImages() {

}

// http://docs.docker.io/reference/api/index_api/#repository-authorization
// PUT /v1/repositories/(namespace)/(repo_name)/auth
// authorize a token for a user repo
// Example Request:
//    PUT /v1/repositories/foo/bar/auth HTTP/1.1
//    Host: index.docker.io
//    Accept: application/json
//    Authorization: Token signature=123abc,repository="foo/bar",access=write
// Parameters:
//    namespace – the namespace for the repo
//    repo_name – the name for the repo
// Example Response:
//    HTTP/1.1 200
//    Vary: Accept
//    Content-Type: application/json
//    "OK"
// Status Codes:
//    200 – OK
//    403 – Permission denied
//    404 – Not found
func (this *RepositoryController) PUTRepositoryAuth() {

}

// Undocumented API
// PUT /v1/repositories/:username/:repository/properties
func (this *RepositoryController) PUTProperties() {

}

// Undocumented API
// GET /v1/repositories/:username/:repository/properties
func (this *RepositoryController) GETProperties() {

}

// Undocumented API
// GET /v1/repositories/:username/:repository/json
func (this *RepositoryController) GETRepositoryJSON() {

}

// Undocumented API
// GET /v1/repositories/:username/:repository/tags/:tag/json
func (this *RepositoryController) GETTagJSON() {

}

// Undocumented API
// DELETE /v1/repositories/:username/:repository/tags/:tag/json
func (this *RepositoryController) DELETERepositoryTags() {

}
