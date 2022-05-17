package netInterface

/*IConnManager
  连接管理接口
*/
type IConnManager interface {
	Add(conn IConnection)                                    // 添加链接
	Remove(conn IConnection)                                 // 移除连接
	RemoveById(connId uint64)                                // 移除连接
	Get(connId uint64) (IConnection, bool)                   // 利用ConnID获取链接
	Count() int                                              // 获取个数[内部遍历整个map,调用频率控制在'大于5秒']
	Range(hFunc func(connId uint64, value IConnection) bool) // 遍历
	ClearConn()                                              // 删除并停止所有链接
}
