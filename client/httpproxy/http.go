package httpproxy

import (
	"github.com/dawniii/bargo/config"
	"github.com/dawniii/bargo/util/pac"
	"net/http"
)

// 监听端口
var httpPort string
var socksPort string

// 全局模式或者局部模式
var globalProxy string

// 用户定义的pac域名ip
var userPac string

// 开启http代理 转发数据到socks代理
func Start() {
	httpPort = *config.ClientHttpPort
	socksPort = *config.ClientPort
	globalProxy = *config.ClientSysproxy
	userPac = *config.ClientPac

	// 添加用户自定义规则
	if len(userPac) > 0 {
		pac.AddRules(userPac)
	}
	// 启动服务
	mux := new(BargoHttp)
	http.ListenAndServe(":"+httpPort, mux)
}
