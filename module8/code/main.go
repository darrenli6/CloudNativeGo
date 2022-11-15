package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", myHandler)
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":8081", nil)

}

func healthz(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "ok\n")
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("enter myHandler ")

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
