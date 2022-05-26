package netService

import (
	"fmt"
	"gitee.com/ling-bin/netwebSocket/global"
	"gitee.com/ling-bin/netwebSocket/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"gitee.com/ling-bin/go-utils/pools"
	"github.com/gorilla/websocket"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// Service 服务
type Service struct {
	connId         uint64                                                                                                                       // 客户端ID
	ConnMgr        netInterface.IConnManager                                                                                                    // 连接属性
	OnConnAuth     func(netInterface.IConnection, http.ResponseWriter, *http.Request) bool                                                      // 连接建立前验证（true:成功，false:失败,关闭连接）
	OnConnStart    func(netInterface.IConnection)                                                                                               // 连接完成回调
	OnConnStop     func(netInterface.IConnection)                                                                                               // 关闭回调
	onLogHandle    func(level netInterface.ErrLevel, msg ...interface{})                                                                        // 设置异常处理
	onReceive      func(conn netInterface.IConnection, msgType int, data []byte)                                                                // 数据上传完成
	onOneReceive   func(conn netInterface.IConnection, msgType int, data []byte)                                                                // [第一包]数据上传完成
	onReply        func(conn netInterface.IConnection, msgType int, data []byte, ok bool, cmdCode string, param interface{}, err error)         // 数据下发完成
	runTaskError   func(conn netInterface.IConnection, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode, task interface{}) // 处理上发或下行任务处理异常函数
	receiveHandler pools.ITaskWorkerPool                                                                                                        // 当前Server的消息管理模块(工作池)
	replyHandle    pools.ITaskWorkerPool                                                                                                        // 消息发送处理器(工作池)
	AcceptHandle   pools.ITaskWorkerPool                                                                                                        // 连接处理池
	isStart        bool                                                                                                                         // 是否启动
	config         *Config                                                                                                                      // 配置
	startTime      time.Time                                                                                                                    // 启动时间
	upgrader       *websocket.Upgrader                                                                                                          // 协议升级器
}

// NewService 初始化
func NewService(config *Config) netInterface.IService {
	s := &Service{
		ConnMgr:        NewConnManager(),
		receiveHandler: pools.NewTaskWorkerPool("数据接收处理器", config.ReceiveWorkerSize, config.ReceiveTaskQueueSize),
		replyHandle:    pools.NewTaskWorkerPool("数据回复处理器", config.ReplyWorkerSize, config.ReplyTaskQueueSize),
		AcceptHandle:   pools.NewTaskWorkerPool("连接接收处理器", config.AcceptWorkerSize, config.AcceptTaskQueueSize),
		config:         config,
		isStart:        false,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  config.RBufferSize, // 读取最大值
			WriteBufferSize: config.WBufferSize, // 写最大值
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return s
}

// GetConnMgr 得到链接管理
func (s *Service) GetConnMgr() netInterface.IConnManager {
	return s.ConnMgr
}

// GetConn 获取连接
func (s *Service) GetConn(connId uint64) (netInterface.IConnection, bool) {
	return s.GetConnMgr().Get(connId)
}

// follow 转列表
func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))
	for index, val := range strArr {
		res[index], _ = strconv.Atoi(val)
	}
	return res
}

// wsHandler http到websocket协议升级处理
func (s *Service) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.CallLogHandle(netInterface.Error, "[websocket]websocket协议升级处理异常:", r)
		return
	}

	// 如果用户传入自己的uid 就使用uid 同时传入一个公司company_id,partner_id
	//var uid uint64 = atomic.AddUint64(&s.connId, 1) // 用户ID
	var uid uint64               // 用户ID
	var partner_id uint64 = 0    // 用户Partner_ID
	var company_id uint64 = 0    // 用户组织ID
	var partner_name string = "" // 用户名称
	var follow []uint64          // 关注者
	var messageType = ""         //  消息类型 room ：聊天室 | radio：广播  | orient：定向
	if r.URL.RawQuery != "" {
		values, _ := url.ParseQuery(r.URL.RawQuery)
		intUid, _ := strconv.Atoi(values["uid"][0])
		uid = uint64(intUid)
		ok := utils.CheckUidUnique(uid) // 校验UID 是否唯一
		if !ok {
			log.Println("该用户UID已经存在:", r)
			conn.Close()
			s.ConnMgr.RemoveById(uid)
		}
		intPartnerId, _ := strconv.Atoi(values["partner_id"][0])
		partner_id = uint64(intPartnerId)
		ok = utils.CheckPartnerIDUnique(partner_id) // 校验partner_id 是否唯一
		if !ok {
			log.Println("该伙伴ID已经存在:", r)
			conn.Close()
			s.ConnMgr.RemoveById(uid)
		}
		global.PartnerMap[uid] = partner_id
		fmt.Println("uid:", uid, "partner_id:", partner_id, "service.go line:116")
		intCompanyId, _ := strconv.Atoi(values["company_id"][0])
		company_id = uint64(intCompanyId)
		partner_name = values["name"][0]
		//followTmp = values["follow"][0]
		follow, _ = utils.ConvertString2IntSlice(values["follow"][0])
		messageType = values["type"][0]
		var UserInfo = global.User{
			UID:       uid,
			PartnerID: partner_id,
			CompanyID: company_id,
			Name:      partner_name,
			Type:      messageType,
			Follow:    follow,
		}
		fmt.Println(UserInfo, "Service.go ---line:130")
		global.GlobalUsers[uid] = &UserInfo // 存储全部的用户信息
	}

	acceptTask := newAcceptTask()
	acceptTask.Conn = conn
	acceptTask.ConnId = uid // 创建用户ID
	acceptTask.PartnerId = partner_id
	acceptTask.CompanyId = company_id
	acceptTask.Name = partner_name
	acceptTask.Follow = follow
	// acceptTask.ConnId = atomic.AddUint64(&s.connId, 2) //创建用户ID
	acceptTask.Request = r
	acceptTask.Response = w
	acceptTask.OnAccept = s.runAcceptTask
	acceptTask.RunError = s.runAcceptOutTime
	// 新连接处理
	if s.config.IsOutTime {
		acceptTask.Duration = s.config.AcceptOutTime
	} else {
		acceptTask.Duration = 0
	}
	if !s.config.OverflowDiscard {
		s.AcceptHandle.SendToTaskQueueWait(acceptTask)
		return
	}
	err = s.AcceptHandle.SendToTaskQueue(acceptTask)
	if err != nil {
		s.CallLogHandle(netInterface.Warn, "连接接入队列已满：", err)
	}
}

