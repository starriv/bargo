package bargo

import (
	"encoding/binary"
	"io"
)

// Protocol 传输协议（4字节定长包头）
type Protocol struct{}

// Encode 按照协议打包
func (Protocol) Encode(buf []byte) []byte {
	bufLen := len(buf)
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(bufLen))
	result := append(head, buf...)

	return result
}

// Decode 按照传输协议解包
func (Protocol) Decode(read io.Reader) ([]byte, error) {
	head := make([]byte, 4)
	// 获得头部
	_, err := io.ReadFull(read, head)
	if err != nil {
		return nil, err
	}
	// 获得包长
	size := binary.BigEndian.Uint32(head)
	data := make([]byte, size)
	_, err = io.ReadFull(read, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}