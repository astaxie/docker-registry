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
  Quota    int64     `xorm:"default 100"`
  Size     int64     `xorm:"default 2048"`
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

type Organization struct {
  Id      int64
  Name    string `xorm:"unique not null"`
  Email   string
  Quota   int64     `xorm:"default 200"`
  Size    int64     `xorm:"default 4096"`
  Actived bool      `xorm:"default 0"`
  Created time.Time `xorm:"created"`
  Updated time.Time `xorm:"updated"`
  Version int       `xorm:"version"`
}

type Member struct {
  Id             int64
  OrganizationId int64     `xorm:"not null 'organization_id'"`
  UserId         int64     `xorm:"not null 'user_id'"`
  Actived        bool      `xorm:"default 0"`
  Created        time.Time `xorm:"created"`
  Updated        time.Time `xorm:"updated"`
  Version        int       `xorm:"version"`
}
