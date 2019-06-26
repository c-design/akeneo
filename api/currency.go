package akeneo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Currency struct {
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
}

type CurrencyItem struct {
	Currency
	ResponseLinks `json:"_links"`
}

type CurrencyResponse struct {
	Response
	Data struct {
		Items []CurrencyItem `json:"items"`
	} `json:"_embedded"`
}

type CurrencyApi ApiService

func (service *CurrencyApi) GetAll(opts RequestOpts) (*CurrencyResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "currencies", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &CurrencyResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *CurrencyApi) Get(code string) (*Currency, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("currencies/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var successResponse = &Currency{}
	if err = json.NewDecoder(response.Body).Decode(&successResponse); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return successResponse, nil
}