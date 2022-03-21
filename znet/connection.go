package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//当前连接的ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经退出/停止的channel
	ExitChan chan bool

	//该连接处理的方法Router
	Router ziface.IRouter
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID = ", c.ConnID, " Reader exits, Remote Address is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的msg head 8字节二进制流
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
		}
		//拆包，得到msgID和msgDataLen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack head error ", err)
			break
		}
		//根据dataLen读取data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//从路由中，找到注册绑定的Conn对应的router调用
		//c.Router.PreHandle(req) 这句为啥报错
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// Start 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start...ConnID = ", c.ConnID)
	//启动从当前连接读数据的业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务

}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop...ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()
	//关闭管道
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态ip:port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包
	//msgDataLen+msgID+data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack msg error")
	}
	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id = ", msgId, " error : ", err)
		return errors.New("conn Write error")
	}
	return nil
}
