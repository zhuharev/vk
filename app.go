package main

import (
	"log"
	"sync"
	"runtime"

	"github.com/kr/pretty"

	"vk_test/vk_api"
)


func main() {
	runtime.GOMAXPROCS(4)

	var api vk_api.Api
	err := api.Auth(
		"duke565@mail.ru",
		"575HPDP55",
		"3087104",
		"wall,offline",
	)
	if err != nil {
		log.Println(err)
	}

	params := make(map[string]string)
	params["domain"] = "happierall"
	params["count"] = "1"

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		strResp := api.Request("wall.get", params)
		if strResp != "" {
			pretty.Println("first okey")
		}
	}()
	go func() {
		defer wg.Done()
		strResp := api.Request("wall.get", params)
		if strResp != "" {
			pretty.Println("second okey")
		}
	}()

	wg.Wait()

	pretty.Println("complete programm")
}
