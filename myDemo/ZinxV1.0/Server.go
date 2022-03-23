package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于Zinx框架开发的服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Ping Router Handle...")
	//先读取客户端的数据，再回写ping
	fmt.Println("receive from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping..ping..ping.."))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloZinxRouter 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Router Handle...")
	//先读取客户端的数据，再回写ping
	fmt.Println("receive from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnBegin 创建连接之后执行的Hook函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnBegin is called.")
	if err := conn.SendMsg(202, []byte("Do Conn BEGIN")); err != nil {
		fmt.Println(err)
	}
	//给当前的连接设置一些属性
	fmt.Println("Set Conn Name, GitHub...")
	conn.SetProperty("Name", "ChenQ96")
	conn.SetProperty("GitHub", "https://github.com/superIKS/my_zinx")
}

// DoConnEnd 销毁连接之前执行的Hook函数
func DoConnLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnLost is called.")
	fmt.Println("connID = ", conn.GetConnID(), " is lost.")
	//获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("[Name] ", name)
	}
	if git, err := conn.GetProperty("GitHub"); err == nil {
		fmt.Println("[GitHub] ", git)
	}
}

func main() {
	//1 创建一个server句柄
	s := znet.NewServer("[zinx V1.0]")
	//2 注册连接的Hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)
	//3 给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//4 启动server
	s.Serve()
}
