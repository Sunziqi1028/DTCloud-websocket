package netService

import (
	"sync"
	"time"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// 回发处理对象池
var replyTaskPool = sync.Pool{
	New: func() interface{} { return new(replyTask) },
}

// replyTask 发送数据TCP
type replyTask struct {
	ConnId       uint64                                                                                            // 连接id
	Data         []byte                                                                                            // 发送数据
	MsgType      int                                                                                               // 消息类型
	Param        interface{}                                                                                       // 参数
	CmdCode      string                                                                                            // 业务指定指令码
	Duration     time.Duration                                                                                     // 超时时间
	CallFunc     func(netInterface.IConnection, int, []byte, bool, string, interface{}, error)                     // 回调方法
	RunReplyTask func(replyTask *replyTask)                                                                        // 下发完成回调方法
	RunError     func(replyTask *replyTask, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) // 运行异常回调
}

// newReplyTask 创建接收对象
func newReplyTask() *replyTask {
	return replyTaskPool.Get().(*replyTask)
}

// free 回收释放
func (r *replyTask) free() {
	r.Data = nil
	r.ConnId = 0
	r.Param = nil
	r.Duration = 0
	r.MsgType = 0
	r.RunError = nil
	r.RunReplyTask = nil
	r.CallFunc = nil
	replyTaskPool.Put(r)
}

// GetTaskId 获取任务ID
func (r *replyTask) GetTaskId() uint64 {
	return r.ConnId
}

// RunTask 运行
func (r *replyTask) RunTask() {
	r.RunReplyTask(r)
	r.free()
}

// GetDuration 获取超时时间
func (r *replyTask) GetDuration() time.Duration {
	return r.Duration
}

// CallError 任务超时
func (r *replyTask) CallError(taskCode int) {
	if r.RunError != nil {
		if taskCode == 0 {
			r.RunError(r, netInterface.ReplyTask, netInterface.TaskOutTime)
		} else {
			r.RunError(r, netInterface.ReplyTask, netInterface.TaskErr)
		}
	}
}
