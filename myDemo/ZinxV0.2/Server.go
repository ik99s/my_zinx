package main

import "zinx/znet"

/*
	基于Zinx框架开发的服务器端应用程序
*/

func main() {
	//创建一个server句柄
	s := znet.NewServer("[zinx V0.2]")
	//启动server
	s.Serve()
}
