package netService

import (
	"gitee.com/ling-bin/netwebSocket/global"
	"sync"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections sync.Map // 连接记录
}

// NewConnManager 实例化管理
func NewConnManager() netInterface.IConnManager {
	return &ConnManager{}
}

// Add 添加链接[新连接处理]
func (c *ConnManager) Add(conn netInterface.IConnection) {
	connId := conn.GetConnId()
	userInfo := global.GlobalUsers[connId]
	c.connections.Store(userInfo.UUID, conn)
}

// Remove 删除连接
func (c *ConnManager) Remove(conn netInterface.IConnection) {
	connId := conn.GetConnId()
	c.RemoveById(connId)
}

// RemoveById 删除连接
func (c *ConnManager) RemoveById(connId uint64) {
	c.connections.Delete(connId)
}

// Get 利用ConnID获取链接
func (c *ConnManager) Get(connId uint64) (netInterface.IConnection, bool) {
	load, ok := c.connections.Load(connId)
	if ok {
		return load.(netInterface.IConnection), ok
	}
	return nil, false
}

// Range 遍历连接
func (c *ConnManager) Range(hFunc func(connId uint64, value netInterface.IConnection) bool) {
	c.connections.Range(func(key, value interface{}) bool {
		return hFunc(key.(uint64), value.(netInterface.IConnection))
	})
}

// Count 获取个数[内部遍历整个map,调用频率控制在'大于5秒']
func (c *ConnManager) Count() int {
	count := 0
	c.connections.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// ClearConn 清除并停止所有连接
func (c *ConnManager) ClearConn() {
	c.connections.Range(func(key, value interface{}) bool {
		c.RemoveById(key.(uint64))
		return true
	})
}
