package akeneo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	API_VERSION_V1 = "v1"
)

type Client struct {
	httpClient *http.Client
	auth       *ClientAuth
	baseUrl    string
	apiVersion string
	userAgent  string
}

type ClientConfig struct {
	BaseUrl        string
	UserAgent      string
	Username       string
	Password       string
	ClientId       string
	SecretKey      string
	RequestTimeout time.Duration
}

func NewClient(config *ClientConfig) *Client {

	httpClient := &http.Client{
		Timeout: config.RequestTimeout,
	}

	return &Client{
		auth: &ClientAuth{
			authURL:   config.BaseUrl + "/api/oauth/v1/token",
			clientId:  config.ClientId,
			secretKey: config.SecretKey,
			userAgent: config.UserAgent,
			body: &AuthBody{
				Username:  config.Username,
				Password:  config.Password,
				GrantType: "password",
			},
			client: httpClient,
		},
		baseUrl:    config.BaseUrl,
		apiVersion: API_VERSION_V1,
		userAgent:  config.BaseUrl,
		httpClient: httpClient,
	}
}

func (c *Client) callAPI(request *http.Request) (*http.Response, error) {
	return c.httpClient.Do(request)
}

func (c *Client) prepareRequestUrl(path string) string {
	return fmt.Sprintf("%s/api/rest/%s/%s", c.baseUrl, c.apiVersion, path)
}

func (c *Client) getHeadersForRequest() *http.Header {
	headers := &http.Header{}
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("Accept", "application/json")
	headers.Set("User-Agent", c.userAgent)

	return headers
}

func (c *Client) getHeadersForBatchRequest() *http.Header {
	headers := &http.Header{}
	headers.Set("Content-Type", "application/vnd.akeneo.collection+json")
	headers.Set("Accept", "application/json")
	headers.Set("User-Agent", c.userAgent)

	return headers
}

func (c *Client) DoRequest(method string, uri string, headers *http.Header, bodyParams []byte, queryParams *url.Values) (response *http.Response, err error) {
	var body = &bytes.Buffer{}
	uri = c.prepareRequestUrl(uri)

	if bodyParams != nil {
		_, err := body.Write(bodyParams)
		if err != nil {
			return nil, err
		}
	}

	reqUrl, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		query := reqUrl.Query()
		for k, v := range *queryParams {
			for _, iv := range v {
				query.Add(k, iv)
			}
		}

		reqUrl.RawQuery = query.Encode()
	}

	request, err := http.NewRequest(method, reqUrl.String(), body)

	if err != nil {
		return nil, err
	}

	request.Header = *headers

	token, err := c.auth.GetToken()
	if err != nil {
		return nil, err
	}

	token.SetAuthHeader(request)
	return c.callAPI(request)
}

type Token struct {
	oauth2.Token
	ExpiresIn int `json:"expires_in"`
}

type AuthBody struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GrantType string `json:"grant_type"`
}

type ClientAuth struct {
	sync.RWMutex
	authURL   string
	clientId  string
	secretKey string
	userAgent string
	body      *AuthBody
	client    *http.Client
	token     *Token
}

func (ca *ClientAuth) GetToken() (*Token, error) {
	var err error

	ca.Lock()
	defer ca.Unlock()

	if ca.token == nil || !ca.token.Valid() {
		ca.token, err = ca.receiveToken()
		if err != nil {
			return nil, err
		}
	}

	return ca.token, nil
}

func (ca *ClientAuth) receiveToken() (*Token, error) {
	jsonBody, err := json.Marshal(ca.body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ca.authURL, bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", ca.userAgent)

	req.SetBasicAuth(ca.clientId, ca.secretKey)

	resp, err := ca.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	token := &Token{}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	token.Expiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return token, nil
}