package models

import (
	"time"
)

type Repository struct {
	Id          int64
	Namespace   string `xorm:"unique"`
	Repository  string
	Description string    `xorm:"text"`
	JSON        string    `xorm:"text 'json'"`
	TagName     string    `xorm:"text 'tag_name'"`
	TagJSON     string    `xorm:"text 'tag_json'"`
	Tag         string    `xorm:"text"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
	Version     int       `xorm:"version"`
}

func PutOneTag(upRegistryRepositorieTag *Repository) (affected int64, err error) {
	rows, err := Engine.Where("tag_name=? and namespace=? and repository=?",
		upRegistryRepositorieTag.TagName,
		upRegistryRepositorieTag.Namespace,
		upRegistryRepositorieTag.Repository).Rows(upRegistryRepositorieTag)
	defer rows.Close()
	if rows.Next() {
		Engine.Id(upRegistryRepositorieTag.Id).Delete(upRegistryRepositorieTag)
	}
	affected, err = Engine.InsertOne(upRegistryRepositorieTag)
	return affected, err
}
