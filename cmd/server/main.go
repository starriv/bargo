package main

import (
	"flag"
	"os"
	"fmt"
	"net"
	"bargo"
)

// 服务端端口
var port = flag.String("p", "50088", "Server listen port")
// 服务密码
var password = flag.String("k", "bargo", "Transmission password")
// 访问日志
var access_log = flag.Bool("a", false, "open access log")
// 加密器
var encryptor *bargo.Encryptor
// 通讯协议
var proto = bargo.Protocol{}

// 检测输入参数
func checkArgs() {
	exit := false
	if len(*port) == 0 {
		fmt.Println("Please input server port.")
		fmt.Println("Example: -p 50088")
		exit = true
	}
	if exit {
		os.Exit(0)
	}
	fmt.Println("----------------------------------")
	fmt.Println("Bargo server start success!")
	fmt.Println("server listen port:", *port)
	fmt.Println("password:", *password)
	fmt.Println("----------------------------------")
}

// 开始服务
func Start()  {
	// 初始化加密器
	encryptor = bargo.NewEncryptor([]byte(*password))
	// tcp服务
	server, err := net.Listen("tcp", ":" + *port)
	defer server.Close()
	if err != nil {
		bargo.Log(err.Error())
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			bargo.Log(err)
		}
		// 处理每个链接
		go onConnection(conn)
	}
}

// 处理每个连接
func onConnection(conn net.Conn)  {
	defer conn.Close()
	remoteConn, err := linkRemoteConn(conn)
	if err != nil {
		return
	}
	defer remoteConn.Close()
	// 代理
	proxy(conn, remoteConn)
	return
}

// 代理
func proxy(conn, remoteConn net.Conn) error {
	errCh := make(chan error,2)

	go func() {
		_, err := bargo.DecryptCopy(remoteConn, conn, encryptor)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		_, err := bargo.EncryptCopy(conn, remoteConn, encryptor)
		if err != nil {
			errCh <- err
		}
	}()

	for i:=0; i<2; i++ {
		err := <-errCh
		return err
	}

	return nil
}

// 建立远程连接
func linkRemoteConn(conn net.Conn) (net.Conn,error) {
	firstData, err := proto.Decode(conn)
	if err != nil {
		return nil, err
	}

	socks5Head, err := bargo.NewSocks5Head(encryptor.Decrypt(firstData))
	if err != nil {
		return nil, err
	}

	// 建立远程连接
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", socks5Head.Addr, socks5Head.Port))
	if err != nil {
		return nil, err
	}

	// 告诉客户端可以开始转发
	_, err = conn.Write(proto.Encode(encryptor.Encrypt([]byte("bargo"))))
	if err != nil {
		return nil, err
	}
	// 记录访问日志
	if *access_log {
		bargo.Log(conn.RemoteAddr(),"->",socks5Head.Addr,socks5Head.Port)
	}

	return remoteConn, nil
}

func main() {
	// 解析参数
	flag.Parse()
	// 判断参数
	checkArgs()
	// 开始服务
	Start()
}