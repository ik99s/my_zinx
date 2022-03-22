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

//hello zinx test 自定义路由
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

func main() {
	//1 创建一个server句柄
	s := znet.NewServer("[zinx V0.7]")
	//2 给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//3 启动server
	s.Serve()
}
