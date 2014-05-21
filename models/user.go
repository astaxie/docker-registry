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

func (user *User) GetByUsername(username string) (has bool, err error) {
  user.Username = username
  has, err = x.Get(user)
  return has, err
}

func GetRegistryUserByToken(mUserName string, mToken string) (returnRegistryUser *User, err error) {
  returnRegistryUser = new(User)
  rows, err := x.Where("username=? and token=?", mUserName, mToken).Rows(returnRegistryUser)
  if rows.Next() {
    rows.Scan(returnRegistryUser)
    return returnRegistryUser, nil
  } else {
    return nil, OrmError("get user by token error")
  }
}

func (user *User) UpRegistryUser() (err error) {
  _, err = x.Id(user.Id).Update(user)
  if err != nil {
    return err
  } else {
    return nil
  }
}

func GetRegistryUserAuth(authUsername string, authPassword string) (err error) {
  mRegistryUser := new(User)
  rows, err := x.Where("username=? and password=?", authUsername, authPassword).Rows(mRegistryUser)

  if rows.Next() {
    return nil
  } else {
    return AuthError("Auth Error")
  }
}
