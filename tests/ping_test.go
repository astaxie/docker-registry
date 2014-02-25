package test

import (
  "encoding/json"
  "github.com/astaxie/beego"
  _ "github.com/docker-registry/routers"
  . "github.com/smartystreets/goconvey/convey"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestV1Ping(t *testing.T) {
  r, _ := http.NewRequest("GET", "/v1/_ping", nil)
  w := httptest.NewRecorder()
  beego.BeeApp.Handlers.ServeHTTP(w, r)

  Convey("Subject: Test /v1/_ping Endpoint\n", t, func() {
    Convey("Status Code Should Be 200", func() {
      So(w.Code, ShouldEqual, 200)
    })
    Convey("The Result Should Not Be Empty", func() {
      So(w.Body.Len(), ShouldBeGreaterThan, 0)
    })
    Convey("The X-Docker-Registry-Version Should Be Exist", func() {
      if _, found := w.Header()["X-Docker-Registry-Version"]; found {
        So(found, ShouldEqual, true)
      }
    })
    Convey("The X-Docker-Registry-Version Should Be 0.6.0", func() {
      So(w.Header()["X-Docker-Registry-Version"], ShouldContain, "0.6.0")
    })
    Convey("The Result of _ping is JSON format", func() {
      body := make(map[string]interface{})
      err := json.Unmarshal(w.Body.Bytes(), &body)
      So(err, ShouldBeNil)
    })
    Convey("The Result of _ping has a Result key", func() {
      body := make(map[string]interface{})
      json.Unmarshal(w.Body.Bytes(), &body)
      if _, found := body["Result"]; found {
        So(found, ShouldEqual, true)
      }
    })
    Convey("The Result of _ping is true", func() {
      body := make(map[string]interface{})
      json.Unmarshal(w.Body.Bytes(), &body)
      So(body["Result"], ShouldEqual, true)
    })
  })
}

func Test_Ping(t *testing.T) {
  r, _ := http.NewRequest("GET", "/_ping", nil)
  w := httptest.NewRecorder()
  beego.BeeApp.Handlers.ServeHTTP(w, r)

  Convey("Subject: Test /_ping Endpoint\n", t, func() {
    Convey("Status Code Should Be 200", func() {
      So(w.Code, ShouldEqual, 200)
    })
    Convey("The Result Should Not Be Empty", func() {
      So(w.Body.Len(), ShouldBeGreaterThan, 0)
    })
    Convey("The X-Docker-Registry-Version Should Be Exist", func() {
      if _, found := w.Header()["X-Docker-Registry-Version"]; found {
        So(found, ShouldEqual, true)
      }
    })
    Convey("The X-Docker-Registry-Version Should Be 0.6.0", func() {
      So(w.Header()["X-Docker-Registry-Version"], ShouldContain, "0.6.0")
    })
    Convey("The Result of _ping is JSON format", func() {
      body := make(map[string]interface{})
      err := json.Unmarshal(w.Body.Bytes(), &body)
      So(err, ShouldBeNil)
    })
    Convey("The Result of _ping has a Result key", func() {
      body := make(map[string]interface{})
      json.Unmarshal(w.Body.Bytes(), &body)
      if _, found := body["Result"]; found {
        So(found, ShouldEqual, true)
      }
    })
    Convey("The Result of _ping is true", func() {
      body := make(map[string]interface{})
      json.Unmarshal(w.Body.Bytes(), &body)
      So(body["Result"], ShouldEqual, true)
    })
  })
}
