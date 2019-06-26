package main

import (
	"fmt"
	"log"
)
import "../api"

var channelA = &akeneo.Channel{
	Code:            "channel_a",
	Locales:         []string{"en_US", "ru_RU"},
	Currencies:      []string{"EUR", "USD"},
	CategoryTree:    categoryRoot.Code,
	Labels:          map[string]string{"en_US": "Channel A"},
}

var channelB = &akeneo.Channel{
	Code:            "channel_b",
	Locales:         []string{"en_US", "ru_RU"},
	Currencies:      []string{"EUR", "USD"},
	CategoryTree:    categoryRoot.Code,
	Labels:          map[string]string{"en_US": "Channel B"},
}

func runChannelsMethods() {
	createChannel()
	upsertChannel()
	batchUpsertChannel()
	getChannel()
	getAllChannels()
}

func createChannel() {
	if err := akeneoApi.Channel.Create(channelA); err != nil {
		log.Println(fmt.Sprintf("[CHANNEL_CREATE_ERROR]: %s", err.Message))
	}
}

func upsertChannel() {
	channelA.Labels["en_US"] = "Change channel A title"
	if err := akeneoApi.Channel.Upsert(channelA); err != nil {
		log.Println(fmt.Sprintf("[CHANNEL_UPSERT_ERROR]: %s", err.Message))
	}
}

func batchUpsertChannel() {
	items := append([]*akeneo.Channel{}, channelA, channelB)
	resp, err := akeneoApi.Channel.BatchUpsert(items)

	if err != nil {
		log.Println(fmt.Sprintf("[CHANNEL_BATCH_UPSERT_ERROR]: %s", err.Message))
	} else {
		for _, respLine := range resp {
			if respLine.StatusCode >= 300 {
				log.Println(fmt.Sprintf("[CHANNEL_BATCH_UPSERT_ERROR]: %s => %s", respLine.Identifier, respLine.Message))
			} else {
				log.Println(fmt.Sprintf("[CHANNEL_BATCH_UPSERT]: %d", respLine.StatusCode))
			}
		}
	}
}

func getChannel() {
	channel, err := akeneoApi.Channel.Get(channelA.Code)
	if err != nil {
		log.Println(fmt.Sprintf("[CHANNEL_GET_ERROR]: %s", err.Message))
	} else {
		log.Println(fmt.Sprintf("[CHANNEL_GET]: %s", channel.Code))
	}
}

func getAllChannels() {
	opts := akeneo.RequestOpts{}
	resp, err := akeneoApi.Channel.GetAll(opts)

	if err != nil {
		log.Println(fmt.Sprintf("[CHANNEL_GET_ALL_ERROR]: %s", err.Message))
	} else {
		for _, item := range resp.Data.Items {
			log.Println(fmt.Sprintf("[CHANNEL_GET_ALL]: %s", item.Code))
		}
	}
}
