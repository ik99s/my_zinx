package ziface

// IRouter 为不同消息对应不同处理方式
// 路由抽象接口，路由里的数据都是IRequest请求
type IRouter interface {
	// PreHandle 在处理conn业务之前的hook方法
	PreHandle(request IRequest)

	// Handle 在处理conn业务的主方法
	Handle(request IRequest)

	// PostHandle 在处理conn业务之后的hook方法
	PostHandle(request IRequest)
}
