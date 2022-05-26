package netService

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// 连接接收对象池
var acceptTaskPool = sync.Pool{
	New: func() interface{} { return new(acceptTask) },
}

// acceptTask 连接接入任务
type acceptTask struct {
	Conn      *websocket.Conn                                                                                 // 连接
	ConnId    uint64                                                                                          // 客户端ID
	CompanyId uint64                                                                                          // 公司ID
	PartnerId uint64                                                                                          // 伙伴ID
	Name      string                                                                                          // 用户名称
	Follow    []uint64                                                                                        // 向关注者发送数据
	Response  http.ResponseWriter                                                                             // 响应
	Request   *http.Request                                                                                   // 请求
	Duration  time.Duration                                                                                   // 超时时间
	OnAccept  func(accept *acceptTask)                                                                        // 接收到一次数据(未分包开始处理不确定是否完成)
	RunError  func(accept *acceptTask, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) // 运行异常回调
}

// newAcceptTask 创建接收对象
func newAcceptTask() *acceptTask {
	return acceptTaskPool.Get().(*acceptTask)
}

// free 回收释放
func (a *acceptTask) free() {
	a.Conn = nil
	a.ConnId = 0
	a.Duration = 0
	a.OnAccept = nil
	a.RunError = nil
	acceptTaskPool.Put(a)
}

// GetTaskId 获取连接id
func (a *acceptTask) GetTaskId() uint64 {
	return a.ConnId
}

// RunTask 处理数据
func (a *acceptTask) RunTask() {
	a.OnAccept(a)
	a.free()
}

// GetDuration 获取超时时间
func (a *acceptTask) GetDuration() time.Duration {
	return a.Duration
}

// CallError 超时回调
func (a *acceptTask) CallError(taskCode int) {
	if a.RunError != nil {
		if taskCode == 0 {
			a.RunError(a, netInterface.AcceptTask, netInterface.TaskOutTime)
		} else {
			a.RunError(a, netInterface.AcceptTask, netInterface.TaskErr)
		}
	}
}
