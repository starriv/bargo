package util

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

// 心跳间隔时间
const HeartbeatInterval = 30

// Protocol 传输协议（4字节定长包头）
type Protocol struct {
	encryptor *Encryptor
}

// 创建协议解析器
func NewProtocol(e *Encryptor) *Protocol {
	p := new(Protocol)
	p.encryptor = e
	return p
}

// 加密数据并打包
func (p *Protocol) Encode(data []byte) []byte {
	// 加密
	buf := p.encryptor.Encrypt(data)
	// 打包
	bufLen := len(buf)
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(bufLen))
	result := append(head, buf...)

	return result
}

// 解密数据并解析协议
func (p *Protocol) Decode(read io.Reader) ([]byte, error) {
	head := make([]byte, 4)
	// 获得头部
	_, err := io.ReadFull(read, head)
	if err != nil {
		return nil, err
	}
	// 获得包长
	size := binary.BigEndian.Uint32(head)
	// 包异常判断
	if size > 4096 {
		return nil, err
	}
	data := make([]byte, size)
	_, err = io.ReadFull(read, data)
	if err != nil {
		return nil, err
	}
	// 解密
	result := p.encryptor.Decrypt(data)

	return result, nil
}

// 转发数据
func (p *Protocol) Pipe(decryptRead, normalRead net.Conn) {
	go func() {
		for {
			err := decryptRead.SetDeadline(time.Now().Add(HeartbeatInterval * time.Second))
			if err != nil {
				break
			}
			err = normalRead.SetDeadline(time.Now().Add(HeartbeatInterval * time.Second))
			if err != nil {
				break
			}
			data, err := p.Decode(decryptRead)
			if err != nil {
				break
			}
			_, err = normalRead.Write(data)
			if err != nil {
				break
			}
		}
	}()

	buf := make([]byte, 1024)
	for {
		err := decryptRead.SetDeadline(time.Now().Add(HeartbeatInterval * time.Second))
		if err != nil {
			break
		}
		err = normalRead.SetDeadline(time.Now().Add(HeartbeatInterval * time.Second))
		if err != nil {
			break
		}
		nr, err := normalRead.Read(buf)
		if err != nil {
			break
		}
		if nr > 0 {
			_, err = decryptRead.Write(p.Encode(buf[0:nr]))
			if err != nil {
				break
			}
		}
	}
}
