package netInterface

//ErrLevel 错误级别
type ErrLevel int

const (
	Fatal ErrLevel = 0 //致命
	Error ErrLevel = 1 //错误
	Warn  ErrLevel = 2 //警告
	Info  ErrLevel = 3 //一般信息
)
