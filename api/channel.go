package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Channel struct {
	Code            string            `json:"code"`
	Locales         []string          `json:"locales,omitempty"`
	Currencies      []string          `json:"currencies,omitempty"`
	CategoryTree    string            `json:"category_tree,omitempty"`
	ConversionUnits map[string]string `json:"conversion_units,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
}

type ChannelItem struct {
	Channel
	ResponseLinks `json:"_links"`
}

type ChannelResponse struct {
	Response
	Data struct {
		Items []ChannelItem `json:"items"`
	} `json:"_embedded"`
}

type ChannelApi ApiService

func (service *ChannelApi) GetAll(opts RequestOpts) (*ChannelResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "channels", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &ChannelResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *ChannelApi) Get(code string) (*Channel, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("channels/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var successResponse = &Channel{}
	if err = json.NewDecoder(response.Body).Decode(&successResponse); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return successResponse, nil
}

func (service *ChannelApi) Create(category *Channel) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(category)

	response, err := service.client.DoRequest("POST", "channels", headers, body, nil)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		msg, _ := ioutil.ReadAll(response.Body)
		return &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	return nil
}

func (service *ChannelApi) Upsert(category *Channel) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("channels/%s", category.Code)
	body, _ := json.Marshal(category)

	response, err := service.client.DoRequest("PATCH", uri, headers, body, nil)
	if err != nil {
		return &ApiError{Code: 0, Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	return nil
}

func (service *ChannelApi) BatchUpsert(categories []*Channel) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range categories {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "channels", headers, body, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var apiResponse []*ResponseBody
	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		var responseLine *ResponseBody
		var reader = bytes.NewReader(scanner.Bytes())
		if err = json.NewDecoder(reader).Decode(&responseLine); err != nil {
			return nil, &ApiError{Message: err.Error()}
		}

		apiResponse = append(apiResponse, responseLine)
	}

	return apiResponse, nil
}