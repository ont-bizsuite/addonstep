# 简介
该repo旨在教用户如何用golang开发一个AddonConfig服务。Config服务是指在启动addon runtime之前需要的一些配置。

## 结构描述
1. pkg/meta: 描述此addon需要多少步骤来配置等一些AddonConfig元信息
2. pkg/path: 定义了所有服务的HTTP path
3. pk/service: 所有配置服务的处理函数模块

## Meta
Addon通过AddonIndex提交PR之后注册到AddonStore，之后AddonStore便会通过AddonConfig注册的信息来与AddonConfig交互。交互的第一个步骤就是获取AddonConfig的元信息。只有获取到元信息之后，AddonStore才可以渲染和用户交互的配置步骤。

AddonStore与AddonConfig之间的这次请求通过HTTP Get request/response来实现。而Path是AddonConfig与AddonStore之间的约定，也即这个一个硬编码（可以理解为C的main一样），所以所有AddonConfig的Step信息必须放置在下面的路径应答中。
```
const (
    MetaSteps = "/meta/steps"
)
```
下面的Localhost地址可以理解为AddonStore知道的该AddonConfig的服务地址。
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

按照上面的应答的返回AddonStore知道该AddonRuntime的配置步骤有两部，注意Index从1开始。根据Tx标记位（标记是否是笔交易步骤），AddonStore会有不同的处理方式。
### Tx处理方式
AddonStore会发送Get请求来获取Ontology的扫码协议返回，之后进行交易确认和签名，在交易发送到链上之后，AddonStore会将Callback proxy到扫码协议返回的callback中来处理。
请求路径如下，HTTP method： POST
```
    PayCallbackPath = "/api/v1/feeback"
```
### 非Tx处理方式
AddonStore直接将该步骤需要的Param交由前端渲染，待用户输入之后将数据Post给AddonConfig。

###  Step之前的数据传递
如果某一个配置步骤依赖于前一项的配置的输入内容，此时获取数据的方式是通过prev Chain来获取。也即AddonStore在Post当前步骤数据给AddonConfig时候会在对象里加上prev字段来携带之前**已**配置的信息。格式如下：
```
{
...
   prev: {
     step_result: {...},
     render_result: {...}
   }
```


## Service

Service就是每一个步骤的配置处理。这个Demo主要提供了Ontology扫码转账协议以及ONS注册。具体的返回结构请参考Ontology扫码协议。
