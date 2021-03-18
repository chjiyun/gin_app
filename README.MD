# Gin App

Go Web FrameWork Test

[Gin中文文档](https://gin-gonic.com/zh-cn/docs/)


## 环境配置
1. 安装golang，[官网地址](https://golang.google.cn/)
2. 设置环境变量 GOPATH，可自定义文件夹
3. 设置module代理：
```
go env -w GOPROXY=https://goproxy.cn,direct
```
4. 初始化air
```
# install for macOS, Linux, Windows
# binary will be $(go env GOPATH)/bin/air
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# 配置文件：./.air.toml
```
5. 安装依赖
```
go mod download
```


## 启动

### 原生方式
```
go run ./main.go
```

### air插件
根目录直接运行如下命令
```
air
```

## 打包

 Windows下编译 Linux平台的64位可执行程序：
 1. cmd

```cmd
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o -x .\build\ .\main.go
```

2. powershell

```powershell
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o .\build\ .\main.go
```

```
# Mac下编译Linux, Windows平台的64位可执行程序：
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

# xshell将文件上传到服务器
# -v, --verbose:输出传输过程中的提示信息
# -y, --overwrite:存在同名文件则替换
rz -y -v

# 赋可执行权限:
chmod 773 main
# 可直接运行：
./main
```

## linux后台运行

```shell
# 如果已经运行先杀进程 
lsof -i:3000
kill pid
nohup ./main 1>gin_app.log 2>error.log &
```

* 其中数字1代表标准输出,2代表错误信息输出,还有未使用的标准输入即数字0.
* 1>/dev/null 2>&1 &的意思是标准信息输出到空设备即不做任何处理,也不做任何显示,/dev/null也可以替换为文本文件名(如 out.log),此时标准信息就会输出到指定文件.
* 2>&1表示错误信息输出到标准输出,即输出同1>后面的位置,当然也可以直接指定新文件名(如 err.log).
* 最后的& 表示程序后台运行



## Gin 是什么？

Gin 是一个用 Go (Golang) 编写的 HTTP web 框架。 它是一个类似于 martini 但拥有更好性能的 API 框架, 由于 httprouter，速度提高了近 40 倍。如果你需要极好的性能，使用 Gin 吧。

## 特性

Gin v1 稳定的特性:

* 零分配路由。

* 仍然是最快的 http 路由器和框架。

* 完整的单元测试支持。

* 实战考验。

* API 冻结，新版本的发布不会破坏你的代码。


## 使用 [Gin](https://github.com/gin-gonic/gin) web 框架的知名项目：

* [gorush](https://github.com/appleboy/gorush)：Go 编写的通知推送服务器。

* [fnproject](https://github.com/fnproject/fn)：原生容器，云 serverless 平台。

* [photoprism](https://github.com/photoprism/photoprism)：由 Go 和 Google TensorFlow 提供支持的个人照片管理工具。

* [krakend](https://github.com/devopsfaith/krakend)：拥有中间件的超高性能 API 网关。

* [picfit](https://github.com/thoas/picfit)：Go 编写的图像尺寸调整服务器。

* [gotify](https://github.com/gotify/server)：使用实时 web socket 做消息收发的简单服务器。

* [cds](https://github.com/ovh/cds)：企业级持续交付和 DevOps 自动化开源平台。