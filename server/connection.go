package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dawniii/bargo/util"
)

// 处理每个链接
func onConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer conn.Close()
	// 设置连接过期
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}
	// 创建远程连接
	remoteConn, err := NewRemoteConn(conn)
	if err != nil {
		return
	}
	defer remoteConn.Close()
	// 设置远程连接过期
	err = remoteConn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}
	// 代理
	protocol.Pipe(conn, remoteConn)
}

// 创建远程连接
func NewRemoteConn(conn net.Conn) (net.Conn, error) {
	// 客户端发来的第一条消息
	firstData, err := protocol.Decode(conn)
	if err != nil {
		return nil, err
	}
	// 解析socks5协议头
	socks5Head, err := util.NewSocks5Head(firstData)
	if err != nil {
		return nil, err
	}
	// 建立远程连接
	remoteConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", socks5Head.Addr, socks5Head.Port), 10*time.Second)
	if err != nil {
		return nil, err
	}
	// 远程连接建立成功 告诉客户端可以开始转发
	_, err = conn.Write(protocol.Encode([]byte("bargo")))
	if err != nil {
		return nil, err
	}

	return remoteConn, nil
}
