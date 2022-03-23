package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动
	Start()
	// Stop 停止
	Stop()
	// Serve 运行
	Serve()

	// AddRouter 路由功能：给当前服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前Server的连接管理器
	GetConnMgr() IConnManager
	//注册OnConnStart Hook函数
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop Hook函数
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
