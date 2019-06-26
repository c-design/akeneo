package main

import (
	"fmt"
	"log"
)
import "../api"

var localeEnCode = "en_US"

func runLocalesMethods() {
	getLocale()
	getAllLocales()
}

func getLocale() {
	channel, err := akeneoApi.Locale.Get(localeEnCode)
	if err != nil {
		log.Println(fmt.Sprintf("[LOCALE_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[LOCALE_GET]: %s", channel.Code))
	}
}

func getAllLocales() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Locale.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[LOCALE_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, item := range resp.Data.Items {
			log.Println(fmt.Sprintf("[LOCALE_GET_ALL]: %s", item.Code))
		}
	}
}
