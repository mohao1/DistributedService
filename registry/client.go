package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService 用于其他服务注册服务
func RegisterService(r Registration) error {

	// 生成心跳检测的url地址
	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		return err
	}
	// 因为心跳检测逻辑简单所以直接写一个函数不需要定义单独的结构体
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 生成更新可用服务的url地址
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandle{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service responded with code %d", res.StatusCode)
	}

	return nil
}

type serviceUpdateHandle struct {
}

func (s *serviceUpdateHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	doc := json.NewDecoder(r.Body)
	var p patch
	err := doc.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Printf("Updated received %v\n", p)
	// 更新可用服务列表数据
	prov.Update(p)
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL,
		bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to dergister service. Registry service responded with code %d", res.StatusCode)
	}
	return nil
}

// Providers 每个服务的可用的所需服务的存储结构
type Providers struct {
	services map[ServiceName][]string // 存储了用的服务（需要的服务中已经注册了的服务）
	mutex    *sync.RWMutex
}

// Update 更新所需的服务的可用服务列表
func (p *Providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {

		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL) // 如果因为网络问题导致心跳一次续不上第二次连接会导致出现数重复的问题
	}

	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i, providerURL := range providerURLs {
				if providerURL == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

// 通过名称获取服务
func (p *Providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for service %s", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))

	return providers[idx], nil
}

// GetProviders 对外暴露获取服务方法
func GetProviders(name ServiceName) (string, error) {
	return prov.get(name)
}

// 创建Providers的变量
var prov = &Providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
