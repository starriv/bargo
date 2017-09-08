package client

import (
	"github.com/getlantern/sysproxy"
	"log"
)

// 开启系统代理
func OpenSysproxy(port string)  {
	err := sysproxy.EnsureHelperToolPresent("bargo-sysproxy", "Input your password and see the world!", "")
	if err != nil {
		log.Printf("Error EnsureHelperToolPresent: %s\n", err)
		return
	}
	_, err = sysproxy.On("127.0.0.1:"+port)
	if err != nil {
		log.Printf("Error set proxy: %s\n", err)
		return
	}
}
