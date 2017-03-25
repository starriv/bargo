package bargo

import "log"

// 软件名称
const SOFT_NAME = "Bargo"
// 记录日志
func Log(msg ...interface{}) {
	log.Println("["+SOFT_NAME+"]",msg)
}
