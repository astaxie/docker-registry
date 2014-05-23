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
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
	Version     int       `xorm:"version"`
}

type Tag struct {
	Id         int64
	Name       string
	JSON       string
	ImageId    string
	Repository int64     `xorm:"repository_id"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	Version    int       `xorm:"version"`
}
