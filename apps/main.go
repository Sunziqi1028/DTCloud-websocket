package main

import (
	"bytes"
	"fmt"
	"gitee.com/ling-bin/netwebSocket/server"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/websocket"

	"gitee.com/ling-bin/netwebSocket/netInterface"
	"gitee.com/ling-bin/netwebSocket/netService"
)

func main() {
	fmt.Println("============================================")
	fmt.Println("|                                          |")
	fmt.Println("|        中亿丰websocket物联网平台！       |")
	fmt.Println("|             微信：amoserp                |")
	fmt.Println("|                                          |")
	fmt.Println("============================================")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	port := []string{"32771"} // 监听端口
	addrAry := make([]string, 0, 5)
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				for _, val := range port {
					addrAry = append(addrAry, ipnet.IP.String()+":"+val)
				}
			}
		}
	}
	for _, val := range port {
		addrAry = append(addrAry, fmt.Sprint("127.0.0.1:", val))
	}
	config := netService.DefaultConfig("ws", addrAry, nil)
	service := netService.NewService(config)

	// 连接验证
	service.SetOnConnAuth(func(conn netInterface.IConnection, w http.ResponseWriter, r *http.Request) bool {
		return true
	})

	// 连接启动
	service.SetOnConnStart(func(connection netInterface.IConnection) {
		connId := connection.GetConnId()
		fmt.Println("[", connection.GetRemoteAddr().String(), "]新连接：", connId)
	})
	// 连接关闭
	service.SetOnConnStop(func(connection netInterface.IConnection) {
		connId := connection.GetConnId()
		fmt.Println("[", connection.GetRemoteAddr().String(), "]连接关闭：", connId)
	})

	// 内部任务运行异常
	service.SetRunTaskError(func(conn netInterface.IConnection, taskType netInterface.TaskTypeCode, taskCode netInterface.TaskErrCode, task interface{}) {
		fmt.Println("内部任务运行异常！:", taskType, "  ", taskCode)
	})

	// 第一包数据[调用后会调用SetOnReceive，在同一协程执行]
	service.SetOnOneReceive(func(connection netInterface.IConnection, msgType int, data []byte) {
		connId := connection.GetConnId()
		fmt.Println("[", connection.GetRemoteAddr().String(), " ", connection.GetLocalAddr().String(), "][", connId, "]连接上传第一次上数据：", len(data), "---", string(data))
	})
	// 接收到的数据
	service.SetOnReceive(func(connection netInterface.IConnection, msgType int, data []byte) {
		// time.Sleep(time.Second*15)
		//connId := connection.GetConnId()
		//
		//
		// fmt.Println("======")
		//fmt.Println("当前用户uid:", connId, "main.go line:75")
		// fmt.Println("[", connection.GetRemoteAddr().String(), " ", connection.GetLocalAddr().String(), "][", connId, "]客户端路径[", connection.GetRequest().URL, "]连接上传新数据：", len(data), "---", string(data))
		// //回复
		by := bytes.Buffer{}
		// by.Write([]byte("收到 => "))
		by.Write(data)
		connection.SendData(websocket.TextMessage, by.Bytes(), "")
	})
	// 数据回复成功
	service.SetOnReply(func(connection netInterface.IConnection, msgType int, data []byte, isOk bool, cmdCode string, param interface{}, err error) {
		connId := connection.GetConnId()
		fmt.Println("[下行回调][", connection.GetRemoteAddr().String(), "][", connId, "]：", len(data), "---", string(data))
		if !isOk {
			log.Println("数据下发异常：", err)
		}
	})
	// 日志
	service.SetLogHandle(func(level netInterface.ErrLevel, msg ...interface{}) {
		switch level {
		case netInterface.Fatal:
			log.Panicln("[致命]", msg)
			break
		case netInterface.Error:
			log.Println("[错误]", msg)
			break
		case netInterface.Warn:
			log.Println("[警告]", msg)
			break
		case netInterface.Info:
			log.Println("[消息]", msg)
			break
		}
	})
	service.Start()
	server.HttpStart() // 启动http服务

	go func() {
		for {
			receiveWorkerPool := service.GetReceiveWorkerPool()
			queueCount, totalCount, timeCount := receiveWorkerPool.GetTaskPoolInfo()
			replyWorkerPool := service.GetReplyWorkerPool()
			rqueueCount, rtotalCount, rtimeCount := replyWorkerPool.GetTaskPoolInfo()
			log.Println(
				"启动时间:", service.GetStartTime().Format("1982-05-21 00:00:00"),
				" 连接数:", service.GetConnMgr().Count(),
				" [接收]处理数：", totalCount,
				" 剩余数:", queueCount,
				" 处理超时数:", timeCount,
				" [回复]处理数:", rtotalCount,
				" 剩余数:", rqueueCount,
				" 处理超时数:", rtimeCount,
				" 协程数：", runtime.NumGoroutine())
			time.Sleep(time.Second * 10)
		}
	}()
	select {}
}
