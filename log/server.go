package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

// 存储log对象函数
var log *stlog.Logger

// 作为log信息类型
type fileLog string

// 创建一个fileLog的写入的方法（将数据写入fileLog中）
// 为了保证实现io.Write的接口才可以作为log的传输类型
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(data)
}

// Run 运行并且创建log的对象（使用我们自己定义的log的数据传输类型）
func Run(destination string) {
	log = stlog.New(fileLog(destination), "go: ", stlog.LstdFlags)
}

// RegisterHandlers 创建接口处理方法(接收服务请求过来的log信息并且处理)
func RegisterHandlers() {
	// 往http中绑定处理函数
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			writer(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

// 调用自定义的log进行输出
func writer(message string) {
	log.Printf("%v\n", message)
}
