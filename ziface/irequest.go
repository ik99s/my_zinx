package ziface

// IRequest 把客户端请求的连接信息、请求的数据包装到一个request中
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte
}
