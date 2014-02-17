package main

import (
  "github.com/astaxie/beego"
  _ "github.com/docker-registry/routers"
)

func main() {
  beego.Run()
}
