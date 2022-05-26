package v1

import (
	"encoding/json"
	"gitee.com/ling-bin/netwebSocket/global"
	"gitee.com/ling-bin/netwebSocket/netService"
	"gitee.com/ling-bin/netwebSocket/utils"
	"github.com/gorilla/websocket"
)

func WriteData() error {
	for message := range BroadcastChan {
		if len(message) > 0 {
			var data global.UserData
			json.Unmarshal(message, &data)
			follow, _ := utils.ConvertString2IntSlice(data.Follow)
			for _, v := range follow {
				if _, ok := global.GlobalUsers[v]; ok {
					toClient := netService.GlobalClient[v]
					err := netService.ReplyForUid(toClient, websocket.TextMessage, message, "", "", nil)
					if err != nil {
						return err
					}
				}
			}
			close(BroadcastChan)
		}
	}
	return nil
}
