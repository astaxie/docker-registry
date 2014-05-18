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
  "github.com/astaxie/beego"
)

type RepositoryController struct {
  beego.Controller
}

func (this *RepositoryController) Prepare() {

}

func (this *RepositoryController) GETTags() {

}

func (this *RepositoryController) GETTag() {

}

func (this *RepositoryController) DELETETag() {

}

func (this *RepositoryController) PUTTag() {

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
