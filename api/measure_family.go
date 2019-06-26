package akeneo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type MeasureFamily struct {
	Code     string        `json:"code"`
	Standard string        `json:"standard"`
	Units    []*FamilyUnit `json:"units"`
}

type FamilyUnit struct {
	Code    string            `json:"code"`
	Symbol  string            `json:"symbol"`
	Convert map[string]string `json:"convert"`
}

type MeasureFamilyItem struct {
	MeasureFamily
	ResponseLinks `json:"_links"`
}

type MeasureFamilyResponse struct {
	Response
	Data struct {
		Items []MeasureFamilyItem `json:"items"`
	} `json:"_embedded"`
}

type MeasureFamilyApi ApiService

func (service *MeasureFamilyApi) GetAll(opts RequestOpts) (*MeasureFamilyResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	for _, key := range []string{"page", "limit", "withCount"} {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "measure-families", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &MeasureFamilyResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *MeasureFamilyApi) Get(code string) (*MeasureFamily, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("measure-families/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var successResponse = &MeasureFamily{}
	if err = json.NewDecoder(response.Body).Decode(&successResponse); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return successResponse, nil
}
