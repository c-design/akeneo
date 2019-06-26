package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type AttributeGroup struct {
	Code       string            `json:"code"`
	SortOrder  int32             `json:"sort_order,omitempty"`
	Attributes []string          `json:"attributes,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type AttributeGroupApi ApiService

type AttributeGroupItem struct {
	Attribute
	ResponseLinks `json:"_links"`
}

type AttributeGroupsResponse struct {
	Response
	Data struct {
		Items []AttributeGroupItem `json:"items"`
	} `json:"_embedded"`
}

func (service *AttributeGroupApi) GetAll(opts RequestOpts) (*AttributeGroupsResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "attribute-groups", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &AttributeGroupsResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *AttributeGroupApi) Get(code string) (*AttributeGroup, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("attribute-groups/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var group = &AttributeGroup{}
	if err = json.NewDecoder(response.Body).Decode(&group); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return group, nil
}

func (service *AttributeGroupApi) Create(group *AttributeGroup) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(group)

	response, err := service.client.DoRequest("POST", "attribute-groups", headers, body, nil)
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

func (service *AttributeGroupApi) Upsert(group *AttributeGroup) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("attribute-groups/%s", group.Code)
	body, _ := json.Marshal(group)

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

func (service *AttributeGroupApi) BatchUpsert(groups []*AttributeGroup) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range groups {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "attribute-groups", headers, body, nil)
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