// runAcceptTask 运行接入任务
func (s *Service) runAcceptTask(accept *acceptTask) {
	defer func() {
		if r := recover(); r != nil {
			s.CallLogHandle(netInterface.Error, "[websocket]连接接入处理异常:", r)
		}
	}()
	dealConn := NewConnection(s, accept.Conn, accept.ConnId, accept.Request)
	dealConn.Start()
	if s.OnConnAuth != nil && !s.OnConnAuth(dealConn, accept.Response, accept.Request) {
		dealConn.Stop()
	}
}

// runAcceptOutTime 运行接收处理任务超时
func (s *Service) runAcceptOutTime(accept *acceptTask, taskTypeCode netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) {
	s.CallLogHandle(netInterface.Error, "[ws]连接处理超时", "远程地址=>", accept.Conn.RemoteAddr(), " 任务ID=> ", accept.ConnId)
}

// SetLogHandle 设置日志处理
func (s *Service) SetLogHandle(hookFunc func(level netInterface.ErrLevel, msg ...interface{})) {
	s.onLogHandle = hookFunc
}

// GetIsStart 获取是否启动[true 启动，false 未启动]
func (s *Service) GetIsStart() bool {
	return s.isStart
}

// Start 启动
func (s *Service) Start() {
	if s.isStart {
		return
	}
	s.isStart = true

	// 开启工作线程
	s.receiveHandler.StartWorkerPool(func(errString string) {
		s.CallLogHandle(netInterface.Error, fmt.Sprint("消息处理工作池：", errString))
	})
	s.replyHandle.StartWorkerPool(func(errString string) {
		s.CallLogHandle(netInterface.Error, fmt.Sprint("消息发送工作池异常：", errString))
	})
	s.AcceptHandle.StartWorkerPool(func(errString string) {
		s.CallLogHandle(netInterface.Error, fmt.Sprint("连接接入工作池：", errString))
	})

	wsMux := http.NewServeMux() // 添加websocket 路由多路复用
	if len(s.config.PathAry) == 0 {
		// 默认监听
		wsMux.HandleFunc("/", s.wsHandler)
		//http.HandleFunc("/", s.wsHandler)
	} else {
		// 监听地址
		for _, val := range s.config.PathAry {
			wsMux.HandleFunc(fmt.Sprint("/", val), s.wsHandler)
			//http.HandleFunc(fmt.Sprint("/", val), s.wsHandler)
		}
	}
	s.startTime = time.Now()
	// 开启监听
	var wg sync.WaitGroup
	wg.Add(len(s.config.AddrAry))
	// 监听IP和端口
	for _, addr := range s.config.AddrAry {
		go func(addr string, wg *sync.WaitGroup) {
			var err error
			time.AfterFunc(time.Second*2, wg.Done)
			s.CallLogHandle(netInterface.Info, "[开启] 服务监听 [", s.config.Scheme, "]地址[", addr, "]")
			if s.config.Scheme == "wss" { // 安全连接
				server := &http.Server{
					Addr: addr,
					//cers.config.CertFile,
					//s.config.KeyFile,
					Handler: wsMux,
				}
				err = server.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
				//err = http.ListenAndServeTLS(addr, s.config.CertFile, s.config.KeyFile, nil)
			} else {
				server := http.Server{
					Addr:    addr,
					Handler: wsMux,
				}
				err = server.ListenAndServe()
				//err = http.ListenAndServe(addr, nil)
			}
			if err != nil {
				s.CallLogHandle(netInterface.Error, "[webSocket]server start listen error::", err)
			}
		}(addr, &wg)
	}
	wg.Wait()
}

