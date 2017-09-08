package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dawniii/bargo/server"
	"github.com/dawniii/bargo/client"
	"github.com/dawniii/bargo/client/httpproxy"
)

// 运行模式
var mode = flag.String("mode", getEnvArgs("bargo_mode", "server"), "run mode: server or client")
// 服务端地址
var serverHost = flag.String("server-host", getEnvArgs("bargo_server_host", ""), "Server Host")
// 服务端监听端口
var serverPort = flag.String("server-port", getEnvArgs("bargo_server_port", "50088"), "Server listen port")
// 密码
var key = flag.String("key", getEnvArgs("bargo_key", "bargo"), "Transmission password")
// 本地socks监听端口
var clientPort = flag.String("client-port", getEnvArgs("bargo_client_port", "1080"), "client listen socks port")
// 本地http监听端口
var clientHttpPort = flag.String("client-http-port", getEnvArgs("bargo_client_http_port", "1081"), "client listen http port")
// 本地开启全局代理
var clientSysproxy = flag.String("client-sysproxy", getEnvArgs("bargo_client_sysproxy", "off"), "client open global system proxy")
// 添加需要代理的域名
var clientPac = flag.String("client-pac", getEnvArgs("bargo_client_pac", ""), "client pac domain or ip, use | split")

// 优先获取环境变量作为默认值
func getEnvArgs(key string, def string) string {
	v := os.Getenv(key)
	if len(v) != 0 {
		return v
	}
	return def
}

func main() {
	flag.Parse()
	// 判断运行模式
	switch *mode {
	case "server": // 服务端
		fmt.Println("mode:", "server")
		fmt.Println("listen port:", *serverPort)
		fmt.Println("password:", *key)

		server.Start(*serverPort, *key)
	case "client": // 客户端
		if len(*serverHost) == 0 {
			fmt.Println("Please input server host. Example: -server-host 123.123.123.123")
			return
		}
		fmt.Println("mode:", "client")
		fmt.Println("socks5 proxy listen port:", *clientPort)
		fmt.Println("http proxy listen port", *clientHttpPort)

		go client.OpenSysproxy(*clientHttpPort)
		go client.Start(*serverHost, *serverPort, *clientPort, *key)
		httpproxy.Start(*clientPort, *clientHttpPort, *clientSysproxy, *clientPac)
	default:
		fmt.Println("Please input correct mode. server or client")
		return
	}
}
