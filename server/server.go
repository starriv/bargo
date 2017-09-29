package server

import (
	"github.com/dawniii/bargo/util"
	"log"
	"net"
	"strconv"
)

// 协议解析器
var protocol *util.Protocol

// 开始服务
func Start(port, key, connLimitNum string) {
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
	serv, err := net.Listen("tcp", ":"+port)
	defer serv.Close()
	if err != nil {
		log.Println(err)
	}
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Panic(err.Error())
		}
		// 连接数量控制
		if climit != 0 && connCount.Get() >= climit {
			conn.Close()
			continue
		}
		go onConnection(conn, connCount)
	}
}
