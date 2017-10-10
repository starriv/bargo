package util

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 获得socks5请求Host
func GetSocksHost(conn net.Conn) (string, error) {
	head := make([]byte, 4)
	_, err := io.ReadFull(conn, head)
	if err != nil {
		return "", err
	}
	cmd := head[1]
	atype := head[3]

	// 获得请求地址和端口
	var addr, port string

	switch atype {
	case 1: // ipv4
		data := make([]byte, 6)
		_, err := io.ReadFull(conn, data)
		if err != nil {
			return "", err
		}
		addr = fmt.Sprintf("%v.%v.%v.%v", data[0], data[1], data[2], data[3])
		port = fmt.Sprintf("%v", binary.BigEndian.Uint16(data[4:]))
	case 3: // domainname
		addrLenData := make([]byte, 1)
		_, err := io.ReadFull(conn, addrLenData)
		if err != nil {
			return "", err
		}
		addrLen := int(addrLenData[0])

		data := make([]byte, addrLen+2)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			return "", err
		}
		addr = string(data[:addrLen])
		port = fmt.Sprintf("%v", binary.BigEndian.Uint16(data[addrLen:]))
	case 4: // ipv6
		data := make([]byte, 18)
		_, err := io.ReadFull(conn, data)
		if err != nil {
			return "", err
		}
		addr = fmt.Sprintf("%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x",
			data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7],
			data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15])
		port = fmt.Sprintf("%v", binary.BigEndian.Uint16(data[16:]))
	default:
		return "", fmt.Errorf("Bad socks5 head")
	}

	if cmd != 1 {
		return "", fmt.Errorf("only support tcp")
	}

	// 加入混淆
	host := addr + ":" + port + "\n"
	host += GetRandomString(GetRandInt(100, 1000))

	return host, nil
}
