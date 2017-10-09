package httpproxy

import (
	"encoding/binary"
	"fmt"
	"github.com/dawniii/bargo/util/pac"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// http代理
type BargoHttp struct{}

// 接收http请求
func (b *BargoHttp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 重置log到标准输出
	log.SetOutput(os.Stdout)

	if globalProxy == "off" { // 智能模式
		if pac.InBlack(req.URL.Hostname()) {
			b.hideProxy(rw, req)
		} else {
			b.nomalProxy(rw, req)
		}
	} else { // 全局模式
		b.hideProxy(rw, req)
	}
}

// 正常代理
func (b *BargoHttp) nomalProxy(rw http.ResponseWriter, req *http.Request) {
	// https websocket 隧道代理
	if req.Method == http.MethodConnect {
		// 获得tcp连接
		hij, ok := rw.(http.Hijacker)
		if !ok {
			return
		}
		client, _, err := hij.Hijack()
		if err != nil {
			return
		}
		// 连接远端
		server, err := net.Dial("tcp", req.URL.Host)
		if err != nil {
			return
		}
		// 响应客户端远端连接成功可以开始通讯
		_, err = client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))
		if err != nil {
			return
		}
		// 开始转发
		go io.Copy(server, client)
		io.Copy(client, server)
		server.Close()
		client.Close()
		return
	} else {
		// http 代理
		transport := http.DefaultTransport
		// 请求远端获得响应
		res, err := transport.RoundTrip(req)
		if err != nil {
			return
		}
		// 补充响应header
		for key, value := range res.Header {
			for _, v := range value {
				rw.Header().Add(key, v)
			}
		}
		// 返回响应
		rw.WriteHeader(res.StatusCode)
		io.Copy(rw, res.Body)
		res.Body.Close()
	}
}

// 隐藏代理
func (b *BargoHttp) hideProxy(rw http.ResponseWriter, req *http.Request) {
	// https websocket 隧道代理
	if req.Method == http.MethodConnect {
		// 获得tcp连接
		hij, ok := rw.(http.Hijacker)
		if !ok {
			return
		}
		client, _, err := hij.Hijack()
		if err != nil {
			return
		}
		// 连接远端
		socksConn, err := b.connectSocks(req.URL.Hostname(), req.URL.Port(), socksPort)
		if err != nil {
			return
		}
		// 响应客户端远端连接成功可以开始通讯
		_, err = client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))
		if err != nil {
			return
		}
		// 开始转发
		go io.Copy(socksConn, client)
		io.Copy(client, socksConn)
		socksConn.Close()
		client.Close()
		return
	} else {
		// http 代理
		// 配置socks 代理
		dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+socksPort, nil, &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		transport := &http.Transport{
			Proxy:               nil,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		// 请求远端获得响应
		res, err := transport.RoundTrip(req)
		if err != nil {
			return
		}
		// 补充响应header
		for key, value := range res.Header {
			for _, v := range value {
				rw.Header().Add(key, v)
			}
		}
		// 返回响应
		rw.WriteHeader(res.StatusCode)
		io.Copy(rw, res.Body)
		res.Body.Close()
	}
}

// 连接socks5服务
func (b *BargoHttp) connectSocks(addr, port, localSocketPort string) (net.Conn, error) {
	socksConn, err := net.DialTimeout("tcp", "127.0.0.1:"+localSocketPort, 10*time.Second)
	if err != nil {
		return nil, err
	}
	// 模拟socks5客户端
	// 客户端第一次发送请求
	_, err = socksConn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return nil, err
	}
	// 服务端第一次响应
	info := []byte{0, 0}
	_, err = io.ReadFull(socksConn, info)
	if err != nil {
		return nil, err
	}
	// 客户端发送连接信息
	_, err = socksConn.Write(b.newSocks5Head(addr, port))
	if err != nil {
		return nil, err
	}
	// 服务端响应ok 转发信息
	ok := make([]byte, 10)
	_, err = socksConn.Read(ok)
	if err != nil || ok[1] != 0x00 {
		return nil, fmt.Errorf("conn socks fail")
	}

	return socksConn, nil
}

// 获得socks5头
func (b *BargoHttp) newSocks5Head(addr, port string) []byte {
	socks5Header := []byte{0x05, 0x01, 0x00}
	// 判断addr是ip地址还是字符串域名
	reg := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)
	if reg.MatchString(addr) { // 是ip地址
		socks5Header = append(socks5Header, byte(0x01))
		// 组合ip地址到协议头
		dstAddr := make([]byte, 4)
		temp := strings.Split(addr, ".")
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
