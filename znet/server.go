package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

// Server 定义一个Server服务器模块
type Server struct {
	//服务器名称
	Name string
	//监听ip版本
	IPVersion string
	//监听ip
	IP string
	//监听端口
	Port int
}

// CallbackToClient 定义当前客户端连接所绑定的handle api(目前是写死的，以后优化应该由用户自定义handle方法)
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf error ", err)
		return errors.New("CallbackToClient error")
	}
	return nil
}

// Start 启动
func (s *Server) Start() {
	fmt.Printf("[start] Server Listenner at IP %s, Port %d, is starting...\n", s.IP, s.Port)
	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IP, "error ", err)
			return
		}
		fmt.Println("start Zinx server ", s.Name, "is listening...")

		var cid uint32
		cid = 0
		//3 阻塞等待客户端进行连接，处理客户端连接业务(读写)
		for {
			//如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error ", err)
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, CallbackToClient)
			cid++

			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// Stop 停止
func (s *Server) Stop() {
	//TODO 将服务器资源、状态、已经开辟的连接信息进行停止或回收
}

// Serve 运行
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

// 初始化Server模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
