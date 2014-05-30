package models

import (
	"time"
)

type Repository struct {
	Id          int64
	Namespace   string    `xorm:"not null"`
	Repository  string    `xorm:"not null"`
	Description string    `xorm:"text"`
	JSON        string    `xorm:"text 'json'"`
	Dockerfile  string    `xorm:"text"`
	Size        int64     `xorm:"default 0"`
	Download    int64     `xorm:"default 0"`
	Uploaded    bool      `xorm:"default 0 'uploaded'"`
	CheckSumed  bool      `xorm:"default 0 'checksumed'"`
	Star        int64     `xorm:"default 0"`
	Trusted     bool      `xorm:"default 0"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
	Version     int       `xorm:"version"`
}

type Tag struct {
	Id         int64
	Name       string    `xorm:"not null"`
	Agent      string    `xorm:"text 'agent'"`
	ImageId    string    `xorm:"text not null 'image_id'"`
	Repository int64     `xorm:"not null repository_id"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	Version    int       `xorm:"version"`
}

type Comment struct {
	Id         int64
	Repository int64     `xorm:"not null repository_id"`
	User       int64     `xorm:"not null user_id"`
	Tag        int64     `xorm:"tag_id"`
	Content    string    `xorm:"text"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	Version    int       `xorm:"version"`
}

type Star struct {
	Id         int64
	Repository int64     `xorm:"repository_id"`
	User       int64     `xorm:"user_id"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	Version    int       `xorm:"version"`
}