// Stop 停止
func (s *Service) Stop() {
	if !s.isStart {
		return
	}
	s.isStart = false

	// 连接接入工作池
	s.AcceptHandle.StopWorkerPool()
	// 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
	// 消息处理工作池
	s.receiveHandler.StopWorkerPool()
	// 消息发送工作池
	s.replyHandle.StopWorkerPool()
}

// SetOnConnAuth 连接之前回调
func (s *Service) SetOnConnAuth(hookStart func(conn netInterface.IConnection, w http.ResponseWriter, r *http.Request) bool) {
	s.OnConnAuth = hookStart
}

// SetOnConnStart 连接完成回调
func (s *Service) SetOnConnStart(hookStart func(conn netInterface.IConnection)) {
	s.OnConnStart = hookStart
}

// SetOnConnStop 关闭之前回调
func (s *Service) SetOnConnStop(hookStop func(conn netInterface.IConnection)) {
	s.OnConnStop = hookStop
}

// SetOnReceive 数据上传完成处理函数[分包后]
func (s *Service) SetOnReceive(hookFunc func(netInterface.IConnection, int, []byte)) {
	s.onReceive = hookFunc
}

// SetOnOneReceive 【第一包数据】数据上传完成处理函数[分包后]
func (s *Service) SetOnOneReceive(hookFunc func(netInterface.IConnection, int, []byte)) {
	s.onOneReceive = hookFunc
}

// SetOnReply 数据回复完成后处理函数
func (s *Service) SetOnReply(hookFunc func(netInterface.IConnection, int, []byte, bool, string, interface{}, error)) {
	s.onReply = hookFunc
}

// SetRunTaskError 任务运行异常回调
func (s *Service) SetRunTaskError(h func(conn netInterface.IConnection, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode, task interface{})) {
	s.runTaskError = h
}

// GetReceiveWorkerPool 消息处理模块(工作池)
func (s *Service) GetReceiveWorkerPool() pools.ITaskWorkerPool {
	return s.receiveHandler
}

// GetReplyWorkerPool 消息发送处理器(工作池)
func (s *Service) GetReplyWorkerPool() pools.ITaskWorkerPool {
	return s.replyHandle
}

// GetAcceptWorkerPool 连接接收处理器(工作池)
func (s *Service) GetAcceptWorkerPool() pools.ITaskWorkerPool {
	return s.AcceptHandle
}

// CallRunTaskError 回调运行异常[taskType:0上行处理任务，1下行处理任务，task 任务对象]
func (s *Service) CallRunTaskError(conn netInterface.IConnection, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode, task interface{}) {
	if s.runTaskError != nil {
		defer func() {
			if r := recover(); r != nil {
				s.CallLogHandle(netInterface.Error, "[websocket]下发后回调业务逻辑异常:", r)
			}
		}()
		s.runTaskError(conn, taskType, taskCode, task)
	}
}

// CallOnReceive 数据上传完成回调
func (s *Service) CallOnReceive(conn netInterface.IConnection, msgType int, data []byte) {
	if s.onReceive != nil {
		s.onReceive(conn, msgType, data)
	}
}

// CallOnOneReceive [第一包]数据上传完成回调
func (s *Service) CallOnOneReceive(conn netInterface.IConnection, msgType int, data []byte) {
	if s.onOneReceive != nil {
		s.onOneReceive(conn, msgType, data)
	}
}

// CallOnReply 下发后回调
func (s *Service) CallOnReply(conn netInterface.IConnection, msgType int, data []byte, ok bool, cmdCode string, param interface{}, err error) {
	if s.onReply != nil {
		defer func() {
			if r := recover(); r != nil {
				s.CallLogHandle(netInterface.Error, "[webSocket]下发后回调业务逻辑异常:", r)
			}
		}()
		s.onReply(conn, msgType, data, ok, cmdCode, param, err)
	}
}

// CallLogHandle 错误消息处理
func (s *Service) CallLogHandle(level netInterface.ErrLevel, msgAry ...interface{}) {
	if s.onLogHandle != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Println("[webSocket]CallLogHandle 错误消息处理调用业务逻辑异常:", r)
			}
		}()
		s.onLogHandle(level, msgAry)
	}
}

// CallOnConnStart 调用连接之前
func (s *Service) CallOnConnStart(conn netInterface.IConnection) {
	if s.OnConnStart != nil {
		defer func() {
			if r := recover(); r != nil {
				s.CallLogHandle(netInterface.Error, "[webSocket]调用开始连接业务逻辑异常：", r)
			}
		}()
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用关闭之前
func (s *Service) CallOnConnStop(conn netInterface.IConnection) {
	if s.OnConnStart != nil {
		defer func() {
			if r := recover(); r != nil {
				s.CallLogHandle(netInterface.Error, "[webSocket]调用断开连接业务逻辑异常：", r)
			}
		}()
		s.OnConnStop(conn)
	}
}

// GetStartTime 获取连接启动时间
func (s *Service) GetStartTime() time.Time {
	return s.startTime
}
