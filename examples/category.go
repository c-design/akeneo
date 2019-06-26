package main

import (
	"fmt"
	"log"
)
import "../api"

var categoryRoot = &akeneo.Category{
	Code:   "root_category",
	Labels: map[string]string{"en_US": "Root Category"},
}

var categoryA = &akeneo.Category{
	Code:   "category_a",
	Parent: &categoryRoot.Code,
	Labels: map[string]string{"en_US": "Category A"},
}

var categoryB = &akeneo.Category{
	Code:   "category_b",
	Parent: &categoryRoot.Code,
	Labels: map[string]string{"en_US": "Category B"},
}

func runCategoriesMethods() {
	createCategory()
	upsertCategory()
	batchUpsertCategory()
	getCategory()
	getAllCategories()
}

func createCategory() {
	if err := akeneoApi.Category.Create(categoryRoot); err != nil {
		log.Println(fmt.Sprintf("[CATEGORY_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertCategory() {
	categoryRoot.Labels["en_US"] = "Change root category title"
	if err := akeneoApi.Category.Upsert(categoryRoot); err != nil {
		log.Println(fmt.Sprintf("[CATEGORY_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertCategory() {
	categories := append([]*akeneo.Category{}, categoryA, categoryB)
	resp, err := akeneoApi.Category.BatchUpsert(categories)

	if err != nil {
		log.Println(fmt.Sprintf("[CATEGORY_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[CATEGORY_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[CATEGORY_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getCategory() {
	category, err := akeneoApi.Category.Get(categoryRoot.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[CATEGORY_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[CATEGORY_GET]: %s", category.Code))
	}
}

func getAllCategories() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Category.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[CATEGORY_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, prod := range resp.Data.Items {
			log.Println(fmt.Sprintf("[CATEGORY_GET_ALL]: %s", prod.Code))
		}
	}
}
