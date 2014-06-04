docker-registry
===============

The dockboard's [docker-registry](https://github.com/dockboard/docker-registry) is a [Golang](http://golang.org) version what clone offical [docker-registry](https://github.com/dotcloud/docker-registry), and we add user manage, UI and more features. We will add more backend storage services support like [Qiniu](http://qiniu.com), [Aliyun OSS](http://www.aliyun.com/product/oss), [Baidu Storage](http://developer.baidu.com/cloud/stor), [Tencent COS](http://www.qcloud.com/product/product.php?item=cos) and [OpenStack Swift](http://docs.openstack.org/developer/swift).


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
```

How to Deploy?
==============

How to Initlization MySQL Database?
===================================

```
INSERT INTO mysql.user(Host,User,Password) VALUES ('localhost', 'docker', password('docker'));
FLUSH PRIVILEGES;
CREATE DATABASE `registry` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON registry.* TO docker@localhost IDENTIFIED BY 'docker';
FLUSH PRIVILEGES;
```

Nginx Conf
==========

```
upstream index_upstream {
  server 127.0.0.1:9911;
}

server {
  listen 80; 
  server_name index.dockboard.org;
  rewrite  ^/(.*)$  https://index.dockboard.org/$1  permanent;
}

server {
  listen 443;

  server_name index.dockboard.org;

  access_log /var/log/nginx/index-dockboard.log;
  error_log /var/log/nginx/index-dockboard-error.log;

  ssl on; 
  ssl_certificate /etc/nginx/ssl/index.dockboard/ssl-bundle.crt;
  ssl_certificate_key /etc/nginx/ssl/index.dockboard/index_dockboard.key;

  client_max_body_size 1024m;
  chunked_transfer_encoding on; 

  proxy_redirect     off;
  proxy_set_header   X-Real-IP $remote_addr;
  proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header   X-Forwarded-Proto $scheme;
  proxy_set_header   Host $http_host;
  proxy_set_header   X-NginX-Proxy true;
  proxy_set_header   Connection ""; 
  proxy_http_version 1.1;

  location / { 
    proxy_pass         http://index_upstream;
  }
}
```
