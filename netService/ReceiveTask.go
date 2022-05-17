package netService

import (
	"sync"
	"time"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// 接收处理对象池
var receiveTaskPool = sync.Pool{
	New: func() interface{} { return new(receiveTask) },
}

// receiveTask 客户端请求内容
type receiveTask struct {
	ConnId      uint64                                                                                            // 连接ID
	Data        []byte                                                                                            // 客户端请求的数据
	MsgType     int                                                                                               // 消息类型
	Duration    time.Duration                                                                                     // 超时时间
	OnCompleted func(receive *receiveTask)                                                                        // 完成回调
	RunError    func(receive *receiveTask, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) // 运行异常回调
}

// newReceiveTask 创建接收对象
func newReceiveTask() *receiveTask {
	return receiveTaskPool.Get().(*receiveTask)
}

// free 回收释放
func (r *receiveTask) free() {
	r.Data = nil
	r.ConnId = 0
	r.MsgType = 0
	r.Duration = 0
	r.RunError = nil
	r.OnCompleted = nil
	receiveTaskPool.Put(r)
}

// GetTaskId 获取连接id
func (r *receiveTask) GetTaskId() uint64 {
	return r.ConnId
}

// RunTask 处理数据
func (r *receiveTask) RunTask() {
	r.OnCompleted(r)
	r.free()
}

// GetDuration 获取超时时间
func (r *receiveTask) GetDuration() time.Duration {
	return r.Duration
}

// CallError 任务超时
func (r *receiveTask) CallError(taskCode int) {
	if r.RunError != nil {
		if taskCode == 0 {
			r.RunError(r, netInterface.ReceiveTask, netInterface.TaskOutTime)
		} else {
			r.RunError(r, netInterface.ReceiveTask, netInterface.TaskErr)
		}
	}
}
