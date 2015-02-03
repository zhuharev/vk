package main

import (
	"log"

	"github.com/yanple/vk_api"
)

func main() {
	// Login/pass auth
	var api = &vk_api.Api{}
	err := api.LoginAuth(
		"email/phone",
		"pass",
		"3087104",      // client id
		"wall,offline", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}

	// OR url auth

	//	authUrl, err := api.GetAuthUrl(
	//		"domain.com/method_get_access_token", // redirect URI
	//		"token", // response type
	//		"4672050", // client id
	//		"wall,offline", // permissions https://vk.com/dev/permissions
	//	)
	//	if err != nil {
	//		panic(err)
	//	}
	//	YourRedirectFunc(authUrl)
	//
	//	//	And receive token on the special method (redirect uri)
	//	currentUrl := getCurrentUrl() // for example "yoursite.com/get_access_token#access_token=3304fdb7c3b69ace6b055c6cba34e5e2f0229f7ac2ee4ef46dc9f0b241143bac993e6ced9a3fbc111111&expires_in=0&user_id=1"
	//	accessToken, userId, expiresIn, err := api.ParseResponseUrl(currentUrl)
	//	if err != nil {
	//		panic(err)
	//	}
	//	api.AccessToken = accessToken
	//	api.UserId = userId
	//	api.ExpiresIn = expiresIn

	// Make query
	params := make(map[string]string)
	params["domain"] = "yanple"
	params["count"] = "1"

	strResp := api.Request("wall.get", params)
	if strResp != "" {
		log.Println(strResp)
	}
}
