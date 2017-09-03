package client

import (
	"path/filepath"
	"github.com/getlantern/sysproxy"
	"log"
)

// 开启系统代理
func OpenSysproxy(port string)  {
	helperFullPath := "bargo-sysproxy"
	iconFullPath, _ := filepath.Abs("./icon.png")
	err := sysproxy.EnsureHelperToolPresent(helperFullPath, "Input your password and see the world!", iconFullPath)
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
