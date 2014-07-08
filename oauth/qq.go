package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	QQAuthUrl             = "https://graph.qq.com/oauth2.0"
	QQApiUrl              = "https://graph.qq.com/user"
	QQTokenUrl     string = QQAuthUrl + "/token"
	QQAuthorizeUrl string = QQAuthUrl + "/authorize"
	QQOpenIDUrl    string = QQAuthUrl + "/me"
)

type QQClient struct {
	ClientID     string
	ClientSecret string
	CallBack     string
	State        string
	OpenID       string
	AccessToken  string
}
type QQUser struct {
	Name      string `json:"nickname"`
	Figureurl string `json:"figureurl_2"`
	Gender    string `json:"gender"`
}

func (qq *QQClient) SetAccessToken(token string) {
	qq.AccessToken = token
}

func (qq *QQClient) SetOpenID(openid string) {
	qq.OpenID = openid
}

func (qq *QQClient) GetAuthorization() string {
	url_, _ := url.Parse(QQAuthorizeUrl)

	q := url.Values{
		"client_id":     {qq.ClientID},
		"redirect_uri":  {qq.CallBack},
		"response_type": {"code"},
		"state":         {qq.State},
	}.Encode()
	if url_.RawQuery == "" {
		url_.RawQuery = q
	} else {
		url_.RawQuery += "&" + q
	}

	return url_.String()
}
func (qq *QQClient) GetAccessToken(code string) (*Token, error) {
	v := url.Values{
		"client_id":     {qq.ClientID},
		"client_secret": {qq.ClientSecret},
		"redirect_uri":  {qq.CallBack},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}

	resp, err := http.PostForm(QQTokenUrl, v)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	/*
		var m struct {
			Access    string        `json:"access_token"`
			Refresh   string        `json:"refresh_token"`
			ExpiresIn time.Duration `json:"expires_in"`
		}*/
	result, _ := url.ParseQuery(string(body))
	//json.Unmarshal([]byte(body), &result)
	tok := Token{
		AccessToken:  result.Get("access_token"),
		RefreshToken: result.Get("refresh_token"),
		Expiry:       time.Now(), //time.Now().Add(result.Get("expires_in") * time.Second),
	}

	return &tok, nil
}

func (qq *QQClient) GetOpenID() (string, error) {
	q := url.Values{
		"access_token": {qq.AccessToken},
	}
	resp, err := http.PostForm(QQOpenIDUrl, q)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	//beego.Info(string(body))
	if openId, err := extractDataByRegex(string(body), `"openid":"(.*?)"`); err == nil {
		return openId, nil
	} else {
		return "", err
	}
}

func (qq *QQClient) GetUserInfo() (*QQUser, error) {
	q := url.Values{
		"oauth_consumer_key": {qq.ClientID},
		"access_token":       {qq.AccessToken},
		"openid":             {qq.OpenID},
	}
	resp, err := http.PostForm(QQApiUrl+"/get_user_info", q)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var user QQUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}
func extractDataByRegex(content, query string) (string, error) {
	rx := regexp.MustCompile(query)
	value := rx.FindStringSubmatch(content)

	if len(value) == 0 {
		return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
	}

	return strings.TrimSpace(value[1]), nil
}

//提交HTTP请求
func postRequest(callUrl string, urlvs map[string][]string) (map[string]interface{}, error) {
	var client *http.Client

	resp, err := client.PostForm(callUrl, urlvs)
	if err != nil {
		//fmt.Println("error")
		return nil, err
	}
	defer resp.Body.Close()
	body, tmpErr := ioutil.ReadAll(resp.Body)
	if tmpErr != nil {
		return nil, tmpErr
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(body), &result)
	return result, nil

}
