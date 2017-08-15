# bargo - 加密的socks5和http代理服务

### 使用方法

1.需要架设服务端和客户端，需要代理的应用，连接客户端的监听端口。

2.直接下载bin目录里编译好的可执行文件或者自行编译当前机器的版本，下面例子是服务器为linux64位，客户端为mac。

#### 服务器端架设

`./bargo-linux-amd64 -mode server -key 设置你的密码 -server-port 监听端口`

#### 客户端（一般是本机，或者路由器）

`./bargo-mac-amd64 -mode client -key 服务端设置的密码 -server-host 服务器ip -server-port 服务器端口`

客户端默认会监听1080端口作为socks5协议端口，1081端口作为http协议端口。`./bargo-xx-xx -h` 查看更多设置参数。

#### 完成

好了，现在你的各种设备或者应用（全局、浏览器、手机等等），都可以连接这两种加密代理了。
