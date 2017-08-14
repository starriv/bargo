package client

import (
	"net"
	"bargo/util"
	"log"
	"fmt"
)

// 协议解析器
var protocol *util.Protocol
// 服务端地址
var serverHost string
// 服务端口
var serverPort string

// 开始服务
func Start(sHost, sPort, clientPort, key string)  {
	// 初始化参数
	serverHost = sHost
	serverPort = sPort
	// 协议解析器
	encryptor := util.NewEncryptor([]byte(key))
	protocol = util.NewProtocol(encryptor)
	// tcp服务
	serv, err := net.Listen("tcp", ":" + clientPort)
	defer serv.Close()
	if err != nil {
		log.Panic(err.Error())
	}
	// 启动欢迎信息
	welcome(clientPort)
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Println(err)
		}
		go onConnection(conn)
	}
}

// 启动欢迎信息
func welcome(clientPort string)  {
	fmt.Println("Bargo Socks5 proxy service start success!")
	fmt.Println("mode:", "client")
	fmt.Println("listen port:", clientPort)
}
