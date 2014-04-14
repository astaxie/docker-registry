docker-registry [![Build Status](https://drone.io/github.com/dockboard/docker-registry/status.png)](https://drone.io/github.com/dockboard/docker-registry/latest)
===============

The [docker-registry](https://github.com/dockboard/docker-registry) is a [Golang](http://golang.org) version with [Beego](http://beego.me) and [Qiniu](http://qiniu.com) what clone offical [docker-registry](https://github.com/dotcloud/docker-registry). We will add more backend storage service support like [Aliyun OSS](http://www.aliyun.com/product/oss), [Baidu Storage](http://developer.baidu.com/cloud/stor), [Tencent COS](http://www.qcloud.com/product/product.php?item=cos) and [OpenStack Swift](http://docs.openstack.org/developer/swift).


What's Registry?
================

A Registry is a hosted service containing [repositories](http://docs.docker.io/en/latest/terms/repository/#repository-def) of [images](http://docs.docker.io/en/latest/terms/image/#image-def) which responds to the Registry API.


What's FQIN?
============

A Fully Qualified Image Name (FQIN) can be made up of 3 parts:

```
[registry_hostname[:port]/][user_name/](repository_name[:version_tag])
```

`version_tag` defaults to `latest`, `username` and `registry_hostname` default to an empty string. When `registry_hostname` is an empty string, then docker push will push to *index.docker.io:80*.


Why image file named layer?
===========================

In Docker terminology, a read-only [Layer](http://docs.docker.io/en/latest/terms/layer/#layer-def) is called an image. An image never changes.


What's docker 

How to build?
=============

The project build with [gopm](https://github.com/gpmgo/gopm). 

Install [gopm](https://github.com/gpmgo/gopm)

```
go get -u github.com/gpmgo/gopm
```

Build commands:

```
gopm build 
#or
gopm build -v
```

How to Deploy?
==============