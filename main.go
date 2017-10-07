package main

import (
	"fmt"

	"github.com/dawniii/bargo/client"
	"github.com/dawniii/bargo/client/httpproxy"
	"github.com/dawniii/bargo/config"
	"github.com/dawniii/bargo/server"
)

func main() {
	// 初始化参数
	config.Parse()
	// 判断运行模式
	switch *config.Mode {
	case "server": // 服务端
		fmt.Println("\033[31m----------------------\033[0m")
		fmt.Printf("Bargo Server Start\n")
		fmt.Printf("%7s: %s\n", "mode", "server")
		fmt.Printf("%7s: %s\n", "port", *config.ServerPort)
		fmt.Printf("%7s: %s\n", "key", *config.Key)
		fmt.Println("\033[31m----------------------\033[0m")

		server.Start()
	case "client": // 客户端
		if len(*config.ServerHost) == 0 {
			fmt.Println("Please input server host. Example: -server-host 123.123.123.123")
			return
		}
		fmt.Println("\033[31m----------------------\033[0m")
		fmt.Printf("Bargo Client Start\n")
		fmt.Printf("%12s: %s\n", "mode", "client")
		fmt.Printf("%12s: %s\n", "socks5 port", *config.ClientPort)
		fmt.Printf("%12s: %s\n", "http port", *config.ClientHttpPort)
		fmt.Println("\033[31m----------------------\033[0m")

		go client.OpenSysproxy(*config.ClientHttpPort)
		go client.Start()
		httpproxy.Start()
	default:
		fmt.Println("Please input correct mode. server or client")
		return
	}
}
