package main

import (
	"fmt"
	"log"
)
import "../api"

var groupA =  &akeneo.AttributeGroup{
	Code:       "drop_base",
	SortOrder:  0,
	Attributes: nil,
	Labels:     map[string]string{"ru_RU": "Основная", "en_US": "Base"},
}

var groupB =  &akeneo.AttributeGroup{
	Code:       "drop_specifications",
	SortOrder:  0,
	Attributes: nil,
	Labels:     map[string]string{"ru_RU": "Характеристики", "en_US": "Specifications"},
}

func runAttributeGroupMethods() {
	createAttributeGroup()
	upsertAttributeGroup()
	batchUpsertAttributeGroup()
	getAttributeGroup()
	getAllAttributeGroups()
}


func createAttributeGroup() {
	if err := akeneoApi.AttributeGroup.Create(groupA); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertAttributeGroup() {
	groupA.Labels["en_US"] = "Change attribute group title"
	if err := akeneoApi.AttributeGroup.Upsert(groupA); err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertAttributeGroup() {
	groups := append([]*akeneo.AttributeGroup{}, groupA, groupB)
	resp, err := akeneoApi.AttributeGroup.BatchUpsert(groups)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getAttributeGroup() {
	group, err := akeneoApi.AttributeGroup.Get(groupA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_GET]: %s", group.Code))
	}
}

func getAllAttributeGroups() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.AttributeGroup.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[ATTRIBUTE_GROUP_GET_ALL]: %s", prod.Code))
		}
	}
}
