package test

import (
  "github.com/astaxie/beego"
  _ "github.com/docker-registry/routers"
  . "github.com/smartystreets/goconvey/convey"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestUsersGet(t *testing.T) {
  r, _ := http.NewRequest("GET", "/users", nil)
  w := httptest.NewRecorder()
  beego.BeeApp.Handlers.ServeHTTP(w, r)

  Convey("Subject: Test /users Endpoint\n", t, func() {
    Convey("Status Code Should Be 200", func() {
      So(w.Code, ShouldEqual, 404)
    })
    Convey("The Result Should Not Be Empty", func() {
      So(w.Body.Len(), ShouldBeGreaterThan, 0)
    })
    Convey("The X-Docker-Registry-Version Should Be Exist", func() {
      if _, found := w.Header()["X-Docker-Registry-Version"]; found {
        So(found, ShouldEqual, true)
      }
    })
    // Convey("The X-Docker-Registry-Version Should Be 0.6.5", func() {
    //   So(w.Header()["X-Docker-Registry-Version"], ShouldContain, "0.6.5")
    // })
    Convey("The X-Docker-Registry-Standalone Should Be Exist", func() {
      if _, found := w.Header()["X-Docker-Registry-Standalone"]; found {
        So(found, ShouldEqual, true)
      }
    })
    // Convey("The X-Docker-Registry-Standalone Should Be True", func() {
    //   So(w.Header()["X-Docker-Registry-Standalone"], ShouldContain, "true")
    // })
    Convey("The Users Body Should Be \"OK\"", func() {
      So(string(w.Body.Bytes()), ShouldEqual, "")
    })
  })
}
