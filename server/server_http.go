package server

import (
	"gitee.com/ling-bin/netwebSocket/router"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func HttpStart() {
	// 添加HTTP路由多路复用
	newRouter := router.NewRouter()

	httpServer := &http.Server{
		Addr:    ":32770",
		Handler: newRouter,
	}
	var g errgroup.Group

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
}
