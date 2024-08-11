package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	// 绑定结构体的用法
	handler := new(studentHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

// 实现ServeHTTP接口用于处理服务请求
type studentHandler struct{}

func (sh studentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	switch len(pathSegments) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh studentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentMutex.RLock()
	defer studentMutex.RUnlock()

	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	return
}

func (sh studentHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentMutex.RLock()
	defer studentMutex.RUnlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	return
}

func (sh studentHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentMutex.Lock()
	defer studentMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	return
}

// 转换成JSON的字节流返回
func (sh studentHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students: %v", err)
	}

	return b.Bytes(), nil
}
