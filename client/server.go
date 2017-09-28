package client

import (
	"log"
	"net"

	"github.com/dawniii/bargo/util"
	"strconv"
)

// 协议解析器
var protocol *util.Protocol

// 服务端地址
var serverHost string

// 服务端口
var serverPort string

// 开始服务
func Start(sHost, sPort, clientPort, key, connLimitNum string) {
	// 初始化参数
	serverHost = sHost
	serverPort = sPort
	// conn limit
	climit, err := strconv.Atoi(connLimitNum)
	if err != nil {
		log.Fatalln("conn limit err:", err)
	}
	connCount := new(util.ConnectionCount)

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
		// 连接数量计数
		if connCount.Get() < climit {
			go onConnection(conn, connCount)
		}
	}
}
