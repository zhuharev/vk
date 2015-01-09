package vk_api

import (
	"errors"
	"strings"
	"io/ioutil"
	"net/url"
	"net/http"
	"net/http/cookiejar"

	"github.com/PuerkitoBio/goquery"
)

const API_METHOD_URL = "https://api.vk.com/method/"

type Api struct {
	AccessToken string
	UserId      string
	ExpiresIn   string
}

func ParseResponseUrl(responseUrl string) (string, string, string) {
	u, err := url.Parse(strings.Replace(responseUrl, "#", "?", 1))
	if err != nil {
		panic(err)
	}

	q := u.Query()
	return q.Get("access_token"), q.Get("user_id"), q.Get("expires_in")
}

func parse_form(doc *goquery.Document) (url.Values, string, error) {
	_origin, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=_origin]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not _origin attr in vk form")
	}

	ip_h, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=ip_h]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not ip_h attr in vk form")
	}

	to, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form input[name=to]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not 'to' attr in vk form")
	}

	formData := url.Values{}
	formData.Add("_origin", _origin)
	formData.Add("ip_h", ip_h)
	formData.Add("to", to)

	url, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form").Attr("action")
	if exists == false {
		return nil, "", errors.New("Not action attr in vk form")
	}
	return formData, url, nil
}

func auth_user(email string, password string, client_id string, scope string, client *http.Client) (*http.Response, error) {
	var auth_url = "http://oauth.vk.com/oauth/authorize?" +
			"redirect_uri=http://oauth.vk.com/blank.html&response_type=token&" +
			"client_id=" + client_id + "&v=5.0&scope=" + scope + "&display=wap"

	res, e := client.Get(auth_url)
	if e != nil {
		return nil, e
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	formData, url, err := parse_form(doc)
	if err != nil {
		return nil, err
	}
	formData.Add("email", email)
	formData.Add("pass", password)

	res, e = client.PostForm(url, formData)
	if e != nil {
		return nil, e
	}
	return res, nil
}

func get_permissions(response *http.Response, client *http.Client) (*http.Response, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	formData, url, err := parse_form(doc)
	if err != nil {
		return nil, err
	}

	res, err := client.PostForm(url, formData)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vk Api) Request(methodName string, params map[string]string) string {
	u, err := url.Parse(API_METHOD_URL + methodName)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("access_token", vk.AccessToken)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(content)
}


func (vk Api) Auth(email string, password string, client_id string, scope string) error{
	// client_id = "28909846" # Vk application ID (Android) doesn't work.
	// client_id = "4672050"  # my APP (yanple)
	// client_id = "3502561"  # Vk application ID (Windows Phone)
	// client_id = "3087106"  # Vk application ID (iPhone)
	// client_id = "3682744"  # Vk application ID (iPad)

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	res, err := auth_user(email, password, client_id, scope, client)
	if err != nil {
		return err
	}

	if res.Request.URL.Path != "/blank.html" {
		res, err = get_permissions(res, client)
		if err != nil {
			return err
		}

		if res.Request.URL.Path != "/blank.html" {
			return errors.New("Not auth")
		}
	}

	fragment, err := url.ParseQuery(res.Request.URL.Fragment)

	vk.AccessToken = fragment["access_token"][0]
	vk.ExpiresIn = fragment["expires_in"][0]
	vk.UserId = fragment["user_id"][0]

	return nil
}

