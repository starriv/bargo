package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dawniii/bargo/client"
	"github.com/dawniii/bargo/client/httpproxy"
	"github.com/dawniii/bargo/server"
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

// 并发处理连接数量
var connectLimt = flag.String("connect-limit", getEnvArgs("bargo_connect_limit", "100"), "connection limit number")

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
		fmt.Println("\033[31m----------------------\033[0m")
		fmt.Printf("Bargo Server Start\n")
		fmt.Printf("%7s: %s\n", "mode", "server")
		fmt.Printf("%7s: %s\n", "port", *serverPort)
		fmt.Printf("%7s: %s\n", "key", *key)
		fmt.Println("\033[31m----------------------\033[0m")

		server.Start(*serverPort, *key, *connectLimt)
	case "client": // 客户端
		if len(*serverHost) == 0 {
			fmt.Println("Please input server host. Example: -server-host 123.123.123.123")
			return
		}
		fmt.Println("\033[31m----------------------\033[0m")
		fmt.Printf("Bargo Client Start\n")
		fmt.Printf("%12s: %s\n", "mode", "client")
		fmt.Printf("%12s: %s\n", "socks5 port", *clientPort)
		fmt.Printf("%12s: %s\n", "http port", *clientHttpPort)
		fmt.Println("\033[31m----------------------\033[0m")

		go client.OpenSysproxy(*clientHttpPort)
		go client.Start(*serverHost, *serverPort, *clientPort, *key, *connectLimt)
		httpproxy.Start(*clientPort, *clientHttpPort, *clientSysproxy, *clientPac)
	default:
		fmt.Println("Please input correct mode. server or client")
		return
	}
}
