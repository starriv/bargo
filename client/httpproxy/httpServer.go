package httpproxy

import (
	"net"
	"log"
	"fmt"
	"bufio"

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
func Start(clientPort, clientHttpPort, clientSysproxy, clientPac string)  {
	httpPort = clientHttpPort
	socksPort = clientPort
	globalProxy = clientSysproxy
	userPac = clientPac
	serv, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Panic(err.Error())
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
			log.Println(err)
		}
	}()
	defer conn.Close()
	// 读取客户端数据
	connReader := bufio.NewReader(conn)
	httpFirstLine, err := connReader.ReadBytes('\n')
	if err != nil {
		return
	}
	// 解析http请求头
	var method, host string
	fmt.Sscanf(string(httpFirstLine), "%s%s", &method, &host)
	if globalProxy == "on" {
		// 全局科学代理
		hideProxy(httpFirstLine, host, method, conn, connReader)
		return
	}
	// 添加用户自定义规则
	if len(userPac) > 0 {
		pac.AddRules(userPac)
	}
	// pac黑名单判断
	if pac.InBlack(host) {
		// 科学代理
		hideProxy(httpFirstLine, host, method, conn, connReader)
	} else {
		// 正常代理
		defaultProxy(httpFirstLine, host, method, conn, connReader)
	}
}
