package znet

import "zinx/ziface"

// BaseRouter 实现router时，先嵌入这个BaseRouter基类，然后根据需求对这个基类的方法进行重写(继承BaseRouter)
type BaseRouter struct{}

// 这里之所以BaseRouter的方法都为空，是因为有的Router不需要PreHandle或者PostHandle这两个业务
// 所以Router全部继承BaseRouter，而不是实现IRouter接口的好处就是，可以不需要实现PreHandle、PostHandle
// 接口隔离原则

// PreHandle 在处理conn业务之前的hook方法
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 在处理conn业务的主方法
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的hook方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
