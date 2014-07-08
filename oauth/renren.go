package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	VERSION        = "0.1"
	DefaultBaseURL = "https://api.renren.com/v2/"
	AouthURL       = "https://graph.renren.com/oauth/authorize"
	TokenURL       = "https://graph.renren.com/oauth/token"
	userAgent      = "renren/" + VERSION
	GET            = "GET"
	POST           = "POST"
)

type APIError struct {
	ErrorMessage struct {
		ErrorCode string `json:"code"`
		Message   string `json:"message"`
	} `json:"error"`
}

func (err *APIError) Error() (str string) {
	return fmt.Sprintf("APIError ErrorMessage:%v ErrorCode %v", err.ErrorMessage.Message, err.ErrorMessage.ErrorCode)
}

type APIToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	expires      int64  `json:"expires_in"`
}
type RenRenUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Response struct {
	UserInfo RenRenUser `json:"response"`
}

/**
 * check ApiToken is Expired
 * @return bool
 */
func (api *APIToken) isExpired() bool {
	return api.AccessToken == "" || api.RefreshToken == "" || api.expires < time.Now().Unix()
}

/*
 * client
 */
type ApiClient struct {
	api *APIRenRen
}

/**
 * api client
 */
type APIRenRen struct {
	ApiKey         string
	ApiSecret      string
	RedirectURL    string
	ResponseType   string
	OAuth2URL      string
	AccessTokenURL string
	AccessToken    string
	Version        string
}

func NewRenRenAPI(appkey string, secret string, redirect_url string, responseType string) (api *APIRenRen) {
	api = &APIRenRen{
		ApiKey:         appkey,
		ApiSecret:      secret,
		RedirectURL:    redirect_url,
		ResponseType:   responseType,
		Version:        VERSION,
		OAuth2URL:      AouthURL,
		AccessTokenURL: TokenURL,
	}
	return api
}

func (token *APIToken) SetAccessToken(atoken string, expires int64) {
	token.AccessToken = atoken
	token.expires = expires
}

/**
 * get the oathor url
 */
func (api *APIRenRen) OAuthorURL() (uri string, err error) {
	_url, err := url.Parse(api.OAuth2URL)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	v := url.Values{
		"client_id":     {api.ApiKey},
		"redirect_uri":  {api.RedirectURL},
		"response_type": {"code"},
	}
	q := v.Encode()
	fmt.Println("qcode ", q)
	if _url.RawQuery == "" {
		_url.RawQuery = q
	} else {
		_url.RawQuery = _url.RawQuery + "&" + q
	}
	fmt.Println("url =", _url.String())
	return _url.String(), nil
}

func (api *APIRenRen) SetAccessToken(token string) {
	api.AccessToken = token
}

/**
 * renren client
 */
func (api *APIRenRen) GetAccessToken(code string) (result *APIToken, err error) {
	_url, err := url.Parse(TokenURL)
	if err != nil {
		return nil, err
	}
	q := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {api.ApiKey},
		"client_secret": {api.ApiSecret},
		"redirect_uri":  {api.RedirectURL},
		"code":          {code},
	}
	ret, err := api.sendRequest(_url, GET, q)
	if err != nil {
		return nil, err
	}
	token := APIToken{}
	err = json.Unmarshal([]byte(ret), &token)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &token, nil
}

/**
 * send net request
 */
func (api *APIRenRen) sendRequest(_url *url.URL, method string, param url.Values) (result string, err error) {
	var body io.Reader
	switch method {
	case GET:
		if _url.RawQuery == "" {
			_url.RawQuery = param.Encode()
		} else {
			_url.RawQuery = _url.RawQuery + "&" + param.Encode()
		}
		body = nil
	case POST:
		body = strings.NewReader(param.Encode())
	}
	client := new(http.Client)
	// fmt.Println(_url.String())
	request, err := http.NewRequest(method, _url.String(), body)
	if err != nil {
		return "", err
	}
	switch method {
	case POST:
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	request.Header.Add("User-Agent", "renren/"+VERSION)
	if strings.EqualFold(api.AccessToken, "") == false {
		request.Header.Add("Authorization", "Bearer "+api.AccessToken)
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	switch response.Header.Get("Conent-Encoding") {
	case "gzip":
		fmt.Println("gzip")
		// need unzip
	default:
		conent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(conent), nil
	}
	return "", nil
}

/**
 * api get request
 */
func (api *APIRenRen) ApiGet(path string, parameters map[string]string) (data string, err error) {
	upath := fmt.Sprintf("%v%v", DefaultBaseURL, path)
	_url, err := url.Parse(upath)
	if err != nil {
		return "", err
	}
	param := url.Values{}
	for k, v := range parameters {
		param.Set(k, v)
	}
	return api.sendRequest(_url, GET, param)
}

/**
 * api post request
 */
func (api *APIRenRen) ApiPost(path string, parameters map[string]string) (data string, err error) {
	upath := fmt.Sprintf("%v%v", DefaultBaseURL, path)
	_url, err := url.Parse(upath)
	if err != nil {
		return "", err
	}
	param := url.Values{}
	for k, v := range parameters {
		param.Set(k, v)
	}
	return api.sendRequest(_url, POST, param)
}

/**
 * Get UserInfo (starchou)
 */
func (api *APIRenRen) GetUserInfo() (*RenRenUser, error) {
	result, _ := api.ApiGet("user/get", nil)
	beego.Info(result)
	var resp Response
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		return nil, err
	} else {
		return &resp.UserInfo, nil
	}
}
