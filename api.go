package vk

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	API_METHOD_URL = "https://api.vk.com/method/"
	AUTH_HOST      = "https://oauth.vk.com/authorize"

	VK_API_VERSION = "5.31"
)

var (
	DEBUG = false
	//todo change freq
	RequestFreq = 350 * time.Millisecond

	mx sync.Mutex
)

/*const (
	METHOD_ACCOUNT_GET_BANNED = "account.getBanned"
	METHOD_MESSAGES_SEND      = "messages.send"
)*/

type Api struct {
	AccessToken string
	UserId      string
	ExpiresIn   string
	debug       bool

	LastCall time.Time
}

func ParseResponseUrl(responseUrl string) (string, string, string, error) {
	u, err := url.Parse("?" + responseUrl)
	if err != nil {
		return "", "", "", err
	}

	q := u.Query()
	return q.Get("access_token"), q.Get("user_id"), q.Get("expires_in"), nil
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
		"client_id=" + client_id + "&v=" + VK_API_VERSION + "&scope=" + scope + "&display=wap"

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

func (vk *Api) Request(methodName string, p ...url.Values) ([]byte, error) {
	mx.Lock()
	defer mx.Unlock()
	u, err := url.Parse(API_METHOD_URL + methodName)
	if err != nil {
		return []byte{}, err
	}

	params := url.Values{}
	if len(p) > 0 {
		params = p[0]
	}
	params.Set("access_token", vk.AccessToken)
	params.Set("v", VK_API_VERSION)
	//u.RawQuery = params.Encode()

	tnow := time.Now()
	dur := tnow.Sub(vk.LastCall)
	if dur < RequestFreq {
		time.Sleep(RequestFreq - dur)
		if vk.debug {
			log.Printf("Slepping %s\n", dur)
		}
	}

	if vk.debug {
		log.Println(u.String())
	}
	resp, err := http.PostForm(u.String(), params)
	if err != nil {
		return []byte{}, err
	}
	vk.LastCall = tnow

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if DEBUG {
		//	log.Println(string(content))
	}
	return content, nil
}

func (vk *Api) LoginAuth(email string, password string, client_id string, scope string) error {
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
	accessToken, userId, expiresIn, err := ParseResponseUrl(res.Request.URL.Fragment)

	if vk.debug {
		log.Printf("Access token %s for user %s", accessToken, userId)
	}

	vk.AccessToken = accessToken
	vk.ExpiresIn = userId
	vk.UserId = expiresIn

	return nil
}

func (vk *Api) GetAuthUrl(redirectUri string, responseType string, client_id string, scope string) (string, error) {
	u, err := url.Parse(AUTH_HOST)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("client_id", client_id)
	q.Set("scope", scope)
	q.Set("redirect_uri", redirectUri)
	q.Set("response_type", responseType)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (vk *Api) SetDebug(s bool) {
	vk.debug = s
}
