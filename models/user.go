package models

import (
	"time"
)

type User struct {
	Id       int64
	Username string `xorm:"unique not null"`
	Password string
	Email    string `xorm:"unique not null"`
	Token    string
	Session  string    `xorm:"unique"`
	Quota    int       `xorm:"default 20"`
	Size     int       `xorm:"default 2048"`
	Actived  bool      `xorm:"default false"`
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
	Version  int       `xorm:"version"`
}

type Profile struct {
	Id       string
	UserId   int64 `xorm:"unique not null 'user_id'"`
	Fullname string
	Company  string
	Location string
	URL      string `xorm:"text 'url'"`
	Gravatar string
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
	Version  int       `xorm:"version"`
}
