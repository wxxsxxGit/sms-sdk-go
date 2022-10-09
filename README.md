# sms-sdk-go

sms http protocol sdk

# 开发帮助文档

- [【V4】短信服务API接入-帮助文档][]
- [【V4】短信服务 客户自助报备短信模板及其模板发送API接入-帮助文档][]
- [【V4】短信服务接入协议错误码对照表][]
- [【V4】短信运营商失败状态码一览表][]

[【V4】短信服务API接入-帮助文档]:https://api-wiki.wxxsxx.com
[【V4】短信服务 客户自助报备短信模板及其模板发送API接入-帮助文档]:https://www.yuque.com/docs/share/8446f03b-5132-4e87-b8d6-48b9cee0846a
[【V4】短信服务接入协议错误码对照表]:https://thoughts.teambition.com/share/5f22592404ce5e001a397794

[【V4】短信运营商失败状态码一览表]:https://thoughts.teambition.com/share/62f9aa40f3d36d0041586a7f#title=运营商短信失败状态码一览表

# GOLANG

### 开启GO111MODULE="on"

### git下载sdk
```
git clone https://github.com/wxxsxxGit/sms-sdk-go
```

### 自动安装依赖文件

进入下载的sms-sdk-go目录执行
```         
go mod tidy
```

### 编译成二进制文件
进入下载的sms-sdk-go目录执行
```
go build demo.go
```
### 配置文件
```
配置文件名字为sms.yaml
依次从/etc目录
命令执行位置的当前目录
命令执行位置的config目录
命令执行位置的上一级config目录
命令执行位置的上两级config目录
命令执行位置的上三级config目录
查找sms.yaml文件，查找到就作为配置文件
```
配置文件内容如下,配置信息联系运营获取
```
spId: xxxxxxxxx
spKey: xxxxxxxxx
smsSendUrl: xxxxxxxxx
reportUrl: xxxxxxxxx
templateUrl: xxxxxxxxx
```

### 运行二进制文件

```
linux环境
./demo
windows环境
双击执行demo.exe
```
