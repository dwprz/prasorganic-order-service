package entity

type ShippingOrder struct {
	Consignee   Consignee   `json:"consignee"`
	Consigner   Consigner   `json:"consigner"`
	Courier     Courier     `json:"courier"`
	Coverage    string      `json:"coverage"`
	Destination Destination `json:"destination"`
	ExternalId  string      `json:"external_id"`
	Origin      Origin      `json:"origin"`
	Package     Package     `json:"package"`
	PaymentType string      `json:"payment_type"`
}

type Consignee struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type Consigner struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type Courier struct {
	COD          bool `json:"cod"`
	RateId       int  `json:"rate_id"`
	UseInsurance bool `json:"use_insurance"`
}

type Destination struct {
	Address string `json:"address"`
	AreaId  int    `json:"area_id"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
}

type Origin struct {
	Address string `json:"address"`
	AreaId  int    `json:"area_id"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
}

type Package struct {
	Height      int     `json:"height"`
	Items       []Item  `json:"items"`
	Length      int     `json:"length"`
	PackageType int     `json:"package_type"`
	Price       int     `json:"price"`
	Weight      float32 `json:"weight"`
	Width       int     `json:"width"`
}

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
}
