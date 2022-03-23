package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	//管理的连接集合
	connections map[uint32]ziface.IConnection
	//保护连接集合的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建当前连接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//将conn加入到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("Connection ", conn.GetConnID(), " add to ConnManager success: conn num = ", cm.Len())
}

// Remove 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//将conn删除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("Connection ", conn.GetConnID(), " delete from ConnManager success: conn num = ", cm.Len())
}

// Get 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		//找到
		return conn, nil
	} else {
		//不存在
		return nil, errors.New("Connection not  FOUND!")
	}
}

// Len 得到当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除并终止所有连接
func (cm *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}
	//TODO 这里调用Len()，写锁有没有问题
	fmt.Println("Clear All Connections success! Conn num = ", cm.Len())
}
