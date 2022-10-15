package main

import "sigs.k8s.io/controller-runtime/pkg/client/config"

func main() {
	conf, err := config.GetConfig
	if err != nil {
		panic(err)
		return
	}

}
