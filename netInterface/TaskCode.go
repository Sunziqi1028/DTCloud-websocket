package netInterface

// TaskTypeCode 任务类型码
type TaskTypeCode int

// TaskErrCode 异常码
type TaskErrCode int

var (
	AcceptTask  TaskTypeCode = 1 // 连接任务
	ReceiveTask TaskTypeCode = 2 // 上行分包任务[UDP不需要分包，使用这个任务池处理数据]
	ReplyTask   TaskTypeCode = 4 // 下行任务
	TaskOutTime TaskErrCode  = 0 // 超时
	TaskErr     TaskErrCode  = 1 // 错误
)
