package netInterface

import (
	"net"
	"net/http"
	"time"
)

// IConnection 连接的接口
type IConnection interface {
	GetNetConn() interface{} // 获取连接
	GetNetwork() string      // 获取网络类型[ws,wss]
	GetConnId() uint64       // 获取客户端ID
	GetRemoteAddr() net.Addr // 获取远程客户端地址信息
	GetLocalAddr() net.Addr  // 获取本地地址
	GetHeartTime() time.Time // 心跳时间
	// SendData 发送消息到客户端：消息类型 1.TextMessage（文本） 2、BinaryMessage(二进制) cmdCode  指令标识[如: rep 普通回复, cmd 用户操作下发 。。]
	SendData(msgType int, data []byte, cmdCode string) error
	// SendDataCall 发送消息到客户端带回调：消息类型 1.TextMessage（文本） 2、BinaryMessage(二进制) cmdCode  指令标识[如: rep 普通回复, cmd 用户操作下发 。。]
	SendDataCall(msgType int, data []byte, cmdCode string, param interface{}, callFunc func(IConnection, int, []byte, bool, string, interface{}, error)) error
	SetProperty(key string, value interface{})           // 设置链接属性
	GetProperty(key string) (interface{}, error)         // 获取链接属性
	RemoveProperty(key string)                           // 移除链接属性
	GetPropertyKeys() []string                           // 获取所有属性key
	GetRecInfo() (count, byteSize uint64)                // 上行当前处理的包总数（处理前，1开始），总大小(字节)
	GetRepInfo() (count, byteSize, errCount uint64)      // 下行当前处理的包总数（处理后），总大小(字节)
	CallLogHandle(level ErrLevel, msgAry ...interface{}) // 设置内部异常抛出处理
	Start()                                              // 启动连接，让当前连接开始工作
	Stop()                                               // 停止连接，结束当前连接状态
	GetIsClosed() bool                                   // 获取的状态（ture:关闭状态，false:未关闭）
	GetRequest() *http.Request                           // 获取请求体
	Incr(val int64) int64                                // 连接提供给业务作为流水号使用,循环累加,(val 为正数为递增，val为负数为递减,val为0则获取值)
}
