package models

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/dockboard/docker-registry/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var x *xorm.Engine

type RegistryUser struct {
	Id           int64
	UserName     string `xorm:"VARCHAR(255)"`
	UserPassword string `xorm:"VARCHAR(255)"`
	UserEmail    string `xorm:"VARCHAR(255)"`
	UserToken    string `xorm:"VARCHAR(255)"`
	//token MD5(UserName+UserPassword+时间戳)
}

type RegistryImage struct {
	Id                     int64
	ImageId                string `xorm:"VARCHAR(255)"`
	ImageJson              string `xorm:"TEXT"`
	ImageParentJson        string `xorm:"TEXT"`
	ImageUploaded          int64
	ImageCheckSumed        int64
	XDockerChecksum        string `xorm:"TEXT"`
	XDockerChecksumPayload string `xorm:"TEXT"`
}

//type RegistryRepositorie struct {
//	Id              int64
//	RepositorieName string `xorm:"VARCHAR(255)"`
//	RepositorieJson string `xorm:"TEXT"`
//}

type RegistryRepositorieTag struct {
	Id                       int64
	RepositorieTagNamespace  string `xorm:"VARCHAR(255)"`
	RepositorieTagRepository string `xorm:"VARCHAR(255)"`
	RepositorieTagName       string `xorm:"VARCHAR(255)"`
	RepositorieTagJson       string `xorm:"TEXT"`
	RepositorieTag           string `xorm:"VARCHAR(255)"`

	//RepositorieIsCheck bool `xorm:default false`

}

func setEngine() {
	dbHost := utils.Cfg.MustValue("db", "db_host")
	dbName := utils.Cfg.MustValue("db", "db_name")
	dbUser := utils.Cfg.MustValue("db", "db_user")
	dbPasswd := utils.Cfg.MustValue("db", "db_passwd")

	var err error
	//root:root@tcp(192.168.168.208:3306)/docker_registry?charset=utf8  <--这个对了
	//root:root@192.168.168.208:3306/docker_registry?charset=utf8  <--这个已经不对了
	connStr := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8",
		dbUser, dbPasswd,
		dbHost, dbName)
	beego.Trace("Initialized database connStr->", connStr)

	x, err = xorm.NewEngine("mysql", connStr)
	if err != nil {
		log.Fatalf("models.init -> fail to conntect database: %v", err)
	}

	//if beego.RunMode != "pro" {
	//}
	x.ShowDebug = true
	x.ShowErr = true
	x.ShowSQL = true

	beego.Trace("Initialized database ->", dbName)

}

// InitDb initializes the database.
func InitDb() {
	setEngine()
	err := x.Sync(new(RegistryUser), new(RegistryImage), new(RegistryRepositorieTag))
	if err != nil {
		log.Fatalf("models.init -> fail to sync database: %v", err)
	}
}

func GetImageById(imageId string) (returnImage *RegistryImage, err error) {
	returnImage = new(RegistryImage)
	rows, err := x.Where("image_id=?", imageId).Rows(returnImage)
	defer rows.Close()
	if err != nil {
		returnImage = nil
		return
	}
	if rows.Next() {
		rows.Scan(returnImage)
	} else {
		returnImage = nil
	}

	return
}

type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

type OrmError string

func (e OrmError) Error() string {
	return string(e)
}

func GetRegistryUserByUserName(mUserName string) (returnRegistryUser *RegistryUser, err error) {
	returnRegistryUser = new(RegistryUser)
	rows, err := x.Where("user_name=?", mUserName).Rows(returnRegistryUser)
	if rows.Next() {
		rows.Scan(returnRegistryUser)
		return returnRegistryUser, nil
	} else {
		return nil, OrmError("get user by name error")
	}

}

func GetRegistryUserByToken(mUserName string, mToken string) (returnRegistryUser *RegistryUser, err error) {
	returnRegistryUser = new(RegistryUser)
	rows, err := x.Where("user_name=? and user_token=?", mUserName, mToken).Rows(returnRegistryUser)
	if rows.Next() {
		rows.Scan(returnRegistryUser)
		return returnRegistryUser, nil
	} else {
		return nil, OrmError("get user by token error")
	}

}

func UpRegistryUser(upRegistryUser *RegistryUser) (err error) {
	_, err = x.Id(upRegistryUser.Id).Update(upRegistryUser)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetRegistryUserAuth(authUsername string, authPassword string) (err error) {
	mRegistryUser := new(RegistryUser)
	rows, err := x.Where("user_name=? and user_password=?", authUsername, authPassword).Rows(mRegistryUser)

	if rows.Next() {
		return nil
	} else {
		return AuthError("Auth Error")
	}
}

func InsertOneImage(putRegistryImage *RegistryImage) (affected int64, err error) {
	affected, err = x.InsertOne(putRegistryImage)
	return
}

func UpOneImage(putRegistryImage *RegistryImage) (affected int64, err error) {
	affected, err = x.Id(putRegistryImage.Id).Update(putRegistryImage)
	fmt.Println("putRegistryImage.ImageCheckSumed:", putRegistryImage.ImageCheckSumed, "___affected:", affected, "___err:", err)
	return
}

func InsertOneTag(insertRegistryRepositorieTag *RegistryRepositorieTag) (affected int64, err error) {
	affected, err = x.InsertOne(insertRegistryRepositorieTag)
	return
}

func UpOneTag(upRegistryRepositorieTag *RegistryRepositorieTag) (affected int64, err error) {
	affected, err = x.Id(upRegistryRepositorieTag.Id).Update(upRegistryRepositorieTag)
	return
}

func PutOneTag(upRegistryRepositorieTag *RegistryRepositorieTag) (affected int64, err error) {
	rows, err := x.Where("repositorie_tag_name=? and repositorie_tag_namespace=? and repositorie_tag_repository=?",
		upRegistryRepositorieTag.RepositorieTagName,
		upRegistryRepositorieTag.RepositorieTagNamespace,
		upRegistryRepositorieTag.RepositorieTagRepository).Rows(upRegistryRepositorieTag)
	defer rows.Close()
	if rows.Next() {
		x.Id(upRegistryRepositorieTag.Id).Delete(upRegistryRepositorieTag)
	}
	affected, err = x.InsertOne(upRegistryRepositorieTag)
	return
}
