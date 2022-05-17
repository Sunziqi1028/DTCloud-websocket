package netService

import "time"

// 配置
type Config struct {
	Scheme               string        // 网络：ws,wss
	AddrAry              []string      // 监听地址和端口:绑定的IP加端口：["192.168.1.24:7018",...]
	PathAry              []string      // 监听路径
	ReceiveWorkerSize    uint          // (上行处理)工作池中工作线程个数,,必须2的N次方
	ReceiveTaskQueueSize uint          // (上行处理)单个工作队列缓存任务大小
	ReceiveOutTime       time.Duration // (上行处理)处理任务超时时间 ，IsOutTime=true 时生效
	ReplyWorkerSize      uint          // (下行处理)工作池中工作线程个数,必须2的N次方
	ReplyTaskQueueSize   uint          // (下行处理)单个工作队列缓存任务大小
	ReplyOutTime         time.Duration // (下行处理)超时时间(完整任务) ，IsOutTime=true 时生效
	AcceptWorkerSize     uint          // (连接接入处理)工作池中工作线程个数,必须2的N次方
	AcceptTaskQueueSize  uint          // (连接接入处理)单个工作队列缓存任务大小
	AcceptOutTime        time.Duration // (连接接入处理)超时时间，IsOutTime=true 时生效
	RBufferSize          int           // 读缓存尺寸(字节)
	WBufferSize          int           // 写缓存尺寸(字节)
	SendOutTime          time.Duration // 下行超时时间(秒)
	CertFile             string        // TLS安全连接文件【wss使用】
	KeyFile              string        // TLS安全连接key【wss使用】
	IsOutTime            bool          // 是否需要支持超时
	OverflowDiscard      bool          // 接收，处理，回复溢出是否丢弃【true:丢弃，false：等待处理】
}

// 默认配置
func DefaultConfig(scheme string, addrAry []string, pathAry []string) *Config {
	config := &Config{
		Scheme:               scheme,
		AddrAry:              addrAry,
		PathAry:              pathAry,
		ReceiveWorkerSize:    512,
		ReceiveTaskQueueSize: 2048,
		ReceiveOutTime:       time.Second * 10,
		ReplyWorkerSize:      512,
		ReplyTaskQueueSize:   1024,
		ReplyOutTime:         time.Second * 10,
		AcceptWorkerSize:     128,
		AcceptTaskQueueSize:  1024,
		AcceptOutTime:        time.Second * 10,
		RBufferSize:          4096,
		WBufferSize:          4096,
		KeyFile:              "",
		CertFile:             "",
		SendOutTime:          time.Second * 5,
		IsOutTime:            true,
		OverflowDiscard:      false,
	}
	return config
}
