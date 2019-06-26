package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type ProductModel struct {
	Code          string                              `json:"code"`
	FamilyVariant string                              `json:"family_variant"`
	Parent        string                              `json:"parent,omitempty"`
	Categories    []string                            `json:"categories,omitempty"`
	Values        map[string][]*ProductAttributeValue `json:"values,omitempty"`
	Created       string                              `json:"created,omitempty"`
	Updated       string                              `json:"updated,omitempty"`
	Metadata      *ProductsMetadata                   `json:"metadata,omitempty"`
}

type ProductModelApi ApiService

type ProductModelItem struct {
	ProductModel
	ResponseLinks `json:"_links"`
}

type ProductModelResponse struct {
	Response
	Data struct {
		Items []ProductModelItem `json:"items"`
	} `json:"_embedded"`
}

func (service *ProductModelApi) GetAll(opts RequestOpts) (*ProductModelResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	keyList := []string{"scope", "search", "locales", "attributes", "pagination_type", "page", "search_after", "limit", "withCount"}

	for _, key := range keyList {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "product-models", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &ProductModelResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *ProductModelApi) Get(code string) (*ProductModel, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("product-models/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var product = &ProductModel{}
	if err = json.NewDecoder(response.Body).Decode(&product); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return product, nil
}

func (service *ProductModelApi) Create(productModel *ProductModel) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(productModel)

	response, err := service.client.DoRequest("POST", "product-models", headers, body, nil)
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

func (service *ProductModelApi) Upsert(productModel *ProductModel) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("product-models/%s", productModel.Code)
	body, _ := json.Marshal(productModel)

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

func (service *ProductModelApi) BatchUpsert(productModels []*ProductModel) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range productModels {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "product-models", headers, body, nil)
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