package main

import (
	"encoding/json"
	"fmt"
	"github.com/sisteamnik/vk_api"
	"log"
)

//WARNING
//this code remove all your follows
//it useful when you expelled from friends
func main() {
	api := &vk_api.Api{}
	err := api.LoginAuth(
		"login",
		"pass",
		"3396837",         // client id
		"offline,friends", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}
	params := make(map[string]string)
	params["out"] = "1"

	strResp := api.Request("friends.getRequests", params)
	if strResp != "" {
		log.Println(strResp)
	}

	var ids map[string][]int

	err = json.Unmarshal([]byte(strResp), &ids)
	if err != nil {
		panic(err)
	}

	for _, v := range ids["response"] {
		api.Request("friends.delete", map[string]string{
			"user_id": fmt.Sprint(v),
		})
	}

}
