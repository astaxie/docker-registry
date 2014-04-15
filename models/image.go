package models

import (
	"errors"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"syscall"
	"unicode/utf16"
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

//字符串转换来unit16，移动文件时需要！
func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

//确认有_checksum 有_checksum并且文件长度为71 才代表这个image是完整的
func (this *Image) LayerIsCompletely() bool {
	fileStat, fileError := os.Stat(beego.AppConfig.String("RegistryPath") + this.ImageId + "/_checksum")
	return (fileStat.Size() == 71) && (fileError == nil || os.IsExist(fileError))
}

//读取Image(Model)的JSON数据
func (this *Image) readJsonData() (string, error) {
	fi, err := os.Open(beego.AppConfig.String("RegistryPath") + this.ImageId + "/json")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd), err
}

//写入Image(Model)的JSON数据
func (this *Image) writeJsonData(JsonData string) error {
	jsonFileName := beego.AppConfig.String("RegistryPath") + this.ImageId + "/json"
	jsonSlices := []byte(JsonData)
	err := ioutil.WriteFile(jsonFileName, jsonSlices, os.ModeAppend)
	return err
}

//外部通过此方法，获得Image(Model)的JSON数据
func (this *Image) GetImageJsonDataById() (string, error) {
	//查找是否有指定的ImageId的Layer数据
	if this.LayerIsCompletely() {
		//--有 _checksum 则读入json并返回
		imageJsonData, err := this.readJsonData()
		return imageJsonData, err
	} else {
		//--没有 _checksum 则返回""
		return "", errors.New("GetImageJsonDataById:" + this.ImageId + ", not exist checksum")
	}
}

//外部通过此方法，设置Image(Model)的JSON数据
func (this *Image) SetImageJsonDataById(ImageJsonData string) error {
	return this.writeJsonData(ImageJsonData)
}

//外部通过此方法，获得Image(Model)的Layer文件位置
func (this *Image) GetImageLayerDataPathById() (string, error) {

	layerFilePath := beego.AppConfig.String("RegistryPath") + this.ImageId + "/layer"

	return layerFilePath, nil
}

//外部通过此方法，设置Image(Model)的Layer文件位置(参数必须指向一个文件，这里只是将文件移动到Layer所在位置)
func (this *Image) SetImageLayerDataPathById(ImageLayerDataPath string) error {
	fromPath := ImageLayerDataPath
	toPath := beego.AppConfig.String("RegistryPath") + this.ImageId + "/layer"
	err := syscall.Rename(fromPath, toPath)
	return err
}

func (this *Image) GetImageLayerChecksumById(ImageLayerId string) (string, error) {
	return "", nil
}
func (this *Image) SetImageLayerChecksumById(ImageLayerId string, ImageLayerChecksum string) (string, error) {
	return "", nil
}
