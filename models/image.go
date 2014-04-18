package models

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"github.com/astaxie/beego"
	"github.com/dotcloud/docker/utils"
	"io/ioutil"
	"os"
	"syscall"
)

type Image struct {
	ImageId           string
	ImageJsonData     string
	ImageJsonLayerUri string
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

//写入checksum的数据
func (this *Image) writeChecksumData(ChecksumData string) error {
	checksumFileName := beego.AppConfig.String("RegistryPath") + this.ImageId + "/_checksum"
	checksumSlices := []byte(ChecksumData)
	err := ioutil.WriteFile(checksumFileName, checksumSlices, os.ModeAppend)
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

//外部通过此方法，获取Image(Model)Layer的checksum
func (this *Image) GetImageLayerChecksumById() (string, error) {
	fi, err := os.Open(beego.AppConfig.String("RegistryPath") + this.ImageId + "/_checksum")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd), err
}

//外部通过此方法，让Image(Model)计算Layer的checksum并对比传出的ImageLayerChecksum
func (this *Image) SetImageLayerChecksumById(ImageLayerPayloadChecksum string) (Checksum string, ChecksumPayload string, err error) {

	layerFilePath := beego.AppConfig.String("RegistryPath") + this.ImageId + "/layer"
	var fileBufioReader *bufio.Reader

	layerFile, layerFileErr := os.Open(layerFilePath)
	if layerFileErr != nil {
		return "", "", layerFileErr
	} else {
		fileBufioReader = bufio.NewReader(layerFile)
		tarsumLayer := &utils.TarSum{Reader: fileBufioReader}

		h := sha256.New()
		imageJsonString, _ := this.GetImageJsonDataById()
		h.Write([]byte(imageJsonString))
		h.Write([]byte{'\n'})

		checksumLayer := &utils.CheckSum{Reader: fileBufioReader, Hash: h}
		checksumPayload := "sha256:" + checksumLayer.Sum()
		if checksumPayload != ImageLayerPayloadChecksum {
			return tarsumLayer.Sum([]byte(imageJsonString)), checksumPayload, errors.New("X-Docker-Checksum-Payload not consistent")
		} else {
			this.writeChecksumData(checksumPayload)
			return tarsumLayer.Sum([]byte(imageJsonString)), checksumPayload, nil
		}
	}
}
