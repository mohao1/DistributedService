package main

import (
	"DistributedService/log"
	"DistributedService/registry"
	"DistributedService/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	// 创建自定义Log服务
	log.Run("./distributed.log")

	// 设定服务端口
	host, port := "127.0.0.1", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	reg := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddress,
		RequiredService:  make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	// 使用Start启动服务
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		log.RegisterHandlers)

	if err != nil {
		stlog.Fatalln(err)
	}

	// 当ctx的Done信号出现时候结束程序
	<-ctx.Done()

	fmt.Println("Shutting down log Service")
}
