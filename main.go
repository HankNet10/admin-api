package main

import (
	"myadmin/router"
	"net/http"
	"os"
	"time"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
func main() {
	// 设置运行环境变量
	ginEngine := router.InitAdminRouter()
	// err := ginEngine.Run()
	// if err != nil {
	// 	panic(err)
	// }
	server := &http.Server{
		Addr:           os.Getenv("GIN_LISTEN"),
		Handler:        ginEngine,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1048576,
	}

	server.ListenAndServe()

}
