package main

import (
	"flag"
	"fmt"
)

// 运行模式
var mode = flag.String("mode", "server", "run mode: server or client")
// 服务端地址
var serverHost = flag.String("s", "", "Server Host")
// 服务端监听端口
var serverPort = flag.String("p", "50088", "Server listen port")
// 密码
var key = flag.String("k", "bargo", "Transmission password")
// 本地socks监听端口
var clientSocksPort = flag.String("lp", "1080", "client listen socks port")
// 本地http监听端口
var clientHttpPort = flag.String("lp", "1081", "client listen http port")

func main() {
	// 判断运行模式
	switch *mode {
	case "server": // 服务端

	case "client": // 客户端
		if len(*serverHost) == 0 {
			fmt.Println("Please input server host. Example: -s 123.123.123.123")
			return
		}
	default:
		fmt.Println("Please input correct mode. server or client")
		return
	}
}
