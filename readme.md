# 简介
该repo旨在教用户如何用golang开发一个AddonConfig服务。Config服务是指在启动addon runtime之前需要的一些配置。

## 结构描述
1. pkg/meta: 描述此addon需要多少步骤来配置等一些AddonConfig元信息
2. pkg/path: 定义了所有服务的HTTP path
3. pk/service: 所有配置服务的处理函数模块

## Meta
Addon通过AddonIndex提交PR之后注册到AddonStore，之后AddonStore便会通过AddonConfig注册的信息来与AddonConfig交互。交互的第一个步骤就是获取AddonConfig的元信息。只有获取到元信息之后，AddonStore才可以渲染和用户交互的配置步骤。

AddonStore与AddonConfig之间的这次请求通过HTTP Get request/response来实现。
```
const (
    MetaSteps = "/meta/steps"
)
```
A sample response is listed below:
```
http -v GET localhost:8080/meta/steps
GET /meta/steps HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.2



HTTP/1.1 200 OK
Content-Length: 315
Content-Type: application/json; charset=utf-8
Date: Fri, 14 Aug 2020 03:05:04 GMT

{
    "steps": [
        {
            "async": false,
            "async_path": "",
            "description": "Operating fee",
            "index": 1,
            "name": "Operating fee",
            "params": {},
            "path": "/api/v1/fee",
            "tx": true
        },
        {
            "async": false,
            "async_path": "",
            "description": "ontology name service register",
            "index": 2,
            "name": "register ONS",
            "params": {
                "Domain": ""
            },
            "path": "/api/v1/ons",
            "tx": true
        }
    ]
}

```

## Service
