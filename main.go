package main

import (
	"flag"
	"fmt"
	"bargo/server"
	"bargo/client"
)

// 运行模式
var mode = flag.String("mode", "server", "run mode: server or client")
// 服务端地址
var serverHost = flag.String("server-host", "", "Server Host")
// 服务端监听端口
var serverPort = flag.String("server-port", "50088", "Server listen port")
// 密码
var key = flag.String("key", "bargo", "Transmission password")
// 本地socks监听端口
var clientPort = flag.String("client-port", "1080", "client listen socks port")
// 本地http监听端口
var clientHttpPort = flag.String("client-http-port", "1081", "client listen http port")

func main() {
	flag.Parse()
	// 判断运行模式
	switch *mode {
	case "server": // 服务端
		server.Start(*serverPort, *key)
	case "client": // 客户端
		if len(*serverHost) == 0 {
			fmt.Println("Please input server host. Example: -server-host 123.123.123.123")
			return
		}
		go client.Start(*serverHost, *serverPort, *clientPort, *key)
		client.HttpStart(*clientPort, *clientHttpPort)
	default:
		fmt.Println("Please input correct mode. server or client")
		return
	}
}
