package netService

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"gitee.com/ling-bin/netwebSocket/netInterface"
)

// Connection 连接
type Connection struct {
	server           *Service               // 当前属于那个server
	conn             *websocket.Conn        // 当前连接的ws
	connId           uint64                 // 连接id
	isClosed         bool                   // 当前连接的关闭状态 [ true:关闭，false:开 ]
	isClosedLock     sync.Mutex             // 连接状态锁
	recPackCount     uint64                 // 上行当前处理到的包数
	recTotalByteSize uint64                 // 上行总大小(字节)
	repPackCount     uint64                 // 下行总包个数
	repTotalByteSize uint64                 // 下行总大小(字节)
	repPackErrCount  uint64                 // 下发异常包个数
	incr             int64                  // 供业务下发使用流水号
	request          *http.Request          // 请求内容
	property         map[string]interface{} // 绑定属性
	propertyLock     sync.RWMutex           // 保护连接属性
	heartTime        time.Time              // 心跳时间(每包更新)
}

// NewConnection 初始化连接方法
func NewConnection(server *Service, conn *websocket.Conn, connId uint64, request *http.Request) *Connection {
	c := &Connection{
		server:    server,
		conn:      conn,
		connId:    connId,
		isClosed:  true,
		property:  make(map[string]interface{}),
		heartTime: time.Now(),
		request:   request,
	}
	// 将当前连接放入ConnMgr
	c.server.GetConnMgr().Add(c)
	return c
}

// GetNetConn 获取连接
func (c *Connection) GetNetConn() interface{} {
	return c.conn
}

// GetNetwork 获取网络类型
func (c *Connection) GetNetwork() string {
	return c.server.config.Scheme
}

// GetConnId 获取客户端ID
func (c *Connection) GetConnId() uint64 {
	return c.connId
}

// GetRemoteAddr 获取远程客户端地址信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// GetLocalAddr 获取本地地址
func (c *Connection) GetLocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// GetHeartTime 心跳时间
func (c *Connection) GetHeartTime() time.Time {
	return c.heartTime
}

// GetRecInfo 上行当前处理的包总数（处理前，1开始），总大小(字节)
func (c *Connection) GetRecInfo() (count, byteSize uint64) {
	return c.recPackCount, c.recTotalByteSize
}

// GetRepInfo 下行当前处理的包总数（处理后），总大小(字节)
func (c *Connection) GetRepInfo() (count, byteSize, errCount uint64) {
	return c.repPackCount, c.repTotalByteSize, c.repPackErrCount
}

// GetIsClosed 获取的状态[脏读][ture:关闭状态，false:未关闭]
func (c *Connection) GetIsClosed() bool {
	return c.isClosed
}

// Incr 连接提供给业务作为流水号使用,循环累加,(val 为正数为递增，val为负数为递减,val为0则获取值)
func (c *Connection) Incr(val int64) int64 {
	return atomic.AddInt64(&c.incr, val)
}

// Start 启动连接，让当前连接，开始工作
func (c *Connection) Start() {
	if c.isClosed {
		c.isClosedLock.Lock()
		if c.isClosed {
			c.isClosed = false
			c.isClosedLock.Unlock()
			// 按照开发者传递的函数来，调用回调函数
			c.server.CallOnConnStart(c)
			// 根据官方文档 读与写只能开一个线程
			// 启动读数据业务
			go c.startReader()
		}
	}
}

// Stop 停止连接，结束当前连接工作
func (c *Connection) Stop() {
	if !c.isClosed {
		c.isClosedLock.Lock()
		if !c.isClosed {
			c.isClosed = true
			c.isClosedLock.Unlock()

			// 按照开发者传递的函数来，调用回调函数,注意在close之前调用
			c.server.CallOnConnStop(c)
			// 关闭连接
			c.conn.Close()
			// 将conn在ConnMgr中删除
			c.server.GetConnMgr().Remove(c)
		}
	}
}

// SendData 发送数据给远程的TCP客户端消息类型 1.TextMessage（文本） 2、BinaryMessage(二进制)
func (c *Connection) SendData(msgType int, data []byte, cmdCode string) error {
	return c.SendDataCall(msgType, data, cmdCode, nil, nil)
}

// SendDataCall 发送数据给远程的TCP客户端(带参数和回调)消息类型 1.TextMessage（文本） 2、BinaryMessage(二进制)
func (c *Connection) SendDataCall(msgType int, data []byte, cmdCode string, param interface{}, callFunc func(netInterface.IConnection, int, []byte, bool, string, interface{}, error)) error {
	if c.isClosed {
		return errors.New("连接关闭，不能发送消息")
	}

	// 对象池获取

	var uid uint64 = 2
	funs := c.server.GetConnMgr()

	// Service, _ = s.GetConnMgr().Get(uid)

	// c.server.config

	// service.GetConnMgr()
	//
	// count := 0
	// c.connections.Range(func(key, value interface{}) bool {
	// 	count++
	// 	return true
	// })

	println(funs)

	if uid == 1 {
		uid = 2
	}

	reply := newReplyTask()
	reply.ConnId = c.connId
	reply.Data = data
	reply.MsgType = msgType
	reply.CallFunc = callFunc
	reply.CmdCode = cmdCode
	reply.Param = param
	reply.RunReplyTask = c.runReplyTask
	reply.RunError = c.runReplyOutTime

	reply1 := newReplyTask()
	reply1.ConnId = 2
	reply1.Data = data
	reply1.MsgType = msgType
	reply1.CallFunc = callFunc
	reply1.CmdCode = cmdCode
	reply1.Param = param
	reply1.RunReplyTask = c.runReplyTask
	reply1.RunError = c.runReplyOutTime

	var err error
	if c.server.config.IsOutTime {
		reply1.Duration = c.server.config.ReplyOutTime
	} else {
		reply.Duration = 0
		reply1.Duration = 0
	}

	if !c.server.config.OverflowDiscard {
		c.server.replyHandle.SendToTaskQueueWait(reply)
		// c.server.replyHandle.SendToTaskQueueWait(reply1)
	} else {
		err = c.server.replyHandle.SendToTaskQueue(reply)
		if err != nil {
			c.server.CallLogHandle(netInterface.Warn, "发送队列已满", err)
		}
	}
	return err
}

