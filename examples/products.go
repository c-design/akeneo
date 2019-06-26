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

var productA = &akeneo.Product{
	Identifier: "product_a",
	Enabled:    false,
	FamilyCode: famA.Code,
	Categories: []string{categoryA.Code},
	Values: map[string][]*akeneo.ProductAttributeValue{
		attrColor.Code: {
			{
				Data: "white",
			},
		},
		attrSize.Code: {
			{
				Data: "XXL",
			},
		},
	},
}

var productB = &akeneo.Product{
	Identifier: "product_b",
	Enabled:    true,
	FamilyCode: famA.Code,
	Categories: []string{categoryB.Code},
	Values: map[string][]*akeneo.ProductAttributeValue{
		attrColor.Code: {
			{
				Data: "black",
			},
		},
		attrSize.Code: {
			{
				Data: "X",
			},
		},
	},
}

var productC = &akeneo.Product{
	Identifier: "product_a_variant",
	Enabled:    true,
	Parent:     productModelA.Code,
	FamilyCode: famA.Code,
	Categories: []string{categoryA.Code},
	Values: map[string][]*akeneo.ProductAttributeValue{
		attrColor.Code: {
			{
				Data: "white",
			},
		},
		attrSize.Code: {
			{
				Data: "XLLL",
			},
		},
	},
}

var productD = &akeneo.Product{
	Identifier: "product_b_variant",
	Enabled:    true,
	Parent:     productModelB.Code,
	FamilyCode: famA.Code,
	Categories: []string{categoryA.Code},
	Values: map[string][]*akeneo.ProductAttributeValue{
		attrColor.Code: {
			{
				Data: "black",
			},
		},
		attrSize.Code: {
			{
				Data: "X",
			},
		},
	},
}

func runProductsMethods() {
	createProduct()
	upsertProduct()
	batchUpsertProduct()
	getProduct()
	deleteProduct()
	getAllProducts()
}

func createProduct() {

	if err := akeneoApi.Product.Create(productA); err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_CREATE_ERROR]: %s", err.Message))
	}

}

func upsertProduct() () {
	if err := akeneoApi.Product.Upsert(productA); err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_UPSERT_ERROR]: %s", err.Message))
	}

}

func batchUpsertProduct() () {
	productB.Enabled = false

	var list = append([]*akeneo.Product{}, productA, productB, productC, productD)
	resp, err := akeneoApi.Product.BatchUpsert(list)

	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[PRODUCT_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[PRODUCT_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getProduct() {
	prod, err := akeneoApi.Product.Get("product_a")
	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[PRODUCT_GET]: %s", prod.Identifier))
	}
}

func deleteProduct() {
	err := akeneoApi.Product.Delete("product_a")
	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_DELETE_ERROR]: %s", err.Message))
	}
}

func getAllProducts() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Product.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[PRODUCT_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[PRODUCT_GET_ALL]: %s", prod.Identifier))
		}
	}
}
