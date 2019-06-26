package main

import (
	"fmt"
	"log"
)
import "../api"

func runMeasureFamilyMethods() {
	getMeasureFamily()
	getAllMeasureFamilies()
}

func getMeasureFamily() {
	channel, err := akeneoApi.MeasureFamily.Get("Length")
	if err != nil {
		log.Println(fmt.Sprintf("[MEASURE_FAMILY_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[MEASURE_FAMILY_GET]: %s", channel.Code))
	}
}

func getAllMeasureFamilies() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.MeasureFamily.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[MEASURE_FAMILY_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, item := range resp.Data.Items {
			log.Println(fmt.Sprintf("[MEASURE_FAMILY_GET_ALL]: %s", item.Code))
		}
	}
}
