package ziface

// IConnManager 连接管理模块抽象层
type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 删除连接
	Remove(conn IConnection)
	// Get 根据connID获取连接
	Get(connID uint32) (IConnection, error)
	// Len 得到当前连接总数
	Len() int
	// ClearConn 清除并终止所有连接
	ClearConn()
}
