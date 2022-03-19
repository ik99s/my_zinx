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
	AddRouter(router IRouter)
}
