package akeneo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Locale struct {
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
}

type LocaleItem struct {
	Locale
	ResponseLinks `json:"_links"`
}

type LocaleResponse struct {
	Response
	Data struct {
		Items []LocaleItem `json:"items"`
	} `json:"_embedded"`
}

type LocaleApi ApiService


func (service *LocaleApi) GetAll(opts RequestOpts) (*LocaleResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "locales", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &LocaleResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *LocaleApi) Get(code string) (*Locale, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("locales/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var successResponse = &Locale{}
	if err = json.NewDecoder(response.Body).Decode(&successResponse); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return successResponse, nil
}