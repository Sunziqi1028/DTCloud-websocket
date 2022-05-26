package v1

import (
	"encoding/json"
	"fmt"
	"gitee.com/ling-bin/netwebSocket/global"
	"github.com/gin-gonic/gin"
)

var BroadcastChan chan []byte

//ws://127.0.0.1:8001/chat?uid=1&partner_id=1&company_id=1&name=张三&follow=1,2&type=orient
func PostDataOfIot(c *gin.Context) {
	var userData global.UserData

	c.ShouldBind(&userData)

	fmt.Println("data_controller.go line:22", userData)

	data, _ := json.Marshal(userData)

	BroadcastChan = make(chan []byte, 10)
	BroadcastChan <- data

	err := WriteData()
	if err != nil {
		fmt.Println("data_controller.go line:45", err)
		c.JSON(500, "send data failed")
	}
	c.JSON(200, "send data successful~")
	return
}
