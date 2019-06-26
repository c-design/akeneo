package main

import (
	"fmt"
	"log"
)
import "../api"

var famA = &akeneo.Family{
	Code:             "t_shirt",
	AttributeAsLabel: "sku",
	AttributeAsImage: attrImage.Code,
	Labels:           map[string]string{"en_US": "T-Shirt"},
	Attributes:       []string{attrColor.Code, attrSize.Code, attrImage.Code},
}

var famB = &akeneo.Family{
	Code:             "cap",
	AttributeAsLabel: "sku",
	AttributeAsImage: attrImage.Code,
	Attributes:       []string{attrColor.Code, attrSize.Code, attrImage.Code},
	Labels:           map[string]string{"en_US": "Cap"},
}

func runFamiliesMethods() {
	createFamily()
	upsertFamily()
	batchUpsertFamily()
	getFamily()
	getAllFamilies()
}

func createFamily() {
	if err := akeneoApi.Family.Create(famA); err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertFamily() {
	famA.Labels["en_US"] = "Change t-shirt family title"
	if err := akeneoApi.Family.Upsert(famA); err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertFamily() {
	items := append([]*akeneo.Family{}, famA, famB)
	resp, err := akeneoApi.Family.BatchUpsert(items)

	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_COLOR_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[FAMALIES_COLOR_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[FAMALIES_COLOR_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getFamily() {
	attr, err := akeneoApi.Family.Get(famA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[FAMALIES_GET]: %s", attr.Code))
	}
}

func getAllFamilies() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Family.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[FAMALIES_GET_ALL]: %s", prod.Code))
		}
	}
}
