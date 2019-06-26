package akeneo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

const (
	AkeneoTypeBoolean      = "pim_catalog_boolean"
	AkeneoTypeMultiSelect  = "pim_catalog_multiselect"
	AkeneoTypeSimpleSelect = "pim_catalog_simpleselect"
	AkeneoTypeNumber       = "pim_catalog_number"
	AkeneoTypeDate         = "pim_catalog_date"
	AkeneoTypeImage        = "pim_catalog_image"
	AkeneoTypeText         = "pim_catalog_text"
	AkeneoTypeTextArea     = "pim_catalog_textarea"
	AkeneoTypeFile         = "pim_catalog_file"
)

type Attribute struct {
	Code                string            `json:"code"`
	Type_               string            `json:"type"`
	Labels              map[string]string `json:"labels,omitempty"`
	Group               string            `json:"group"`
	SortOrder           int32             `json:"sort_order,omitempty"`
	Localizable         bool              `json:"localizable,omitempty"`
	Scopable            bool              `json:"scopable,omitempty"`
	AvailableLocales    []string          `json:"available_locales,omitempty"`
	Unique              bool              `json:"unique,omitempty"`
	UseableAsGridFilter bool              `json:"useable_as_grid_filter,omitempty"`
	MaxCharacters       int32             `json:"max_characters,omitempty"`
	ValidationRule      string            `json:"validation_rule,omitempty"`
	ValidationRegexp    string            `json:"validation_regexp,omitempty"`
	WysiwygEnabled      bool              `json:"wysiwyg_enabled,omitempty"`
	NumberMin           string            `json:"number_min,omitempty"`
	NumberMax           string            `json:"number_max,omitempty"`
	DecimalsAllowed     bool              `json:"decimals_allowed,omitempty"`
	NegativeAllowed     bool              `json:"negative_allowed,omitempty"`
	MetricFamily        string            `json:"metric_family,omitempty"`
	DefaultMetricUnit   string            `json:"default_metric_unit,omitempty"`
	DateMin             *time.Time        `json:"date_min,omitempty"`
	DateMax             *time.Time        `json:"date_max,omitempty"`
	AllowedExtensions   []string          `json:"allowed_extensions,omitempty"`
	MaxFileSize         string            `json:"max_file_size,omitempty"`
}

type AttributeApi ApiService

type AttributeItem struct {
	Attribute
	ResponseLinks `json:"_links"`
}

type AttributesResponse struct {
	Response
	Data struct {
		Items []AttributeItem `json:"items"`
	} `json:"_embedded"`
}

func (service *AttributeApi) GetAll(opts RequestOpts) (*AttributesResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "attributes", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &AttributesResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *AttributeApi) Get(code string) (*Attribute, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("attributes/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var attribute = &Attribute{}
	if err = json.NewDecoder(response.Body).Decode(&attribute); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return attribute, nil
}

func (service *AttributeApi) Create(attribute *Attribute) *ApiError {
	headers := service.client.getHeadersForRequest()
	body, _ := json.Marshal(attribute)

	response, err := service.client.DoRequest("POST", "attributes", headers, body, nil)
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

func (service *AttributeApi) Upsert(attribute *Attribute) *ApiError {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("attributes/%s", attribute.Code)
	body, _ := json.Marshal(attribute)

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

func (service *AttributeApi) BatchUpsert(attributes []*Attribute) ([]*ResponseBody, *ApiError) {
	headers := service.client.getHeadersForBatchRequest()
	var body []byte

	for _, bodyItem := range attributes {
		bodyItem, _ := json.Marshal(bodyItem)
		body = append(body, bodyItem...)
		body = append(body, '\n')
	}

	response, err := service.client.DoRequest("PATCH", "attributes", headers, body, nil)
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
