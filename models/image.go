package models

import (
  "time"
)

type Image struct {
  Id         int64
  ImageId    string `xorm:"unique not null"`
  JSON       string `xorm:"text 'json'"`
  ParentJSON string `xorm:"text 'parent_json'"`
  Checksum   string `xorm:"text"`
  Payload    string `xorm:"text"`
  URLs       string `xorm:"text 'urls'"`
  Uploaded   bool
  CheckSumed bool      `xorm:"'checksumed'"`
  Created    time.Time `xorm:"created"`
  Updated    time.Time `xorm:"updated"`
  Version    int       `xorm:"version"`
}
