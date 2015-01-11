Api wrapper for VKontakte and login/pass authorization (hack).
Plus: masking client_id to the iPhone, Android, iPad, Windows Phone clients.
==========

go (golang) api client for vk.com

#How to use

```Go
var api vk_api.Api
err := api.Auth(
    "email",
    "pass",
    "3087104", // client id (this iphone)
    "wall,offline", // scope
)
if err != nil {
    log.Println(err)
}

// get one post
params := make(map[string]string)
params["domain"] = "happierall"
params["count"] = "1"

strResp := api.Request("wall.get", params)
pretty.Println("first okey")

//you will get string in json format that can be parsed with any json lib
stringResponse := api.Request("getProfiles", m) //{"response":[{"uid":1,"first_name":"Pavel","last_name":"Durov"}]}
```

you can find all api methods on https://vk.com/dev/methods