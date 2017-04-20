package bargo

import (
	"net"
	"time"
)

// 传输协议
var protocol = Protocol{}
// 允许连接空闲时间
const CONNECT_IDLE_TIME = 30 * time.Second

// DecryptCopy 解密转发
func DecryptCopy(dst net.Conn, src net.Conn, encryptor *Encryptor) (written int64, err error) {
	for {
		err = src.SetDeadline(time.Now().Add(CONNECT_IDLE_TIME))
		if err != nil {
			return 0, err
		}
		err = dst.SetDeadline(time.Now().Add(CONNECT_IDLE_TIME))
		if err != nil {
			return 0, err
		}
		data, err := protocol.Decode(src)
		if err != nil {
			return 0, err
		}
		// 转发
		nw, ew := dst.Write(encryptor.Decrypt(data))
		if nw > 0 {
			written += int64(nw)
		}
		if ew != nil {
			err = ew
			break
		}
	}

	return written, err
}

// EncryptCopy 加密转发
func EncryptCopy(dst net.Conn, src net.Conn, encryptor *Encryptor) (written int64, err error) {
	buf := make([]byte, 4096)
	for {
		err = src.SetDeadline(time.Now().Add(CONNECT_IDLE_TIME))
		if err != nil {
			return 0, err
		}
		err = dst.SetDeadline(time.Now().Add(CONNECT_IDLE_TIME))
		if err != nil {
			return 0, err
		}
		nr, er := src.Read(buf)
		if nr > 0 {
			data := protocol.Encode(encryptor.Encrypt(buf[0:nr]))
			nw, ew := dst.Write(data)
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}