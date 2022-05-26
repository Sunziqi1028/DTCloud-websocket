package router

import (
	v1 "gitee.com/ling-bin/netwebSocket/api/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/data_iot", v1.PostDataOfIot)

	return router
}
