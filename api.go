package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ungerik/go-dry"
)

var (
	APIURL         = "https://api.vk.com/method/"
	OAuthURL       = "https://oauth.vk.com"
	AuthHost       = OAuthURL + "/authorize"
	AccessTokenURL = OAuthURL + "/access_token"

	DefaultRedirectURI = "https://oauth.vk.com/blank.html"

	DefaultUserAgent = `VKAndroidApp/4.8.3-1113 (Android 4.4.2; SDK 10; armeabi-v7a; ZTE ZTE Blade L3; ru)`

	// VkAPIVersion is lastest vk.com api version supported
	VkAPIVersion = "5.64"
)

var (
	// RequestFreq is request limiter
	RequestFreq = 333 * time.Millisecond

	mx sync.Mutex
)

// Api main struct
type Api struct {
	AccessToken string
	UserId      string
	ExpiresIn   string
	debug       bool

	PhoneCode string

	ClientId     int
	ClientSecret string

	StdCaptcha bool

	Lang  string
	Https bool

	cacheDir string
	cache    bool

	LastCall time.Time

	httpClient *http.Client

	CaptchaResolver CaptchaResolver
}

func NewApi(at string) *Api {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 50,
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	return &Api{AccessToken: at, httpClient: client}
}

func (a *Api) SetHTTPClient(client *http.Client) {
	a.httpClient = client
}

func (a *Api) HTTPClient() *http.Client {
	if a.httpClient != nil {
		return a.httpClient
	}
	tr := &http.Transport{
		MaxIdleConnsPerHost: 50,
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	a.httpClient = client
	return a.httpClient
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
	_origin, exists := doc.Find("form input[name=_origin]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not _origin attr in vk form")
	}

	ip_h, exists := doc.Find("form input[name=ip_h]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not ip_h attr in vk form")
	}

	lg_h, exists := doc.Find("form input[name=lg_h]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not lg_h attr in vk form")
	}

	to, exists := doc.Find("form input[name=to]").Attr("value")
	if exists == false {
		return nil, "", errors.New("Not 'to' attr in vk form")
	}

	formData := url.Values{}
	formData.Add("_origin", _origin)
	formData.Add("ip_h", ip_h)
	formData.Add("to", to)
	formData.Add("lg_h", lg_h)

	url, exists := doc.Find("#vk_wrap #m #mcont .pcont .form_item form").Attr("action")
	if exists == false {
		return nil, "", errors.New("Not action attr in vk form")
	}
	return formData, url, nil
}

func parse_perm_form(doc *goquery.Document) (string, error) {
	url, exists := doc.Find("form").Attr("action")
	if exists == false {
		return "", errors.New("Not action attr in vk form")
	}
	return url, nil
}

