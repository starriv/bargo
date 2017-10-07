package client

import (
	"log"
	"net"

	"github.com/dawniii/bargo/config"
	"github.com/dawniii/bargo/util"
)

// 协议解析器
var protocol *util.Protocol

// 服务端地址
var serverHost string

// 服务端口
var serverPort string

// 开始服务
func Start() {
	// 初始化参数
	serverHost = *config.ServerHost
	serverPort = *config.ServerPort
	clientPort := *config.ClientPort
	key := *config.Key
	// 协议解析器
	encryptor := util.NewEncryptor([]byte(key))
	protocol = util.NewProtocol(encryptor)
	// tcp服务
	serv, err := net.Listen("tcp", ":"+clientPort)
	defer serv.Close()
	if err != nil {
		log.Panic(err.Error())
	}
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Println(err)
		}

		go onConnection(conn)
	}
}
