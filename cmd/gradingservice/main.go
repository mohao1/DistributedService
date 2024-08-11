package main

import (
	"DistributedService/grades"
	"DistributedService/log"
	"DistributedService/registry"
	"DistributedService/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {

	// 设定服务端口
	host, port := "127.0.0.1", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	reg := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredService:  []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeats",
	}

	// 使用Start启动服务
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		grades.RegisterHandlers)

	if err != nil {
		stlog.Fatalln(err)
	}

	if logProvider, err := registry.GetProviders(registry.LogService); err != nil {
		fmt.Printf("Logging service found at:%s\n", logProvider)
		log.SetClientLogger(logProvider, reg.ServiceName) // 调用Log的客户端去使用Log服务
	}

	// 当ctx的Done信号出现时候结束程序
	<-ctx.Done()
	fmt.Println("Shutting down log Service")
}
