package main

import (
	"fmt"
	"net/http"
)

func main() {

	client := &http.Client{}

	// ...
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("X-Real-IP", "192.168.10.100")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("请求不正常")
		return
	}
	fmt.Println(resp)

}
