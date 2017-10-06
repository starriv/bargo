package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// 连接超时关闭时间
const KeepCloseTime = 30

// 处理每个链接
func onConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer conn.Close()
	// 设置连接过期
	err := conn.SetDeadline(time.Now().Add(KeepCloseTime * time.Second))
	if err != nil {
		return
	}
	// socks5版本验证
	err = checkVersion(conn)
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
	err = remoteConn.SetDeadline(time.Now().Add(KeepCloseTime * time.Second))
	if err != nil {
		return
	}
	// 代理
	protocol.Pipe(remoteConn, conn)
}

// 创建远程连接
func NewRemoteConn(conn net.Conn) (net.Conn, error) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	// 建立远程连接
	remoteConn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", serverHost, serverPort), KeepCloseTime*time.Second)
	if err != nil {
		conn.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return nil, err
	}
	// 向server端加密发送链接信息
	data := protocol.Encode(buf[:length])
	_, err = remoteConn.Write(data)
	if err != nil {
		return nil, err
	}
	// 接收服务端的转发握手
	handdata, err := protocol.Decode(remoteConn)
	if err != nil {
		return nil, fmt.Errorf("not recv remote hander")
	}
	if string(handdata) == "bargo" {
		// 响应客户端消息
		_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		if err != nil {
			return nil, err
		}
		return remoteConn, nil
	} else {
		return nil, fmt.Errorf("remote hander don't right")
	}
}

// checkVersion 判断协议版本和验证方式(不做验证)
func checkVersion(conn net.Conn) error {
	versionInfo := []byte{0, 0, 0}
	_, err := io.ReadFull(conn, versionInfo)
	if err != nil {
		return err
	}
	if versionInfo[0] != 5 {
		_, err = conn.Write([]byte{0x05, 0xff})
		return err
	}
	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return err
	}

	return nil
}
