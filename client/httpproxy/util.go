package httpproxy

import (
	"encoding/binary"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 获得地址和端口
func parseAddrPort(host, method string) (string, string, error) {
	var addr, port string // 地址 端口
	// 获得目标服务器地址和端口
	if method == "CONNECT" { // https
		temp := strings.Split(host, ":")
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
			temp := strings.Split(hostPortURL.Host, ":")
			addr = temp[0]
			port = temp[1]
		}
	}

	return addr, port, nil
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
