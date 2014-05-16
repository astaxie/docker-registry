## Official Docker Registry API Specification

### API Specification

#### Docker Image API

##### GET /v1/images/(image_id)/layer

Get image layer for a given image_id
读取 image 的 layer ：local 存储时读取 /registry/images/(image_id)/layer 文件；第三方存储时重定向到第三方存储文件的 URL。

> http://docs.docker.io/en/latest/reference/api/registry_api/#layer

###### Example Request
```
GET /v1/images/088b4505aa3adc3d35e79c031fa126b403200f02f51920fbd9b7c503e87c7a2c/layer HTTP/1.1
Host: registry-1.docker.io
Accept: application/json
Content-Type: application/json
```

###### Example Response
```
HTTP/1.1 200
Vary: Accept
X-Docker-Registry-Version: 0.6.0
Cookie: (Cookie provided by the Registry)
{layer binary data stream}
```
###### Status Codes
* 200 – OK
* 401 – Requires authorization
* 404 – Image not found

##### PUT /v1/images/(image_id)/layer

Put image layer for a given image_id
写入 image 的 layer ：local 存储时写入 /registry/images/(image_id)/layer 文件；第三方存储时保存到第三方的存储空间。

> http://docs.docker.io/en/latest/reference/api/registry_api/#layer

###### Example Request
```
PUT /v1/images/088b4505aa3adc3d35e79c031fa126b403200f02f51920fbd9b7c503e87c7a2c/layer HTTP/1.1
Host: registry-1.docker.io
Transfer-Encoding: chunked
Authorization: Token signature=123abc,repository="foo/bar",access=write
{layer binary data stream}
```

###### Example Response
```
HTTP/1.1 200
Vary: Accept
Content-Type: application/json
X-Docker-Registry-Version: 0.6.0
""
```

###### Status Codes
* 200 – OK
* 401 – Requires authorization
* 404 – Image not found
