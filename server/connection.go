package server

import (
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/dawniii/bargo/util"
	"strings"
)

// 连接超时关闭时间
const KeepCloseTime = 30

// 处理每个链接
func onConnection(conn net.Conn, connCount *util.ConnectionCount) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	defer func() {
		conn.Close()
		connCount.Add(-1)
	}()
	// 连接数量加1
	connCount.Add(1)
	// 设置连接过期
	err := conn.SetDeadline(time.Now().Add(KeepCloseTime * time.Second))
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
	protocol.Pipe(conn, remoteConn)
}

// 创建远程连接
func NewRemoteConn(conn net.Conn) (net.Conn, error) {
	// 客户端发来的第一条消息
	hostData, err := protocol.Decode(conn)
	if err != nil {
		return nil, err
	}
	// 获得请求地址
	hostIndex := strings.Index(string(hostData), "\n")
	if hostIndex == -1 {
		return nil, fmt.Errorf("bad host")
	}
	host := string(hostData[:hostIndex])
	// 验证后缀
	if string(hostData[hostIndex+1:hostIndex+6]) != "bargo" {
		// 远程连接建立失败
		hand := "error" + string(hostData[hostIndex:])
		conn.Write(protocol.Encode([]byte(hand)))
		return nil, fmt.Errorf("bad host suffix")
	}
	// 建立远程连接
	remoteConn, err := net.DialTimeout("tcp", host, KeepCloseTime*time.Second)
	if err != nil {
		return nil, err
	}
	// 远程连接建立成功 告诉客户端可以开始转发
	hand := string(hostData[hostIndex+1:])
	_, err = conn.Write(protocol.Encode([]byte(hand)))
	if err != nil {
		return nil, err
	}

	return remoteConn, nil
}
