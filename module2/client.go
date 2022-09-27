package main

import (
	"net/http"
)

func main() {

	client := &http.Client{}
	postData := make(map[string]string)

	req, err := http.NewRequest("POST", "http://localhost:8001", "")
}
