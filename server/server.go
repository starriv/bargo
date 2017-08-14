package server

import (
	"net"
	"bargo/util"
	"log"
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
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Println(err)
		}
		go onConnection(conn)
	}
}
