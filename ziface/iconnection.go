package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接，让当前连接准备开始工作
	Start()

	// Stop 停止连接，结束当前连接的工作
	Stop()

	// GetTCPConnection 获取当前连接的绑定socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接的ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态ip:port
	RemoteAddr() net.Addr

	// Send 发送数据，将数据发送给远程的客户端
	SendMsg(msgIs uint32, data []byte) error
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