func auth_user(email string, password string, client_id string, scope string, client *http.Client) (*http.Response, error) {
	var authURL = "http://oauth.vk.com/oauth/authorize?" +
		"redirect_uri=http://oauth.vk.com/blank.html&response_type=token&" +
		"client_id=" + client_id + "&v=" + VkAPIVersion + "&scope=" + scope + "&display=wap"
	log.Printf("Login url = %s\n", authURL)

	res, e := client.Get(authURL)
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

func (vk *Api) get_permissions(response *http.Response, client *http.Client) (*http.Response, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	url, err := parse_perm_form(doc)
	if err != nil {
		return nil, err
	}

	if vk.debug {
		log.Printf("Get permissions = %s\n", url)
	}
	res, err := client.PostForm(url, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

var secCheckRe = regexp.MustCompile(`var params = {code: ge\('code'\).value, to: '', al_page: '', hash: '([0-9a-zA-Z]*)'};`)

func (vk *Api) securityCheck(response *http.Response, client *http.Client) (*http.Response, error) {
	if vk.debug {
		log.Printf("Security check with code = %s\n", vk.PhoneCode)
	}
	arr := secCheckRe.FindStringSubmatch("var params = {code: ge('code').value, to: '', al_page: '', hash: '7595037709139067db'};")
	if len(arr) != 2 {
		return nil, errors.New("Unknown error")
	}
	form := url.Values{}
	form.Set("code", vk.PhoneCode)
	form.Set("hash", arr[1])

	return client.PostForm("https://vk.com/login.php?act=security_check", form)
}

func (vk *Api) Request(methodName string, p ...url.Values) ([]byte, error) {
	mx.Lock()
	if vk.LastCall.IsZero() {
		vk.LastCall = time.Now()
	}

	params := url.Values{}
	if len(p) > 0 {
		params = p[0]
	}

	if vk.Lang == "" {
		vk.Lang = "en"
	}

	var pars, _ = url.QueryUnescape(params.Encode())
	if len(pars) > 80 {
		pars = pars[:80]
	}

	params.Set("lang", vk.Lang)
	if vk.Https {
		params.Set("https", "1")
	}
	if tok := params.Get("access_token"); tok == "" {
		params.Set("access_token", vk.AccessToken)
	}
	params.Set("v", VkAPIVersion)
	//u.RawQuery = params.Encode()

	if vk.cache {
		key := methodName + "?" + params.Encode()
		md5 := dry.StringMD5Hex(key)
		fname := vk.cacheDir + "/" + md5
		if !dry.FileExists(vk.cacheDir) {
			os.MkdirAll(vk.cacheDir, 0777)
		}
		if dry.FileExists(fname) {
			mx.Unlock()
			return ioutil.ReadFile(fname)
		}
	}

	dur := time.Since(vk.LastCall)
	if dur < RequestFreq /*&&  methodName != METHOD_USERS_GET*/ {
		time.Sleep(RequestFreq - dur)
		if vk.debug {
			log.Printf("Slepping %s\n", RequestFreq-dur)
		}
	}
	vk.LastCall = time.Now()
	mx.Unlock()

	if vk.debug {
		log.Println(methodName, pars)
	}
	content, err := vk.request(methodName, params)
	if err != nil {
		return nil, err
	}

	if vk.debug {
		log.Println(string(content))
	}

	//handle error
	if content[2] == 'e' {
		var er ErrResponse
		err = json.Unmarshal(content, &er)
		if err != nil {
			return nil, err
		}
		if er.Error.Code != 0 && er.Error.Code != 14 {
			return nil, fmt.Errorf("%s", er.Error.Msg)
		} else if er.Error.Code == 14 {
			var code string
			if vk.CaptchaResolver != nil {
				code, err = vk.CaptchaResolver.ResolveCaptcha(er.Error.CaptchaSid, er.Error.CapthchaImg)
				if err != nil {
					return nil, err
				}
			} else if vk.StdCaptcha {
				fmt.Printf("Open in your browser %s\n", er.Error.CapthchaImg)
				fmt.Println("Write captcha key here")
				fmt.Scanln(&code)
			} else {
				return nil, fmt.Errorf("captcha needed")
			}
			params.Set("captcha_sid", er.Error.CaptchaSid)
			params.Set("captcha_key", code)
			content, err = vk.request(methodName, params)
			if err != nil {
				return nil, err
			}
		}
		if er.Error.Code == 14 && vk.StdCaptcha {

		}
	}

	if vk.cache {
		key := methodName + "?" + params.Encode()
		md5 := dry.StringMD5Hex(key)
		fname := vk.cacheDir + "/" + md5
		return content, ioutil.WriteFile(fname, content, 0777)
	}

	mx.Lock()
	vk.LastCall = time.Now()
	mx.Unlock()
	return content, nil
}

type ErrResponse struct {
	Error struct {
		Code        int    `json:"error_code"`
		Msg         string `json:"error_msg"`
		CaptchaSid  string `json:"captcha_sid"`
		CapthchaImg string `json:"captcha_img"`
	} `json:"error"`
}

func (a *Api) request(methodName string, p url.Values) ([]byte, error) {
	u, err := url.Parse(APIURL + methodName)
	if err != nil {
		mx.Unlock()
		return nil, err
	}

	// android
	sig := Sig(p, methodName, "hHbZxrka2uZ6jB1inYsH")
	p.Set("sig", sig)
	var (
		req *http.Request
	)

	if a.debug {
		log.Println(u.String() + "?" + p.Encode())
	}

	//if len(p.Encode()) > 1024 {
	req, err = http.NewRequest("POST", u.String(), strings.NewReader(p.Encode()))
	//} else {
	//	req, err = http.NewRequest("POST", u.String()+"?"+p.Encode(), nil)
	//}
	if err != nil {
		if a.debug {
			log.Println(err)
		}
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("X-Get-Processing-Time", "1")
	req.Header.Set("Cookie", "remixlang=0")

	resp, err := a.HTTPClient().Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
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

	if vk.debug {
		log.Printf("Path %s (%s)\n", res.Request.URL.Path, res.Request.URL.Path)
	}

	if res.Request.URL.Path != "/blank.html" {
		if res.Request.URL.Query().Get("act") == "security_check" {
			_, err = vk.securityCheck(res, client)
			if err != nil {
				return err
			}
			res, err = auth_user(email, password, client_id, scope, client)
			if err != nil {
				return err
			}
		} else {
			res, err = vk.get_permissions(res, client)
			if err != nil {
				return err
			}

			//if res.Request.URL.Path != "/blank.html" {
			//	return errors.New("Not auth (path not /blank.html), path: " + res.Request.URL.Path)
			//}
		}
	}
	accessToken, userId, expiresIn, err := ParseResponseUrl(res.Request.URL.Fragment)

	if err != nil {
		log.Println(err)
	}

	if vk.debug {
		log.Printf("Access token %s for user %s", accessToken, userId)
	}

	vk.AccessToken = accessToken
	vk.ExpiresIn = userId
	vk.UserId = expiresIn

	return nil
}

// GetAuthURL return url to auth in VK
func GetAuthURL(redirectURI string, responseType string, clientID string, scope string) (string, error) {
	u, err := url.Parse(AuthHost)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("scope", scope)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", responseType)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// GetAccessTokenByCode return's access_token with authcode flow https://vk.com/dev/authcode_flow_user
func GetAccessTokenByCode(code string, clientID string, clientSecret string, redirectURI string) (accessToken string, userID int64, err error) {
	params := url.Values{}
	params.Set("code", code)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("redirect_uri", redirectURI)

	uri := AccessTokenURL + "?" + params.Encode()

	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	type TokenResponse struct {
		Token  string `json:"access_token"`
		UserID int64  `json:"user_id"`
	}

	var tr = new(TokenResponse)

	jDec := json.NewDecoder(resp.Body)
	err = jDec.Decode(tr)
	if err != nil {
		return
	}

	accessToken = tr.Token
	userID = tr.UserID

	return
}

func (vk *Api) CacheDir(s string) {
	vk.cacheDir = s
	if vk.cacheDir != "" {
		vk.cache = true
	} else {
		vk.cache = false
	}

	if vk.debug {
		log.Println("Cache", vk.cache)
	}
}

func (vk *Api) SetDebug(s bool) {
	if s {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}
	vk.debug = s
}
