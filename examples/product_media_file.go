package main

import (
	"../api"
	"fmt"
	"io/ioutil"
	"log"
)

var productMediaFileA = &akeneo.MediaFileBody{
	Product:      &akeneo.MediaFileProduct{
		Identifier: productB.Identifier,
		Attribute:  attrImage.Code,
	},
	FileName:     "logotype.jpg",
	File:         nil,
}

func runProductMediaFilesMethods() {
	createProductMediaFile()
	getAllProductMediaFile()
	getProductMediaFile()
	downloadMediaFile()
}

func createProductMediaFile() {
	var err error

	if productMediaFileA.File, err = ioutil.ReadFile("files/c-design_logotype.jpg"); err != nil {
		log.Println(fmt.Sprintf("[MEDIA_FILE_CREATE_ERROR]: %s", err.Error()))
	}

	if err := akeneoApi.MediaFile.Create(productMediaFileA); err != nil {
		log.Println(fmt.Sprintf("[MEDIA_FILE_CREATE_ERROR]: %s", err.Message))
	}
}


func downloadMediaFile() () {

	resp := getAllProductMediaFile()

	if len(resp.Data.Items) == 0 {
		log.Println("[MEDIA_FILE_GET_NOTICE]: No images")
		return
	}

	dirPath := "runtime"
	if err := akeneoApi.MediaFile.Download(resp.Data.Items[len(resp.Data.Items) - 1].Code, dirPath); err != nil {
		log.Println(fmt.Sprintf("[MODEL_DOWNLOAD_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[MEDIA_FILE_DOWNLOADED]: %s", dirPath))
	}
}

func getProductMediaFile() {

	resp := getAllProductMediaFile()

	if len(resp.Data.Items) == 0 {
		log.Println("[MEDIA_FILE_GET_NOTICE]: No images")
		return
	}

	prod, err := akeneoApi.MediaFile.Get(resp.Data.Items[0].Code)
	if err != nil {
		log.Println(fmt.Sprintf("[MEDIA_FILE_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[MEDIA_FILE_GET]: %s", prod.Code))
	}
}

func getAllProductMediaFile() *akeneo.ProductMediaFileResponse {
	opts := akeneo.RequestOpts{}
	items, err := akeneoApi.MediaFile.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[MEDIA_FILE_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, item := range items.Data.Items {
			log.Println(fmt.Sprintf("[MEDIA_FILE_GET_ALL]: %s", item.Code))
		}
	}

	return items
}
