package akeneo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
)

type MediaFileBody struct {
	Product      *MediaFileProduct
	ProductModel *MediaFileProductModel
	FileName     string
	File         []byte
}

type MediaFile struct {
	Code             string `json:"code"`
	OriginalFilename string `json:"original_filename"`
	MimiType         string `json:"mime_type"`
	Size             int    `json:"size"`
	Extension        string `json:"extension"`
}

type MediaFileProduct struct {
	Identifier string  `json:"identifier"`
	Attribute  string  `json:"attribute"`
	Scope      *string `json:"scope"`
	Locale     *string `json:"locale"`
}

type MediaFileProductModel struct {
	Code      string  `json:"code"`
	Attribute string  `json:"attribute"`
	Scope     *string `json:"scope"`
	Locale    *string `json:"locale"`
}

type MediaFileApi ApiService

type MediaFileResponseLinks struct {
	Self     ResponseLink `json:"self"`
	Download ResponseLink `json:"download"`
}

type ProductMediaFileItem struct {
	MediaFile
	MediaFileResponseLinks `json:"_links"`
}

type ProductMediaFileResponse struct {
	Response
	Data struct {
		Items []ProductMediaFileItem `json:"items"`
	} `json:"_embedded"`
}

func (service *MediaFileApi) GetAll(opts RequestOpts) (*ProductMediaFileResponse, *ApiError) {
	headers := service.client.getHeadersForRequest()
	queryParams := &url.Values{}

	keyList := []string{"scope", "search", "locales", "attributes", "pagination_type", "page", "search_after", "limit", "withCount"}

	for _, key := range keyList {
		if value, ok := opts[key].(string); ok {
			queryParams.Add(key, value)
		}
	}

	response, err := service.client.DoRequest("GET", "media-files", headers, nil, queryParams)

	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()
	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	resp := &ProductMediaFileResponse{}

	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return resp, nil
}

func (service *MediaFileApi) Get(code string) (*MediaFile, *ApiError) {
	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("media-files/%s", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return nil, &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	var product = &MediaFile{}
	if err = json.NewDecoder(response.Body).Decode(&product); err != nil {
		return nil, &ApiError{Message: err.Error()}
	}

	return product, nil
}

func (service *MediaFileApi) Create(mediaFile *MediaFileBody) *ApiError {
	var body = &bytes.Buffer{}
	form := multipart.NewWriter(body)

	if mediaFile.Product != nil {
		productJson, err := json.Marshal(*mediaFile.Product)
		if err != nil {
			return &ApiError{Message: err.Error()}
		}

		if err := form.WriteField("product", fmt.Sprintf("%s", productJson)); err != nil {
			return &ApiError{Message: err.Error()}
		}
	} else if mediaFile.ProductModel != nil {
		productModelJson, err := json.Marshal(*mediaFile.ProductModel)
		if err != nil {
			return &ApiError{Message: err.Error()}
		}

		if err := form.WriteField("product_model", fmt.Sprintf("%s", productModelJson)); err != nil {
			return &ApiError{Message: err.Error()}
		}
	}

	file, err := form.CreateFormFile("file", mediaFile.FileName)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	reader := bytes.NewReader(mediaFile.File)
	if _, err := io.Copy(file, reader); err != nil {
		return &ApiError{Message: err.Error()}
	}

	if err := form.Close(); err != nil {
		return &ApiError{Message: err.Error()}
	}

	headers := service.client.getHeadersForRequest()
	headers.Set("Content-Type", form.FormDataContentType())

	response, err := service.client.DoRequest("POST", "media-files", headers, body.Bytes(), nil)
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

func (service *MediaFileApi) Download(code string, folderPath string) *ApiError {

	headers := service.client.getHeadersForRequest()
	uri := fmt.Sprintf("media-files/%s/download", code)

	response, err := service.client.DoRequest("GET", uri, headers, nil, nil)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		msg, _ := ioutil.ReadAll(response.Body)
		return &ApiError{Code: response.StatusCode, Status: response.Status, Message: fmt.Sprintf("%s", msg)}
	}

	imagePath := fmt.Sprintf("%s/%s", folderPath, code)

	dir := filepath.Dir(imagePath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return &ApiError{Message: err.Error()}
		}
	}

	out, err := os.Create(imagePath)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return &ApiError{Message: err.Error()}
	}

	return nil
}
