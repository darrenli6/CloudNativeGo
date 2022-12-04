package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var loger *log.Logger

func init() {

	goapppath := os.Getenv("GOAPPPATH")
	if goapppath == "" {
		goapppath = "./"
	}

	file := goapppath + time.Now().Format("2022") + "_log" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[orcale_query]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
}

func main() {

	port := os.Getenv("GOPORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", myHandler)
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":"+port, nil)

	loger.Println("application start ", port)

}

func healthz(w http.ResponseWriter, r *http.Request) {

	loger.Println("ok ")
	io.WriteString(w, "ok\n")
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("enter myHandler ")

	loger.Println("enter myHandler ")
	// 获取环境变量

	log.Println("读取VERSION环境变量:", os.Getenv("VERSION"))
	w.Header().Add("VERSION", fmt.Sprintf("%s", os.Getenv("VERSION")))
	ip := r.Header.Get("X-Real-IP")
	log.Printf("客户端ip:%s\n", ip)

	io.WriteString(w, "detail of header \n")

	for k, v := range r.Header {
		w.Header().Add(k, fmt.Sprintf("%s", v))
		io.WriteString(w, fmt.Sprintf(" %s=%s \n", k, v))
	}

	io.WriteString(w, "well done \n")

}
