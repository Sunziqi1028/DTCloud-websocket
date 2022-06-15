package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/ling-bin/netwebSocket/api/db"
	"gitee.com/ling-bin/netwebSocket/global"
	"gitee.com/ling-bin/netwebSocket/netService"
	"gitee.com/ling-bin/netwebSocket/utils"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

func WriteData() error {
	for message := range BroadcastChan {
		if len(message) > 0 {
			var data global.UserData
			json.Unmarshal(message, &data)
			follow, _ := utils.ConvertString2IntSlice(data.Follow)
			for _, v := range follow {
				if _, ok := global.GlobalUsers[v]; ok {
					uuids := global.OneUId2Uuids[v]
					for _, uuid := range uuids {
						toClient := netService.GlobalClient[uuid]
						err := netService.ReplyForUid(toClient, websocket.TextMessage, message, "", "", nil)
						if err != nil {
							return err
						}
					}

				}
			}
			close(BroadcastChan)
		}
	}
	return nil
}

// 校验前端发送数据是否含有params参数
func CheckDataIsParams(data []byte) *global.Params {
	j, err := simplejson.NewJson(data)
	if err != nil {
		log.Println("NewJson error:%s", err.Error())
		return nil
	}
	var params global.Params

	paramsNode, ok := j.CheckGet("params")
	if ok {
		bytes, err := paramsNode.MarshalJSON()
		if err != nil {
			fmt.Println(err, "client.go line:46")
		}
		err = json.Unmarshal(bytes, &params)
		if err != nil {
			fmt.Println("json unmarshal err:", err)
		}
		fmt.Println(params, "client.go line:58")
		return &params
	}
	return nil
}

func HandlerParams(params *global.Params) ([]byte, error) {
	db := db.DBEngine{
		global.PostgreSqlDBEngine,
	}
	switch params.Type {
	case global.LOG:
		s := HandLog(params)
		return []byte(s), nil

	case global.SQL:
		tableName := utils.ConvertTableName(params.Model)
		ids, _ := utils.ConvertString2IntSlice(params.Ids)
		fmt.Println(ids, "client.go line: 75")
		if params.Method == global.WRITE {
			err := db.Update(tableName, ids, params.Field)
			if err != nil {
				return nil, errors.New("更新失败！")
			}
		}
		if params.Method == global.CREATE {
			err := db.Insert(tableName, params.Field)
			if err != nil {
				return nil, errors.New("插入失败！")
			}
		}

	case global.URL:
		data, err := PostParams(params)
		if err == nil {
			return data, nil
		}
	case global.XML_RPC:

	}

	return nil, nil
}

func PostParams(params *global.Params) ([]byte, error) {
	byteTmp, err := json.Marshal(params)
	if err != nil {
		fmt.Println("json marshal err", err)
		return nil, err
	}

	resp, err := http.Post(params.URL, "application/json;charset=utf-8", bytes.NewBuffer(byteTmp))
	if err != nil {
		fmt.Printf("POST %s err:%v", params.URL, err)
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read response body err:", err)
		return nil, err
	}

	return data, nil
}

func HandLog(params *global.Params) string {
	logPath, err := utils.MkdirLogDir()
	if err != nil {
		return ""
	}
	f, err := utils.CreateLogFile(logPath)
	data, _ := json.Marshal(params)

	f.Write([]byte(data))
	return "写入成功！"
}
