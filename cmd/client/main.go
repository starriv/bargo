package main

import (
	"flag"
	"bargo"
	"fmt"
	"os"
	"net"
	"io"
)

// 本地监听端口
var localPort = flag.String("lp", "1080", "Local listen port")

// 服务端地址
var serverHost = flag.String("s", "", "Server Host")

// 服务端端口
var serverPort = flag.String("p", "50088", "Server listen port")

// 服务密码
var password = flag.String("k", "bargo", "Transmission password")

// 加密工具
var encryptor *bargo.Encryptor

// 通讯协议
var protocol = bargo.Protocol{}

// 检测输入参数
func checkArgs() {
	exit := false
	if len(*serverHost) == 0 {
		fmt.Println("Please input server host.")
		fmt.Println("Example: -s 123.123.123.123")
		exit = true
	}
	if len(*serverPort) == 0 {
		fmt.Println("Please input server port.")
		fmt.Println("Example: -p 50088")
		exit = true
	}
	if exit {
		os.Exit(0)
	}
	fmt.Println("----------------------------------")
	fmt.Println("Bargo local server start success!")
	fmt.Println("server host:", *serverHost)
	fmt.Println("server listen port:", *serverPort)
	fmt.Println("password:", *password)
	fmt.Println("local listen port:", *localPort)
	fmt.Println("----------------------------------")
}

func main()  {
	// 解析输入参数
	flag.Parse()
	// 检测输入参数
	checkArgs()
	// 初始化加密工具
	encryptor = bargo.NewEncryptor([]byte(*password))
	// 开启tcp服务
	startTcpServer()
}

// 开启tcp服务
func startTcpServer() {
	server, err := net.Listen("tcp", ":"+ *localPort)
	if err != nil {
		panic(err.Error())
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			bargo.Log(err.Error())
		}
		// 处理每个链接
		go onConnection(conn)
	}
}

// 处理链接
func onConnection(conn net.Conn) {
	defer conn.Close()
	// socks5版本验证
	err := checkVersion(conn)
	if err != nil {
		return
	}
	// 建立远程连接
	remoteConn, err := linkRemoteConn(conn)
	if err != nil {
		return
	}
	defer remoteConn.Close()
	// 开始代理
	proxy(conn, remoteConn)
	return
}

// 代理
func proxy(conn, remoteConn net.Conn) error {
	errCh := make(chan error,2)

	go func() {
		_, err := bargo.DecryptCopy(conn, remoteConn, encryptor)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		_, err := bargo.EncryptCopy(remoteConn, conn, encryptor)
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

// 连接远程
func linkRemoteConn(conn net.Conn) (net.Conn, error) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	// 建立远程连接
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", *serverHost, *serverPort))
	if err != nil {
		conn.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return nil, err
	}

	// 向server端加密发送链接信息
	data := protocol.Encode(encryptor.Encrypt(buf[:length]))
	_, err = remoteConn.Write(data)
	if err != nil {
		return nil, err
	}

	// 接收服务端的转发握手
	handdata, err := protocol.Decode(remoteConn)
	if err != nil {
		return nil, fmt.Errorf("not recv remote hander")
	}
	if string(encryptor.Decrypt(handdata)) == "bargo" {
		// 响应客户端消息
		_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		if err != nil {
			return nil, err
		}
		return remoteConn, nil
	} else {
		return nil, fmt.Errorf("remote hander don't right")
	}
}

// checkVersion 判断协议版本和验证方式(不做验证)
func checkVersion(conn net.Conn) error {
	versionInfo := []byte{0, 0, 0}
	_, err := io.ReadFull(conn, versionInfo)
	if err != nil {
		return err
	}
	if versionInfo[0] != 5 {
		_, err = conn.Write([]byte{0x05, 0xff})
		return err
	}
	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return err
	}

	return nil
}
