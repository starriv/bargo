package httpproxy

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 所有http请求方法
var httpMethods = map[string]int{
	"GET":     1,
	"POST":    1,
	"DELETE":  1,
	"PUT":     1,
	"HEAD":    1,
	"OPTIONS": 1,
	"CONNECT": 1,
}

// 解析http协议
func httpHeadParse(recvBuf []byte) (int, string, string, error) {
	headerIndex := strings.Index(string(recvBuf), "\r\n\r\n")
	if headerIndex == -1 {
		return 0, "", "", nil
	}
	headerIndex += 4
	// 获得请求方法
	header := recvBuf[:headerIndex]

	var method, url string
	fmt.Sscanf(string(header), "%s%s", &method, &url)
	// 判断method是否合法
	if _, ok := httpMethods[method]; !ok {
		return 0, "", "", fmt.Errorf("method error")
	}

	return getRequestSize(method, string(header)), method, url, nil
}

// 获得请求长度
func getRequestSize(method, header string) int {
	if method == "GET" || method == "OPTIONS" || method == "HEAD" || method == "CONNECT" {
		return len(header)
	}
	reg := regexp.MustCompile("(?i:\r\nContent-Length: ?(\\d+))")
	res := reg.FindAllStringSubmatch(header, -1)
	if len(res) > 0 {
		l, _ := strconv.Atoi(res[0][1])
		return l + len(header)
	}

	return 0
}
