package models

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
)

type Image struct {
	ImageId           string
	ImageJsonData     string
	ImageJsonLayerUri string
}

//{
//	"comment": "Imported from http://get.docker.io/images/base",
//	"container_config": {
//		"Tty": false,
//		"Cmd": null,
//		"MemorySwap": 0,
//		"Image": "",
//		"Hostname": "",
//		"User": "",
//		"Env": null,
//		"Memory": 0,
//		"Detach": false,
//		"Ports": null,
//		"OpenStdin": false
//	},
//	"id": "27cf784147099545",
//	"created": "2013-03-23T12:53:11.10432-07:00"
//}

//确认有_checksum 有_checksum并且文件长度为71 才代表这个image是完整的
func (this *Image) LayerIsCompletely() bool {
	fileStat, fileError := os.Stat(beego.AppConfig.String("RegistryPath") + this.ImageId + "/_checksum")
	return (fileStat.Size() == 71) && (fileError == nil || os.IsExist(fileError))
}

func readJsonFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func (this *Image) GetImageJsonDataById() (ReturnImageJson string, ErrorInfo string) {
	//查找是否有指定的ImageId的Layer数据
	if this.LayerIsCompletely() {
		//--有 _checksum 则读入json并返回
		imageJsonData := readJsonFile(beego.AppConfig.String("RegistryPath") + this.ImageId + "/json")
		return imageJsonData, ""
	} else {
		//--没有 _checksum 则返回""
		return "", "no checksum"
	}

}

func (this *Image) SetImageJsonDataById(ImageJsonData string) (ReturnImageJson string, ErrorInfo string) {
	return "", ""
}

func (this *Image) GetImageLayerDataById(ImageLayerId string) (ReturnImageLayerJson string, ErrorInfo string) {
	return "", ""
}

func (this *Image) SetImageLayerDataById(ImageLayerId string, ImageLayerJson string) (ReturnImageLayerJson string, ErrorInfo string) {
	return "", ""
}

func (this *Image) GetImageLayerChecksumById(ImageLayerId string) (ReturnImageLayerJson string, ErrorInfo string) {
	return "", ""
}
func (this *Image) SetImageLayerChecksumById(ImageLayerId string, ImageLayerChecksum string) (ReturnImageLayerJson string, ErrorInfo string) {
	return "", ""
}
