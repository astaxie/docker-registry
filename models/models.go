package models

import (
  "fmt"
  "github.com/astaxie/beego"
  "github.com/dockboard/docker-registry/utils"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-xorm/xorm"
  "log"
  "time"
)

var x *xorm.Engine

type Image struct {
  Id         int64
  ImageId    string `xorm:"unique not null"`
  JSON       string `xorm:"text 'json'"`
  ParentJSON string `xorm:"text 'parent_json'"`
  Checksum   string `xorm:"text"`
  Payload    string `xorm:"text"`
  Uploaded   bool
  CheckSumed bool      `xorm:"'checksumed'"`
  Created    time.Time `xorm:"created"`
  Updated    time.Time `xorm:"updated"`
  Version    int       `xorm:"version"`
}

type Repository struct {
  Id          int64
  Namespace   string `xorm:"unique"`
  Repository  string
  Description string    `xorm:"text"`
  TagName     string    `xorm:"text 'tag_name'"`
  TagJSON     string    `xorm:"text 'tag_json'"`
  Tag         string    `xorm:"text"`
  Created     time.Time `xorm:"created"`
  Updated     time.Time `xorm:"updated"`
  Version     int       `xorm:"version"`
}

func setEngine() {
  host := utils.Cfg.MustValue("mysql", "Host")
  name := utils.Cfg.MustValue("mysql", "Name")
  user := utils.Cfg.MustValue("mysql", "User")
  passwd := utils.Cfg.MustValue("mysql", "Passwd")

  var err error
  conn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", user, passwd, host, name)
  beego.Trace("Initialized database connStr ->", conn)

  x, err = xorm.NewEngine("mysql", conn)
  if err != nil {
    log.Fatalf("models.init -> fail to conntect database: %v", err)
  }

  x.ShowDebug = true
  x.ShowErr = true
  x.ShowSQL = true

  beego.Trace("Initialized database ->", name)

}

// InitDb initializes the database.
func InitDb() {
  setEngine()
  err := x.Sync(new(User), new(Image), new(Repository))
  if err != nil {
    log.Fatalf("models.init -> fail to sync database: %v", err)
  }
}

func GetImageById(imageId string) (returnImage *Image, err error) {
  returnImage = new(Image)
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

func InsertOneImage(putRegistryImage *Image) (affected int64, err error) {
  affected, err = x.InsertOne(putRegistryImage)
  return
}

func UpOneImage(putRegistryImage *Image) (affected int64, err error) {
  affected, err = x.Id(putRegistryImage.Id).Update(putRegistryImage)
  fmt.Println("putRegistryImage.ImageCheckSumed:", putRegistryImage.CheckSumed, "___affected:", affected, "___err:", err)
  return
}

func InsertOneTag(insertRegistryRepositorieTag *Repository) (affected int64, err error) {
  affected, err = x.InsertOne(insertRegistryRepositorieTag)
  return
}

func UpOneTag(upRegistryRepositorieTag *Repository) (affected int64, err error) {
  affected, err = x.Id(upRegistryRepositorieTag.Id).Update(upRegistryRepositorieTag)
  return
}

func PutOneTag(upRegistryRepositorieTag *Repository) (affected int64, err error) {
  rows, err := x.Where("repositorie_tag_name=? and repositorie_tag_namespace=? and repositorie_tag_repository=?",
    upRegistryRepositorieTag.TagName,
    upRegistryRepositorieTag.Namespace,
    upRegistryRepositorieTag.Repository).Rows(upRegistryRepositorieTag)
  defer rows.Close()
  if rows.Next() {
    x.Id(upRegistryRepositorieTag.Id).Delete(upRegistryRepositorieTag)
  }
  affected, err = x.InsertOne(upRegistryRepositorieTag)
  return
}
