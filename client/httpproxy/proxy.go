package httpproxy

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/dawniii/bargo/util/pac"
)

// 正常的http代理
func proxy(conn net.Conn, mode string) {
	// 当前包长
	curPackLen := 0
	// 接收缓存
	recvBuf := make([]byte, 0)
	// 每次接收数据
	readBuf := make([]byte, 4096)
	// 远端连接
	var remoteConn net.Conn

	// 读取客户端数据
	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			return
		}
		recvBuf = append(recvBuf, readBuf[:n]...)
		// 超过10m没有读到包头 退出
		if len(recvBuf) >= 10485760 {
			return
		}
		// 循环解析buf中的数据
		for {
			if len(recvBuf) <= 0 {
				break
			}
			// 没有读到包头
			if curPackLen == 0 {
				// 尝试解析包头
				var method, url string
				curPackLen, method, url, err = httpHeadParse(recvBuf)
				if err != nil { // 解析出错 退出
					return
				}
				if curPackLen == 0 { // 没有解析到包头 继续读取
					break
				}
				// 读取成功
				if curPackLen > 0 {
					// 连接远端
					if remoteConn == nil {
						if mode == "off" { // 自动模式
							if pac.InBlack(url) {
								remoteConn, err = linkToSocksRemote(conn, method, url)
							} else {
								remoteConn, err = linkToRemote(conn, method, url)
							}
						} else { // 全局模式
							remoteConn, err = linkToSocksRemote(conn, method, url)
						}
						if err != nil {
							return
						}
					}
					// 转发数据到远端
					recvBuf, err = copyToRemote(recvBuf, curPackLen, conn, remoteConn)
					if err != nil {
						remoteConn.Close()
						return
					}
					curPackLen = 0
				}
			}
		}
	}
}

// 连接到远端
func linkToRemote(conn net.Conn, method, url string) (net.Conn, error) {
	addr, port, err := parseAddrPort(url, method)
	if err != nil {
		return nil, err
	}
	var remoteConn net.Conn
	remoteConn, err = net.DialTimeout("tcp", addr+":"+port, 30*time.Second)
	if err != nil {
		return nil, err
	}
	// https 直接转发
	if method == "CONNECT" {
		_, err = fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			return nil, err
		}
		// 直接进行转发
		go io.Copy(remoteConn, conn)
		io.Copy(conn, remoteConn)
		remoteConn.Close()
		return nil, fmt.Errorf("closed")
	}
	go func() {
		defer func() {
			remoteConn.Close()
		}()
		_, err = io.Copy(conn, remoteConn)
		if err != nil {
			return
		}
	}()
	return remoteConn, nil
}

// 连接到本地socks服务
func linkToSocksRemote(conn net.Conn, method, url string) (net.Conn, error) {
	var remoteConn net.Conn
	var err error

	remoteConn, err = net.DialTimeout("tcp", "127.0.0.1:"+socksPort, 10*time.Second)
	if err != nil {
		return nil, err
	}
	// 模拟socks5客户端
	// 客户端第一次发送请求
	_, err = remoteConn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return nil, err
	}
	// 服务端第一次响应
	info := []byte{0, 0}
	_, err = io.ReadFull(remoteConn, info)
	if err != nil {
		return nil, err
	}
	// 客户端发送连接信息
	_, err = remoteConn.Write(newSocks5Head(url, method))
	if err != nil {
		return nil, err
	}
	// 服务端响应ok 转发信息
	info2 := make([]byte, 10)
	_, err = remoteConn.Read(info2)
	if err != nil || info2[1] != 0x00 {
		return nil, fmt.Errorf("conn socks fail")
	}

	// https 直接转发
	if method == "CONNECT" {
		_, err = fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			return nil, err
		}
		// 直接进行转发
		go io.Copy(remoteConn, conn)
		io.Copy(conn, remoteConn)
		remoteConn.Close()
		return nil, fmt.Errorf("closed")
	}
	go func() {
		defer func() {
			remoteConn.Close()
		}()
		_, err = io.Copy(conn, remoteConn)
		if err != nil {
			return
		}
	}()
	return remoteConn, nil
}

// 发送数据到远端
func copyToRemote(recvBuf []byte, curPackLen int, conn, remoteConn net.Conn) ([]byte, error) {
	switch {
	case len(recvBuf) == curPackLen: // 刚好读取了一个包
		err := sendHeader(recvBuf, remoteConn)
		if err != nil {
			return nil, err
		}
		return make([]byte, 0), nil
	case len(recvBuf) > curPackLen: // 读取了一个完整包并多读取到了内容
		header := recvBuf[:curPackLen]
		newRecvBuf := recvBuf[curPackLen:]
		err := sendHeader(header, remoteConn)
		if err != nil {
			return nil, err
		}
		return newRecvBuf, nil
	case len(recvBuf) < curPackLen: // 还有body没读取或者没读取完
		err := sendHeader(recvBuf, remoteConn)
		if err != nil {
			return nil, err
		}
		wantReadlen := curPackLen - len(recvBuf)
		readBodylen := 0
		// 每次接收数据
		readBodyBuf := make([]byte, wantReadlen)
		for {
			n, err := conn.Read(readBodyBuf)
			if err != nil {
				return nil, err
			}
			// 转发
			_, err = remoteConn.Write(readBodyBuf[:n])
			if err != nil {
				return nil, err
			}
			// 计数
			readBodylen += n
			if readBodylen < wantReadlen {
				readBodyBuf = make([]byte, wantReadlen-readBodylen)
				continue
			} else {
				break
			}
		}
		return make([]byte, 0), nil
	}

	return nil, fmt.Errorf("unknow")
}

// 发送头部信息到远端
func sendHeader(recvBuf []byte, remoteConn net.Conn) error {
	firstLineIndex := strings.Index(string(recvBuf), "\r\n")
	firstLine := recvBuf[:firstLineIndex]
	reg := regexp.MustCompile(`(?i:http://.*?/)`)
	firstLine = reg.ReplaceAll(firstLine, []byte("/"))
	newHeader := append(firstLine, recvBuf[firstLineIndex:]...)

	_, err := remoteConn.Write(newHeader)
	if err != nil {
		return err
	}
	return nil
}
