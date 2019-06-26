package main

import (
	"../api"
	"fmt"
	"log"
)

var (
//localeRu = "ru_RU"
//localeEn = "en_US"
)

var productModelA = &akeneo.ProductModel{
	Code:          "product_model_a",
	FamilyVariant: "color",
	Categories:    []string{"test_a"},
}

var productModelB = &akeneo.ProductModel{
	Code:          "product_model_b",
	FamilyVariant: "color",
	Categories:    []string{"test_a"},
}

func runProductModelsMethods() {
	createProductModel()
	upsertProductModel()
	batchUpsertProductModel()
	getProductModel()
	getAllProductModels()
}

func createProductModel() {
	if err := akeneoApi.ProductModel.Create(productModelA); err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertProductModel() () {
	if err := akeneoApi.ProductModel.Upsert(productModelB); err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertProductModel() () {
	var list = append([]*akeneo.ProductModel{}, productModelA, productModelB)
	resp, err := akeneoApi.ProductModel.BatchUpsert(list)

	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[PRODUCT_MODEL_BATCH_UPSERT_ERROR]: %s", respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[PRODUCT_MODEL_BATCH_UPSERT]: %s", respLine.Code))
			}
		}
	}
}

func getProductModel() {
	prod, err := akeneoApi.ProductModel.Get(productModelA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_GET]: %s", prod.Code))
	}
}

func getAllProductModels() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.ProductModel.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_MODEL_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[PRODUCT_MODEL_GET_ALL]: %s", prod.Code))
		}
	}
}
