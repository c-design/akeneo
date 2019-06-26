package akeneo

type ApiService struct {
	client *Client
}

type Api struct {
	Category           *CategoriesApi
	Family             *FamilyApi
	FamilyVariant      *FamilyVariantApi
	Attribute          *AttributeApi
	AttributeOption    *AttributeOptionApi
	AttributeGroup     *AttributeGroupApi
	AssociationTypeApi *AssociationTypeApi
	Product            *ProductApi
	ProductModel       *ProductModelApi
	MediaFile          *MediaFileApi
	Channel            *ChannelApi
	Locale             *LocaleApi
	Currency           *CurrencyApi
	MeasureFamily      *MeasureFamilyApi
}

func NewAkeneoApi(client *Client) *Api {
	akeneoApi := &Api{}
	service := &ApiService{client: client}

	akeneoApi.Category = (*CategoriesApi)(service)
	akeneoApi.Family = (*FamilyApi)(service)
	akeneoApi.FamilyVariant = (*FamilyVariantApi)(service)
	akeneoApi.Attribute = (*AttributeApi)(service)
	akeneoApi.AttributeOption = (*AttributeOptionApi)(service)
	akeneoApi.AttributeGroup = (*AttributeGroupApi)(service)
	akeneoApi.AssociationTypeApi = (*AssociationTypeApi)(service)
	akeneoApi.Product = (*ProductApi)(service)
	akeneoApi.ProductModel = (*ProductModelApi)(service)
	akeneoApi.MediaFile = (*MediaFileApi)(service)
	akeneoApi.Channel = (*ChannelApi)(service)
	akeneoApi.Locale = (*LocaleApi)(service)
	akeneoApi.Currency = (*CurrencyApi)(service)
	akeneoApi.MeasureFamily = (*MeasureFamilyApi)(service)

	return akeneoApi
}

type RequestOpts map[string]interface{}

type Response struct {
	Links       ResponseLinks `json:"_links"`
	CurrentPage int           `json:"current_page"`
}

type ResponseLinks struct {
	Self     ResponseLink `json:"self"`
	First    ResponseLink `json:"first"`
	Previous ResponseLink `json:"previous"`
	Next     ResponseLink `json:"next"`
}

type ResponseLink struct {
	Href string `json:"href"`
}

type ResponseBody struct {
	Line       int32  `json:"line"`
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
	StatusCode int32  `json:"status_code"`
	Message    string `json:"message"`
}

type ResponseBodyLinks struct {
	Documentation ResponseLink `json:"documentation"`
}

type ApiError struct {
	Code    int
	Status  string
	Message string
}