// runReplyOutTime 回复任务运行超时
func (c *Connection) runReplyOutTime(replyTask *replyTask, taskTypeCode netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) {
	c.server.CallRunTaskError(c, taskTypeCode, taskCode, replyTask)
}

// runReplyTask 运行回复任务
func (c *Connection) runReplyTask(replyTask *replyTask) {
	defer func() {
		if r := recover(); r != nil {
			c.CallLogHandle(netInterface.Error, "[ws]运行回复任务异常:", r)
		}
	}()
	// 设置发送超时时间
	err := c.conn.SetWriteDeadline(time.Now().Add(c.server.config.SendOutTime))
	if err != nil {
		c.replyNotice(replyTask, false, err)
		return
	}
	// 发送数据
	if err := c.conn.WriteMessage(replyTask.MsgType, replyTask.Data); err != nil {
		c.replyNotice(replyTask, false, err)
		return
	}
	c.replyNotice(replyTask, true, nil)
}

// replyNotice 数据发送后通知
func (c *Connection) replyNotice(replyTask *replyTask, isOk bool, err error) {
	if replyTask.CallFunc != nil {
		replyTask.CallFunc(c, replyTask.MsgType, replyTask.Data, isOk, replyTask.CmdCode, replyTask.Param, err)
	}
	if isOk {
		c.repTotalByteSize += uint64(len(replyTask.Data))
		c.repPackCount++
	} else {
		c.repPackErrCount++
	}
	c.server.CallOnReply(c, replyTask.MsgType, replyTask.Data, isOk, replyTask.CmdCode, replyTask.Param, err)
}

// onCompleted 数据上传了完整一包的回调
func (c *Connection) onCompleted(receive *receiveTask) {
	// 延期执行
	defer func(receive *receiveTask) {
		if r := recover(); r != nil {
			c.CallLogHandle(netInterface.Error, "[ws]业务处理异常:", r)
		}
	}(receive)
	// 成功记录发送次数
	c.recPackCount++
	if c.recPackCount == 1 {
		// 第一包调用
		c.server.CallOnOneReceive(c, receive.MsgType, receive.Data)
	}
	// 更新心跳时间
	c.heartTime = time.Now()
	c.server.CallOnReceive(c, receive.MsgType, receive.Data)
}

// startReader 读取数据
func (c *Connection) startReader() {
	defer c.Stop()
	// 读业务
	for {
		// 读取数据到内存中 messageType:TextMessage/BinaryMessage
		msgType, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// 获取缓存对象
		receive := newReceiveTask()
		receive.Data = data
		receive.MsgType = msgType
		receive.OnCompleted = c.onCompleted
		receive.RunError = c.runReceiveOutTime

		// 累加上行总字节数
		c.recTotalByteSize += uint64(len(data))
		c.heartTime = time.Now()
		if c.server.config.IsOutTime {
			receive.Duration = c.server.config.ReceiveOutTime
		} else {
			receive.Duration = 0
		}
		if !c.server.config.OverflowDiscard {
			c.server.receiveHandler.SendToTaskQueueWait(receive)
		} else {
			err = c.server.receiveHandler.SendToTaskQueue(receive)
			if err != nil {
				c.CallLogHandle(netInterface.Fatal, "[TCP处理队列缓存池满]", err)
			}
		}
	}
}

// runReceiveOutTime 运行接收处理任务超时
func (c *Connection) runReceiveOutTime(receive *receiveTask, taskTypeCode netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode) {
	c.server.CallRunTaskError(c, taskTypeCode, taskCode, receive)
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("connection getProperty get error key:" + key)
	}
}

// RemoveProperty 移除设置属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

// GetPropertyKeys 获取所有属性key
func (c *Connection) GetPropertyKeys() []string {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	propertyAry := make([]string, 0, len(c.property))
	for key := range c.property {
		propertyAry = append(propertyAry, key)
	}
	return propertyAry
}

// GetRequest 获取请求体
func (c *Connection) GetRequest() *http.Request {
	return c.request
}

// CallLogHandle 调用异常处理
func (c *Connection) CallLogHandle(level netInterface.ErrLevel, msgAry ...interface{}) {
	c.server.CallLogHandle(level, msgAry)
}
