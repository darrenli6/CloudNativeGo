package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/golang/glog"
)

// k8s组件大多数都是带有
func main() {

	glog.V(2).Info("start server ")
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "ok\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("enter root")

	user := r.URL.Query().Get("user")

	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s] \n", user))
	} else {
		io.WriteString(w, "hello stranger \n")
	}

	io.WriteString(w, "detail of header")

	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf(" %s=%s \n", k, v))
	}

}
