package client

import (
	"net"
	"log"
	"fmt"
	"net/url"
	"bytes"
	"strings"
	"regexp"
	"encoding/binary"
	"strconv"
	"time"
	"io"
)
// 监听端口
var httpPort string
var socksPort string

// 开启http代理 转发数据到socks代理
func HttpStart(clientPort, clientHttpPort string)  {
	httpPort = clientHttpPort
	socksPort = clientPort
	serv, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Panic(err.Error())
	}
	// 启动欢迎信息
	httpWelcome(clientHttpPort)
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Panic(err)
		}

		go onHttpConnection(conn)
	}
}

// 启动欢迎信息
func httpWelcome(clientHttpPort string)  {
	fmt.Println("http proxy listen port", clientHttpPort)
}

// 处理每个连接
func onHttpConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	defer conn.Close()
	// 读取客户端数据
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	// 解析http请求头
	var method, host string
	fmt.Sscanf(string(buf[:bytes.IndexByte(buf, '\n')]), "%s%s", &method, &host)

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
	_, err = socksServer.Write(newSocks5Head(host))
	if err != nil {
		return
	}
	// 服务端响应ok 转发信息
	info2 := make([]byte, 512)
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
		_, err = socksServer.Write(buf[:n])
		if err != nil {
			return
		}
	}
	//进行转发
	go io.Copy(socksServer, conn)
	io.Copy(conn, socksServer)
}

func newSocks5Head(host string) []byte {
	socks5Header := []byte{0x05, 0x01, 0x00}
	// 解析host
	hostPortURL, err := url.Parse(host)
	if err != nil {
		return nil
	}
	var addr string // 地址
	var port string // 端口
	// 获得目标服务器地址和端口
	if hostPortURL.Opaque == "443" { //https访问
		addr = hostPortURL.Scheme
		port = "443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			addr = hostPortURL.Host
			port = "80"
		} else {
			temp := strings.Split(hostPortURL.Host,":")
			addr = temp[0]
			port = temp[1]
		}
	}
	// 判断addr是ip地址还是字符串域名
	reg := regexp.MustCompile(`^\d\.\d\.\d\.\d$`)
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
