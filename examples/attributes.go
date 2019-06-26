package main

import (
	"fmt"
	"log"
)
import "../api"

var attrSize = &akeneo.Attribute{
	Code:             "drop_size",
	Type_:            akeneo.AkeneoTypeSimpleSelect,
	Labels:           map[string]string{"ru_RU": "Размер", "en_US": "Size"},
	Group:            "base",
	SortOrder:        0,
	Localizable:      false,
	Scopable:         false,
	AvailableLocales: nil,
}

var attrColor = &akeneo.Attribute{
	Code:             "drop_color",
	Type_:            akeneo.AkeneoTypeSimpleSelect,
	Labels:           map[string]string{"ru_RU": "Цвет", "en_US": "Color"},
	Group:            "base",
	SortOrder:        0,
	Localizable:      false,
	Scopable:         false,
	AvailableLocales: nil,
}

var attrImage = &akeneo.Attribute{
	Code:             "drop_image",
	Type_:            akeneo.AkeneoTypeFile,
	Labels:           map[string]string{"ru_RU": "Изображение", "en_US": "Image"},
	Group:            "base",
	SortOrder:        0,
	Localizable:      false,
	Scopable:         false,
	AvailableLocales: nil,
}

func runAttributesMethods() {
	createAttribute()
	upsertAttribute()
	batchUpsertAttribute()
	getAttribute()
	getAllAttributes()
}

func createAttribute() {
	if err := akeneoApi.Attribute.Create(attrColor); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertAttribute() {
	attrColor.Labels["en_US"] = "Change color attribute title"
	if err := akeneoApi.Attribute.Upsert(attrColor); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertAttribute() {
	attrs := append([]*akeneo.Attribute{}, attrColor, attrImage, attrSize)
	resp, err := akeneoApi.Attribute.BatchUpsert(attrs)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[ATTRIBUTE_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[ATTRIBUTE_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getAttribute() {
	attr, err := akeneoApi.Attribute.Get(attrSize.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GET]: %s", attr.Code))
	}
}

func getAllAttributes() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Attribute.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[ATTRIBUTE_GET_ALL]: %s", prod.Code))
		}
	}
}
