package registry

type ServiceName string

// Registration 注册服务存储的数据的结构
type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredService  []ServiceName // 注册的服务对象依赖的服务
	ServiceUpdateURL string        // 告诉注册的服务我当前的服务有哪几个匹配可用
	HeartbeatURL     string        // 用于心跳检查
}

// LogService 注册服务名称
const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	PortalService  = ServiceName("Portal")
)

// 可用服务更新数据结构(每一个的服务)
type patchEntry struct {
	Name ServiceName
	URL  string
}

// 可用服务更新数据结构(发送更新的全部的综合)
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
