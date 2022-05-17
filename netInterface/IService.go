package netInterface

import (
	"net/http"
	"time"

	"gitee.com/ling-bin/go-utils/pools"
)

// IService 定义服务器接口
type IService interface {
	GetIsStart() bool                                                                                        // 获取是否启动[true 启动，false 未启动]
	Start()                                                                                                  // 启动服务器方法
	Stop()                                                                                                   // 停止服务器方法
	GetConnMgr() IConnManager                                                                                // 得到链接管理
	GetConn(connId uint64) (IConnection, bool)                                                               // 获取连接
	SetLogHandle(func(level ErrLevel, msg ...interface{}))                                                   // 设置内部异常抛出处理
	SetOnConnAuth(h func(conn IConnection, w http.ResponseWriter, r *http.Request) bool)                     // 连接前验证
	SetOnConnStart(h func(IConnection))                                                                      // 设置连接结束处理方法
	SetOnConnStop(h func(IConnection))                                                                       // 设置连接开始处理方法
	SetOnOneReceive(h func(IConnection, int, []byte))                                                        // 连接上传第一包完整数据(连接,上行数据)
	SetOnReceive(h func(IConnection, int, []byte))                                                           // 连接上传一包完整数据(连接,上行数据)
	SetOnReply(h func(IConnection, int, []byte, bool, string, interface{}, error))                           // 设置下发回调(连接,下发数据,下发数据时带的参数,是否成功,异常信息)
	SetRunTaskError(h func(conn IConnection, taskType TaskTypeCode, taskCode TaskErrCode, task interface{})) // 设置处理上发或下行任务处理异常函数
	GetStartTime() time.Time                                                                                 // 获取服务启动时间
	GetReceiveWorkerPool() pools.ITaskWorkerPool                                                             // 消息处理模块(工作池)
	GetReplyWorkerPool() pools.ITaskWorkerPool                                                               // 消息发送处理器(工作池)
	GetAcceptWorkerPool() pools.ITaskWorkerPool                                                              // 连接接收处理器(工作池)
}
