package main

import (
	"fmt"
	"log"
)
import "../api"

var currencyUsdCode = "USD"

func runCurrenciesMethods() {
	getCurrency()
	getAllCurrencies()
}

func getCurrency() {
	channel, err := akeneoApi.Currency.Get(currencyUsdCode)
	if err != nil {
		log.Println(fmt.Sprintf("[CURRENCY_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[CURRENCY_GET]: %s", channel.Code))
	}
}

func getAllCurrencies() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Currency.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[CURRENCY_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, item := range resp.Data.Items {
			log.Println(fmt.Sprintf("[CURRENCY_GET_ALL]: %s", item.Code))
		}
	}
}
