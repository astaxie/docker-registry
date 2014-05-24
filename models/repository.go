package models

import (
	"time"
)

type Repository struct {
	Id          int64
	Namespace   string
	Repository  string
	Description string    `xorm:"text"`
	JSON        string    `xorm:"text 'json'"`
	Size        int       `xorm:"default 0"`
	Download    int       `xorm:"default 0"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
	Version     int       `xorm:"version"`
}

type Tag struct {
	Id         int64
	Name       string
	JSON       string    `xorm:"text 'json'"`
	ImageId    string    `xorm:"text 'image_id"`
	Repository int64     `xorm:"repository_id"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	Version    int       `xorm:"version"`
}
