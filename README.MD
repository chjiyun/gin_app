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

6. 依赖管理
```
# 升级依赖
go get -u -d github.com/gin-gonic/gin
# 补全清理依赖
go mod tidy
```

7. 单元测试

`-v` 参数作用是无论用例是否通过都显示结果
```
# 运行所有测试文件
go test -v
# 运行该目录下含有Basename的测试函数
go test -run Basename ./app/util -v
# 运行该目录下全字匹配Basename的测试函数
go test -run ^Basename$ ./app/util -v
```

8. 数据源dsn配置加密

本地开发环境配置文件（config.dev.yml）不想暴露真实连接，可用对称加密（des）保护敏感数据，key 可用[哈希算法工具](http://www.jsons.cn/allencrypt/)生成一个长度大于16位的文本，根目录下新建`hashkey.txt`并粘贴密钥，然后git忽略跟踪这个文件

```bash
# 加密dsn连接信息，“src=”后填原文，密钥取自 hashkey.txt
go test -run ^TestEncrypt$ .\app\util -v -args "src=*"
# 解密dsn文本，“src=”后填加密文本
go test -run ^TestDecrypt$ .\app\util -v -args "src=*"
# 使用下面的测试函数生成一个sha1散列信息
go test -run ^TestSha1$ .\app\util -v -args "src=*"
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

## 打包（交叉编译）

 Windows下编译 Linux平台的64位可执行程序：
 1. cmd

```cmd
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o .\build\ .\main.go
```

2. powershell
注意：此方式打包后记得关闭当前shell窗口
```powershell
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o .\build\ .\main.go
```

```shell
# Mac下编译Linux, Windows平台的64位可执行程序：
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

# xshell将文件上传到服务器
# -v, --verbose:输出传输过程中的提示信息
# -y, --overwrite:存在同名文件则替换
rz -y -v
```

## linux后台运行

```shell
# 赋可执行权限
chmod 773 main
# 可直接运行（退出即程序终止）
./main
```

或者后台启动
```shell
# 如果已经运行先杀进程 
lsof -i:8000
kill pid
nohup ./main 1>info_2021_03_27.log 2>error.log & echo $! > pid.txt
```

* 其中数字1代表标准输出,2代表错误信息输出,还有未使用的标准输入即数字0.
* 1>/dev/null 2>&1 &的意思是标准信息输出到空设备即不做任何处理,也不做任何显示,/dev/null也可以替换为文本文件名(如 out.log),此时标准信息就会输出到指定文件.
* 2>&1表示错误信息输出到标准输出,即输出同1>后面的位置,当然也可以直接指定新文件名(如 err.log).
* 最后的& 表示程序后台运行

## 一键部署
```shell
chmod 755 ./build.sh
./build.sh
```

# linux安装go

请前往golang [官网下载页面](https://golang.google.cn/dl/)，复制文件链接地址

```shell
# 下载压缩包
wget https://golang.google.cn/dl/go1.16.3.linux-amd64.tar.gz
# 解压
tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
# ~/.profile 文件末尾添加环境变量
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
# 生效
source ~/.profile
```

# 注意
若windows环境下使用logger后控制台出现这样的错误提示 "failed to rotate: failed to create new symlink ***.log ..."，
请按照如下步骤解决：
1. win + r 打开运行弹窗，输入secpol.msc
2. 弹出一个本地安全策略窗口，依次找到安全设置>本地策略>用户权限分配，双击后右侧出现很多行策略，找到创建符号链接双击
3. 新弹出的窗口中点击添加用户或组，把当前登录计算机的用户加进去，输用户名搜索即可，然后一路确定，还可能需要重启电脑

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