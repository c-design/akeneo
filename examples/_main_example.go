package main

import (
	akeneo "../api"
	"time"
)

var akeneoApi *akeneo.Api

func _main() {

	config := &akeneo.ClientConfig{
		BaseUrl:        "http://akeneo-pim-host.com",
		UserAgent:      "ImportToAkeneoService",
		Username:       "admin", // http://akeneo-pim-host.com/#/user/
		Password:       "admin",
		ClientId:       "client_id", // http://akeneo-pim-host.com/#/client/
		SecretKey:      "secret_key",
		RequestTimeout: time.Second * 5,
	}

	client := akeneo.NewClient(config)
	akeneoApi = akeneo.NewAkeneoApi(client)

	runCategoriesMethods()
	runAttributeGroupMethods()
	runAttributesMethods()
	runAttributeOptionMethods()
	runFamiliesMethods()
	runFamilyVariantsMethods()
	runProductModelsMethods()
	runProductsMethods()
	runProductMediaFilesMethods()
	runAssociationTypeMethods()
	runChannelsMethods()
	runLocalesMethods()
	runCurrenciesMethods()
	runMeasureFamilyMethods()
}