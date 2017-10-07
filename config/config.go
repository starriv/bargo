package config

import (
	"flag"
	"os"
)

// 是否初始化过
var isParsed = false

// 运行模式
var Mode = flag.String("mode", getEnvArgs("bargo_mode", "server"), "run mode: server or client")

// 服务端地址
var ServerHost = flag.String("server-host", getEnvArgs("bargo_server_host", ""), "Server Host")

// 服务端监听端口
var ServerPort = flag.String("server-port", getEnvArgs("bargo_server_port", "50088"), "Server listen port")

// 密码
var Key = flag.String("key", getEnvArgs("bargo_key", "bargo"), "Transmission password")

// 本地socks监听端口
var ClientPort = flag.String("client-port", getEnvArgs("bargo_client_port", "1080"), "client listen socks port")

// 本地http监听端口
var ClientHttpPort = flag.String("client-http-port", getEnvArgs("bargo_client_http_port", "1081"), "client listen http port")

// 本地开启全局代理
var ClientSysproxy = flag.String("client-sysproxy", getEnvArgs("bargo_client_sysproxy", "off"), "client open global system proxy")

// 添加需要代理的域名
var ClientPac = flag.String("client-pac", getEnvArgs("bargo_client_pac", ""), "client pac domain or ip, use | split")

// 并发处理连接数量
var ConnectLimt = flag.String("connect-limit", getEnvArgs("bargo_connect_limit", "0"), "max keep alive connection number, 0 not limit")

// 优先获取环境变量作为默认值
func getEnvArgs(key string, def string) string {
	v := os.Getenv(key)
	if len(v) != 0 {
		return v
	}
	return def
}

// 初始化
func Parse() {
	if !isParsed {
		flag.Parse()
		isParsed = true
	}
}
