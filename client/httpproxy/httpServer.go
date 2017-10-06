package httpproxy

import (
	"log"
	"net"

	"github.com/dawniii/bargo/util/pac"
)

// 监听端口
var httpPort string
var socksPort string

// 全局模式或者局部模式
var globalProxy string

// 用户定义的pac域名ip
var userPac string

// 开启http代理 转发数据到socks代理
func Start(clientPort, clientHttpPort, clientSysproxy, clientPac string) {
	httpPort = clientHttpPort
	socksPort = clientPort
	globalProxy = clientSysproxy
	userPac = clientPac
	serv, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Panic(err.Error())
	}
	defer serv.Close()
	// 添加用户自定义规则
	if len(userPac) > 0 {
		pac.AddRules(userPac)
	}
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Panic(err)
		}

		go onHttpConnection(conn)
	}
}

// 处理每个连接
func onHttpConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	defer conn.Close()

	// 添加用户自定义规则
	if len(userPac) > 0 {
		pac.AddRules(userPac)
	}

	if globalProxy == "on" {
		proxy(conn, "auto")
	} else {
		proxy(conn, "all")
	}
}
