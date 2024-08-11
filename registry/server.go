package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

// 注册服务对象（注册服务）
type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

// 注册服务
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 服务注册时服务获取自己所需要的服务（心跳检测会因为这个导致服务数据获取重复）
	err := r.sendRequiredServices(reg)
	// 当服务上线时候去出发全局的通知每一个服务（其他服务感知服务上线）
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return err
}

// 取消注册服务
func (r *registry) remove(url string) error {
	for i, registration := range reg.registrations {
		if registration.ServiceURL == url {
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			// 告诉其他服务这个服务要断开连接了
			r.notify(patch{
				Added: []patchEntry{
					patchEntry{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceURL,
					},
				},
			})
			return nil
		}
	}
	return fmt.Errorf("service at URL %s not found", url)
}

// 注册时候向注册用户发送需要的服务信息（已经存在）
func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	// 循环注册服务
	for _, serviceReg := range r.registrations {
		// 循环注册服务所依赖的服务
		for _, reqService := range reg.RequiredService {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}

	// 发送可用服务列表
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// 发送服务信息
func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return nil
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	return nil
}

// 全局通告（告诉其他服务服务已经上线可以使用）
func (r *registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredService {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false

				// 判断是否存在服务的依赖服务连接
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}

				// 判断是否存在服务的依赖服务断开
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Added, removed)
						sendUpdate = true
					}
				}

				// 如果存在变化则发生变化的值给对应服务
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

// 判断心跳（检查服务是否存活）
func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				// 三次重试机会
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						// 如果被移除 了那么就重新添加回来
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("Heartbeat check failed for %v", reg.ServiceName)
					// 如果出现失败则将服务移除可用服务列表
					if success {
						success = false
						r.remove(reg.ServiceURL)
					}
					// 间隔1s进行发送一次重试
					time.Sleep(1 * time.Second)
				}
			}(reg)
			wg.Wait()
			time.Sleep(freq) // 设置成功的服务的心跳间隔时间
		}
	}
}

// 保证函数只会执行一次
var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

// 存储服务实例结构
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// RegistryService 用于启动一个服务Http的服务对外进行服务
type RegistryService struct{}

// 启动注册服务时候使用的函数（处理服务注册请求）
func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	switch r.Method {
	case http.MethodPost: // 处理注册服务
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)

		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete: // 处理取消注册服务
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		url := string(payload)
		log.Printf("Removeing service at URL: %s\n", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
