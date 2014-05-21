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
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
	Version  int       `xorm:"version"`
}
