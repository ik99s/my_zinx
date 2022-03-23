package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//当前连接的ID
	ConnID uint32

	//当前连接的状态
	isClosed bool

	//告知当前连接已经退出/停止的channel(由Reader告知Writer退出)
	ExitChan chan bool

	//无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler

	//连接属性集合，可以供开发者增加属性
	property map[string]interface{}

	//保护连接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]...")
	defer fmt.Println("[Reader exits] ConnID = ", c.ConnID, " Remote Address is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的msg head 8字节二进制流
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给工作池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//根据绑定的MsgID找到对应API业务执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写消息Goroutine，专门发送消息给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]...")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")
	//不断阻塞等待channel的消息，一旦有消息就写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，则Writer也要退出
			return
		}
	}
}

// Start 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("[Conn Start]...ConnID = ", c.ConnID)
	//启动从当前连接读数据的业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()
	//按照开发者传递的创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("[Conn Stop]...ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()
	//告知Writer关闭
	c.ExitChan <- true
	//按照开发者传递的销毁连接之前需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStop(c)
	//将当前连接从ConnMgr中摘除
	c.TcpServer.GetConnMgr().Remove(c)
	//关闭管道
	close(c.ExitChan)
	close(c.msgChan)
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
	//将数据发送给channel
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加属性
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no Property FOUND")
	}
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
