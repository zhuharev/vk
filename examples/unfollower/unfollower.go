package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/zhuharev/vk"
)

//WARNING
//this code remove all your follows
//it useful when you expelled from friends
func main() {
	api := &vk.Api{}
	err := api.LoginAuth(
		"login",
		"pass",
		"3396837",         // client id
		"offline,friends", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params["out"] = []string{"1"}

	strResp, err := api.Request("friends.getRequests", params)
	if err == nil {
		log.Println(string(strResp))
	}

	var ids map[string][]int

	err = json.Unmarshal([]byte(strResp), &ids)
	if err != nil {
		panic(err)
	}

	for _, v := range ids["response"] {
		api.Request("friends.delete", url.Values{
			"user_id": {fmt.Sprint(v)},
		})
	}

}
