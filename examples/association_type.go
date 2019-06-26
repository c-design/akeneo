package main

import (
	"fmt"
	"log"
)
import "../api"

var typeA = &akeneo.AssociationType{
	Code:             "drop_size",
	Labels:           map[string]string{"ru_RU": "Ассоциация А", "en_US": "Association A"},
}


var typeB = &akeneo.AssociationType{
	Code:             "drop_size",
	Labels:           map[string]string{"ru_RU": "Ассоциация А", "en_US": "Association A"},
}

func runAssociationTypeMethods() {
	createAssociationType()
	upsertAssociationType()
	batchUpsertAssociationType()
	getAssociationType()
	getAllAssociationTypes()
}

func createAssociationType() {
	if err := akeneoApi.AssociationTypeApi.Create(typeA); err != nil {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertAssociationType() {
	typeA.Labels["en_US"] = "Change association A label"
	if err := akeneoApi.AssociationTypeApi.Upsert(typeA); err != nil {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertAssociationType() {
	attrs := append([]*akeneo.AssociationType{}, typeA, typeB)
	resp, err := akeneoApi.AssociationTypeApi.BatchUpsert(attrs)

	if err != nil {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getAssociationType() {
	attr, err := akeneoApi.AssociationTypeApi.Get(typeA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_GET]: %s", attr.Code))
	}
}

func getAllAssociationTypes() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.AssociationTypeApi.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[ASSOCIATION_TYPE_GET_ALL]: %s", prod.Code))
		}
	}
}
