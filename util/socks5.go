package util

import (
	"encoding/binary"
	"fmt"
)

// Socks5Head 协议头
type Socks5Head struct {
	Version byte   // 协议版本
	Cmd     byte   // 连接类型
	Atyp    byte   // host类型
	Addr    string // host地址
	Port    string // 端口
}

// NewSocks5Head 实例化
func NewSocks5Head(data []byte) (Socks5Head, error) {
	s := Socks5Head{}
	s.Version = data[0]
	s.Cmd = data[1]
	s.Atyp = data[3]
	// 错误信息
	var err error
	if s.Version != 5 {
		err = fmt.Errorf("Bad socks version")
		return s, err
	}
	// 判断连接类型
	switch s.Cmd {
	case 1: // tcp
		// 获得地址和端口
		s.Addr, s.Port, err = getAddrPort(data, int(s.Atyp))
		if err != nil {
			return s, err
		}
		return s, nil
	case 2: // bind
		err = fmt.Errorf("Not support bind")
		return s, err
	case 3: // udp
		err = fmt.Errorf("Not support udp")
		return s, err
	}

	err = fmt.Errorf("Bad socks5 head")
	return s, err
}

// 获得地址和端口
func getAddrPort(data []byte, atype int) (string, string, error) {
	var addr string
	var port string
	switch atype {
	case 1: // ipv4
		addr = fmt.Sprintf("%v.%v.%v.%v", data[4], data[5], data[6], data[7])
		port = fmt.Sprintf("%v", binary.BigEndian.Uint16(data[8:10]))
		return addr, port, nil
	case 3: // domainname
		addrLen := int(data[4])
		addr = string(data[5 : 5+addrLen])
		port = fmt.Sprintf("%v", binary.BigEndian.Uint16(data[5+addrLen:7+addrLen]))
		return addr, port, nil
	case 4: // ipv6
		err := fmt.Errorf("Not support ipv6")
		return addr, port, err
	default:
		err := fmt.Errorf("Bad socks5 head")
		return addr, port, err
	}
}
