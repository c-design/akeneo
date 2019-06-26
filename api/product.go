package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Product struct {
	Identifier   string                              `json:"identifier"`
	Enabled      bool                                `json:"enabled"`
	FamilyCode   string                              `json:"family,omitempty"`
	Categories   []string                            `json:"categories,omitempty"`
	Groups       []string                            `json:"groups,omitempty"`
	Parent       string                              `json:"parent,omitempty"`
	Values       map[string][]*ProductAttributeValue `json:"values,omitempty"`
	Associations map[string]*ProductAssociation      `json:"associations,omitempty"`
	Created      string                              `json:"created,omitempty"`
	Updated      string                              `json:"updated,omitempty"`
	Metadata     *ProductsMetadata                   `json:"metadata,omitempty"`
}

type ProductAttributeValue struct {
	Scope  *string     `json:"scope"`
	Locale *string     `json:"locale"`
	Data   interface{} `json:"data"`
}

type ProductAssociation struct {
	Groups        []string `json:"groups,omitempty"`
	Products      []string `json:"products,omitempty"`
	ProductModels []string `json:"products,omitempty"`
}

type ProductsMetadata struct {
	WorkflowStatus string `json:"workflow_status,omitempty"`
}

type ProductApi ApiService

type ProductItem struct {
	Product
	ResponseLinks `json:"_links"`
}

type ProductsResponse struct {
	Response
	Data struct {
		Items []ProductItem `json:"items"`
	} `json:"_embedded"`
}

func (service *ProductApi) GetAll(opts RequestOpts) (*ProductsResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	keyList := []string{"page", "limit", "withCount", "scope", "search", "locales", "attributes", "pagination_type", "search_after"}

	for _, key := range keyList {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "products", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &ProductsResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *ProductApi) Get(code string) (*Product, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("products/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var product = &Product{}
	if err = json.NewDecoder(response.Body).Decode(&product); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return product, nil
}

func (service *ProductApi) Create(product *Product) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(product)

	response, err := service.client.DoRequest("POST", "products", headers, body, nil)
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

func (service *ProductApi) Upsert(product *Product) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("products/%s", product.Identifier)
	body, _ := json.Marshal(product)

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

func (service *ProductApi) BatchUpsert(products []*Product) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range products {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "products", headers, body, nil)
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

func (service *ProductApi) Delete(code string) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("products/%s", code)

	response, err := service.client.DoRequest("DELETE", uri, headers, nil, nil)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	return nil
}
