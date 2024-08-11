package main

import (
	"DistributedService/log"
	"DistributedService/portal"
	"DistributedService/registry"
	"DistributedService/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	// 导入模板
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "127.0.0.1", "5000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	reg := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL:  serviceAddress,
		RequiredService: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeats",
	}

	ctx, err := service.Start(context.Background(),
		host,
		port,
		reg,
		portal.RegisterHandlers)

	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProviders(registry.LogService); err != nil {
		fmt.Printf("Logging service found at:%s\n", logProvider)
		log.SetClientLogger(logProvider, reg.ServiceName) // 调用Log的客户端去使用Log服务
	}

	// 当ctx的Done信号出现时候结束程序
	<-ctx.Done()
	fmt.Println("Shutting down log Service")
}
