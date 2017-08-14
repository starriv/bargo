package server

import (
	"net"
	"bargo/util"
	"log"
	"fmt"
)

// 协议解析器
var protocol *util.Protocol

// 开始服务
func Start(port string, key string)  {
	// 协议解析器
	encryptor := util.NewEncryptor([]byte(key))
	protocol = util.NewProtocol(encryptor)
	// tcp服务
	serv, err := net.Listen("tcp", ":" + port)
	defer serv.Close()
	if err != nil {
		log.Println(err)
	}
	// 启动欢迎信息
	welcome(port, key)
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Panic(err.Error())
		}
		go onConnection(conn)
	}
}

// 启动欢迎信息
func welcome(port string, key string)  {
	fmt.Println("Bargo Socks5 proxy service start success!")
	fmt.Println("mode:", "server")
	fmt.Println("listen port:", port)
	fmt.Println("password:", key)
}
