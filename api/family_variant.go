package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type FamilyVariant struct {
	Code          string                        `json:"code"`
	AttributeSets []*FamilyVariantAttributeSets `json:"variant_attribute_sets"`
	Labels        map[string]string             `json:"labels,omitempty"`
}

type FamilyVariantAttributeSets struct {
	Level      int8     `json:"level"`
	Axes       []string `json:"axes"`
	Attributes []string `json:"attributes"`
}

type FamilyVariantApi ApiService

type FamilyVariantItem struct {
	FamilyVariant
	ResponseLinks `json:"_links"`
}

type FamilyVariantsResponse struct {
	Response
	Data struct {
		Items []FamilyVariant `json:"items"`
	} `json:"_embedded"`
}

func (service *FamilyVariantApi) GetAll(code string, opts RequestOpts) (*FamilyVariantsResponse, *ApiError) {
	uri := fmt.Sprintf("families/%s/variants", code)
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", uri, headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &FamilyVariantsResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *FamilyVariantApi) Get(familyCode string, variantCode string) (*FamilyVariant, *ApiError) {
	uri := fmt.Sprintf("families/%s/variants/%s", familyCode, variantCode)
	headers := service.client.getHeadersForRequest()

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var family = &FamilyVariant{}
	if err = json.NewDecoder(response.Body).Decode(&family); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return family, nil
}

func (service *FamilyVariantApi) Create(familyCode string, variant *FamilyVariant) *ApiError {
	uri := fmt.Sprintf("families/%s/variants", familyCode)
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(variant)

	response, err := service.client.DoRequest("POST", uri, headers, body, nil)
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


func (service *FamilyVariantApi) Upsert(familyCode string, variant *FamilyVariant) *ApiError {

	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("families/%s/variants/%s", familyCode, variant.Code)
	body, _ := json.Marshal(variant)

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

func (service *FamilyVariantApi) BatchUpsert(code string, variants []*FamilyVariant) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	uri := fmt.Sprintf("families/%s/variants", code)

	var body []byte

	for _, bodyItem := range variants {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", uri, headers, body, nil)
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
