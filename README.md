Api client for VKontakte with login/pass authorization (hack) on Go (golang).
==========
###Plus: masking client_id to the iPhone, Android, iPad, Windows Phone clients.

go (golang) api client for vk.com

###Get
```Bash
    go get github.com/yanple/vk_api
    // and dependence
    go get github.com/PuerkitoBio/goquery
```

###Import
```Go
    @import "github.com/yanple/vk_api"
```

##How to use

###Login/pass auth

```Go
	var api vk_api.Api
	err := api.LoginAuth(
		"email/phone",
		"pass",
		"3087104", // client id
		"wall,offline", // scope (permissions)
	)
	if err != nil {
		panic(err)
	}
```

###By user auth (click "allow" on special vk page)
```Go
	var api vk_api.Api
	authUrl, err := api.GetAuthUrl(
		"domain.com/method_get_access_token", // redirect URI
		"token", // response type
		"4672050", // client id
		"wall,offline", // permissions https://vk.com/dev/permissions
	)
	if err != nil {
		panic(err)
	}
	YourRedirectFunc(authUrl)

	//	And receive token on the special method (redirect uri)
	currentUrl := getCurrentUrl() // for example "yoursite.com/get_access_token#access_token=3304fdb7c3b69ace6b055c6cba34e5e2f0229f7ac2ee4ef46dc9f0b241143bac993e6ced9a3fbc111111&expires_in=0&user_id=1"
	accessToken, userId, expiresIn, err := api.ParseResponseUrl(currentUrl)
	if err != nil {
		panic(err)
	}
	api.AccessToken = accessToken
	api.UserId = userId
	api.ExpiresIn = expiresIn
```

###Make query to API
```Go
	params := make(map[string]string)
	params["domain"] = "yanple"
	params["count"] = "1"

	strResp := api.Request("wall.get", params)
	if strResp != "" {
		pretty.Println(strResp)
	}
```

All api methods on https://vk.com/dev/methods

###Client ids (Masking only for login/pass auth)
```Go
    // client_id = "28909846" # Vk application ID (Android) doesn't work.
	// client_id = "3502561"  # Vk application ID (Windows Phone)
	// client_id = "3087106"  # Vk application ID (iPhone)
	// client_id = "3682744"  # Vk application ID (iPad)
```
