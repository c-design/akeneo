package main

import (
	"fmt"
	"log"
)
import "../api"

var variantA = &akeneo.FamilyVariant{
	Code: "drop_var_color",
	AttributeSets: []*akeneo.FamilyVariantAttributeSets{
		{
			Level:      1,
			Axes:       []string{attrColor.Code},
			Attributes: []string{attrColor.Code, "sku"},
		},
	},
	Labels: map[string]string{"ru_RU": "Цветовой вариант", "en_US": "Color variant"},
}

var variantB = &akeneo.FamilyVariant{
	Code: "drop_var_size",
	AttributeSets: []*akeneo.FamilyVariantAttributeSets{
		{
			Level:      1,
			Axes:       []string{attrSize.Code},
			Attributes: []string{attrSize.Code, "sku"},
		},
	},
	Labels: map[string]string{"ru_RU": "Вариант размеров", "en_US": "Size variant"},
}

func runFamilyVariantsMethods() {
	createFamilyVariant()
	upsertFamilyVariant()
	batchUpsertFamilyVariant()
	getFamilyVariant()
	getAllFamilyVariants()
}

func createFamilyVariant() {
	if err := akeneoApi.FamilyVariant.Create(famA.Code, variantA); err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertFamilyVariant() {
	variantA.Labels["en_US"] = "Change color variant title"
	if err := akeneoApi.FamilyVariant.Upsert(famA.Code, variantA); err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertFamilyVariant() {
	items := append([]*akeneo.FamilyVariant{}, variantA, variantB)
	resp, err := akeneoApi.FamilyVariant.BatchUpsert(famA.Code, items)

	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_COLOR_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[FAMALIES_VARIANT_COLOR_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[FAMALIES_VARIANT_COLOR_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getFamilyVariant() {
	attr, err := akeneoApi.FamilyVariant.Get(famA.Code, variantA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_GET]: %s", attr.Code))
	}
}

func getAllFamilyVariants() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.FamilyVariant.GetAll(famA.Code, opts)

	if err != nil {
		log.Println(fmt.Sprintf("[FAMALIES_VARIANT_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[FAMALIES_VARIANT_GET_ALL]: %s", prod.Code))
		}
	}
}
