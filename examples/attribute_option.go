package main

import (
	"fmt"
	"log"
)
import "../api"

var optA = &akeneo.AttributeOption{
	Code:      "black",
	SortOrder: 0,
	Labels:    map[string]string{"ru_RU": "Черный", "en_US": "Black"},
}

var optB = &akeneo.AttributeOption{
	Code:      "white",
	SortOrder: 0,
	Labels:    map[string]string{"ru_RU": "Белый", "en_US": "White"},
}

var optC = &akeneo.AttributeOption{
	Code:      "X",
	SortOrder: 0,
	Labels:    map[string]string{"ru_RU": "Размер X", "en_US": "Размер XL"},
}

var optD = &akeneo.AttributeOption{
	Code:      "XXL",
	SortOrder: 0,
	Labels:    map[string]string{"ru_RU": "Размер XXL", "en_US": "Size XXL"},
}

func runAttributeOptionMethods() {
	createAttributeOption()
	upsertAttributeOption()
	batchUpsertColorAttributeOption()
	batchUpsertSizeAttributeOption()
	getAttributeOption()
	getAllAttributeOptions()
}

func createAttributeOption() {
	if err := akeneoApi.AttributeOption.Create(attrColor.Code, optA); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertAttributeOption() {
	optA.Labels["en_US"] = "Change color attribute title"
	if err := akeneoApi.AttributeOption.Upsert(attrColor.Code, optA); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertColorAttributeOption() {
	options := append([]*akeneo.AttributeOption{}, optA, optB)
	resp, err := akeneoApi.AttributeOption.BatchUpsert(attrColor.Code, options)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_COLOR_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_COLOR_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_COLOR_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func batchUpsertSizeAttributeOption() {
	options := append([]*akeneo.AttributeOption{}, optC, optD)
	resp, err := akeneoApi.AttributeOption.BatchUpsert(attrSize.Code, options)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_SIZE_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_SIZE_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_SIZE_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getAttributeOption() {
	attr, err := akeneoApi.AttributeOption.Get(attrColor.Code, optA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_GET]: %s", attr.Code))
	}
}

func getAllAttributeOptions() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.AttributeOption.GetAll(attrColor.Code, opts)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[ATTRIBUTE_OPTION_GET_ALL]: %s", prod.Code))
		}
	}
}
