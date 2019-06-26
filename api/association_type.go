package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type AssociationType struct {
	Code   string            `json:"code"`
	Labels map[string]string `json:"labels,omitempty"`
}

type AssociationTypeApi ApiService

type AssociationTypeItem struct {
	AssociationType
	ResponseLinks `json:"_links"`
}

type AssociationTypeResponse struct {
	Response
	Data struct {
		Items []AssociationTypeItem `json:"items"`
	} `json:"_embedded"`
}


func (service *AssociationTypeApi) GetAll(opts RequestOpts) (*AssociationTypeResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "association-types", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &AssociationTypeResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *AssociationTypeApi) Get(code string) (*AssociationType, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("association-types/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var associationType = &AssociationType{}
	if err = json.NewDecoder(response.Body).Decode(&associationType); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return associationType, nil
}

func (service *AssociationTypeApi) Create(associationType *AssociationType) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(associationType)

	response, err := service.client.DoRequest("POST", "association-types", headers, body, nil)
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

func (service *AssociationTypeApi) Upsert(associationType *AssociationType) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("association-types/%s", associationType.Code)
	body, _ := json.Marshal(associationType)

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

func (service *AssociationTypeApi) BatchUpsert(associationTypes []*AssociationType) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range associationTypes {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "association-types", headers, body, nil)
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
