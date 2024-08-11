package main

import (
	"DistributedService/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// 启动心跳检查服务
	registry.SetupRegistryService()

	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //结束程序

	// 配置服务
	var srv http.Server
	srv.Addr = registry.ServerPort

	// 启动服务
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		// 监听用户输入退出程序
		fmt.Printf("Registry started. Press any key to exit.\n")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel() // 使用ctx来结束main函数
	}()

	// 等待程序返回Done结束程序
	<-ctx.Done()
	fmt.Println("Shutting down registry service...")
}
