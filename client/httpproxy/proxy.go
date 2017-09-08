package httpproxy

import (
	"strings"
	"regexp"
	"strconv"
	"net"
	"bufio"
	"net/url"
	"time"
	"fmt"
	"io"
	"encoding/binary"
)

// 获得地址和端口
func parseAddrPort(host, method string) (string, string, error) {
	var addr, port string // 地址 端口
	// 获得目标服务器地址和端口
	if method == "CONNECT" { // https
		temp := strings.Split(host,":")
		addr = temp[0]
		port = temp[1]
	} else { // http
		hostPortURL, err := url.Parse(host)
		if err != nil {
			return "", "", err
		}
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			addr = hostPortURL.Host
			port = "80"
		} else {
			temp := strings.Split(hostPortURL.Host,":")
			addr = temp[0]
			port = temp[1]
		}
	}

	return addr, port, nil
}

// 正常代理
func defaultProxy(httpFirstLine []byte, host string, method string, conn net.Conn, connReader *bufio.Reader) {
	addr, port, err := parseAddrPort(host, method)
	if err != nil {
		return
	}
	address := addr + ":" + port
	socksServer, err := net.DialTimeout("tcp", address, 10 * time.Second)
	if err != nil {
		return
	}
	defer socksServer.Close()
	// 开始转发信息 响应
	if method == "CONNECT" {
		_, err := fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			return
		}
	} else {
		bufferedData := make([]byte, connReader.Buffered())
		_, err := connReader.Read(bufferedData)
		if err != nil {
			return
		}
		_, err = socksServer.Write(append(httpFirstLine, bufferedData...))
		if err != nil {
			return
		}
	}
	//进行转发
	go io.Copy(socksServer, conn)
	io.Copy(conn, socksServer)
}

// 科学代理
func hideProxy(httpFirstLine []byte, host string, method string, conn net.Conn, connReader *bufio.Reader) {
	//获得了请求的host和port，就开始拨号吧
	socksServer, err := net.DialTimeout("tcp", "127.0.0.1:"+socksPort, 10 * time.Second)
	if err != nil {
		return
	}
	defer socksServer.Close()
	// 模拟socks5客户端
	// 客户端第一次发送请求
	_, err = socksServer.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return
	}
	// 服务端第一次响应
	info := []byte{0, 0}
	_, err = io.ReadFull(socksServer, info)
	if err != nil {
		return
	}
	// 客户端发送连接信息
	_, err = socksServer.Write(newSocks5Head(host, method))
	if err != nil {
		return
	}
	// 服务端响应ok 转发信息
	info2 := make([]byte, 10)
	_, err = socksServer.Read(info2)
	if err != nil || info2[1] != 0x00 {
		return
	}
	// 开始转发信息 响应
	if method == "CONNECT" {
		_, err := fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			return
		}
	} else {
		bufferedData := make([]byte, connReader.Buffered())
		_, err := connReader.Read(bufferedData)
		if err != nil {
			return
		}
		_, err = socksServer.Write(append(httpFirstLine, bufferedData...))
		if err != nil {
			return
		}
	}
	//进行转发
	go io.Copy(socksServer, conn)
	io.Copy(conn, socksServer)
}

// 组合socks5通讯头
func newSocks5Head(host, method string) []byte {
	socks5Header := []byte{0x05, 0x01, 0x00}
	addr, port, err := parseAddrPort(host, method)
	if err != nil {
		return nil
	}
	// 判断addr是ip地址还是字符串域名
	reg := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)
	if reg.MatchString(addr) { // 是ip地址
		socks5Header = append(socks5Header, byte(0x01))
		// 组合ip地址到协议头
		dstAddr := make([]byte, 4)
		temp := strings.Split(addr,".")
		dstAddr0, _ := strconv.Atoi(temp[0])
		dstAddr1, _ := strconv.Atoi(temp[1])
		dstAddr2, _ := strconv.Atoi(temp[2])
		dstAddr3, _ := strconv.Atoi(temp[3])
		dstAddr[0] = byte(dstAddr0)
		dstAddr[1] = byte(dstAddr1)
		dstAddr[2] = byte(dstAddr2)
		dstAddr[3] = byte(dstAddr3)
		socks5Header = append(socks5Header, dstAddr...)
	} else { // 字符串域名地址
		socks5Header = append(socks5Header, byte(0x03))
		// 组合域名到协议头
		socks5Header = append(socks5Header, byte(len(addr)))
		socks5Header = append(socks5Header, []byte(addr)...)
	}
	// 组合端口到协议头
	dstPort, _ := strconv.Atoi(port)
	dstPortByte := make([]byte, 2)
	binary.BigEndian.PutUint16(dstPortByte, uint16(dstPort))
	socks5Header = append(socks5Header, dstPortByte...)

	return socks5Header
}
