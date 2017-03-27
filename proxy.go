package bargo

import (
	"io"
)

var protocol = Protocol{}

// DecryptCopy 解密转发
func DecryptCopy(dst io.Writer, src io.Reader, encryptor *Encryptor) (written int64, err error) {
	for {
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
func EncryptCopy(dst io.Writer, src io.Reader, encryptor *Encryptor) (written int64, err error) {
	buf := make([]byte, 4096)
	for {
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