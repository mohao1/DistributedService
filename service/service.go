package service

import (
	"DistributedService/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start 启动函数用于统一启动服务
func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	// 绑定对于路由处理函数
	registerHandlersFunc()
	// 启动服务
	ctx = startService(ctx, reg.ServiceName, host, port)
	// 注册服务
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {

	// 为ctx上下文增加取消功能
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server

	srv.Addr = host + ":" + port

	go func() {
		// 监听启动服务然后监听是否出错（如果服务启动失败则返回err）
		log.Println(srv.ListenAndServe()) //srv.ListenAndServe()启动服务
		err := registry.ShutdownService(fmt.Sprintf("http://%s", host+":"+port))
		if err != nil {
			log.Println(err)
		}
		cancel() // 使用ctx来结束main函数
	}()

	go func() {
		// 监听用户输入退出程序
		fmt.Printf("%v started. Press any key to exit.\n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		err := registry.ShutdownService(fmt.Sprintf("http://%s", host+":"+port))
		if err != nil {
			log.Println(err)
		}
		cancel() // 使用ctx来结束main函数
	}()

	return ctx
}
