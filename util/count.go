package util

import "sync"

// 连接计数
type ConnectionCount struct {
	num int
	m   sync.RWMutex
}

// 修改计数
func (c *ConnectionCount) Add(n int) {
	c.m.Lock()
	c.num = c.num + n
	c.m.Unlock()
}

// 获得当前计数
func (c *ConnectionCount) Get() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.num
}
