package entity

import "time"

type OrderStatus string

const (
	PENDING_PAYMENT   OrderStatus = "PENDING_PAYMENT"
	PAID              OrderStatus = "PAID"
	IN_PROGRESS       OrderStatus = "IN_PROGRESS"
	COMPLETED         OrderStatus = "COMPLETED"
	CANCELLED         OrderStatus = "CANCELLED"
	FAILED            OrderStatus = "FAILED"
	REFUND_PROCESSING OrderStatus = "REFUND_PROCESSING"
	REFUND_COMPLETED  OrderStatus = "REFUND_COMPLETED"
	RETURN_PROCESSING OrderStatus = "RETURN_PROCESSING"
	RETURN_COMPLETED  OrderStatus = "RETURN_COMPLETED"
	LOST_OR_DAMAGED   OrderStatus = "LOST_OR_DAMAGED"
)

type Order struct {
	OrderId         string      `json:"order_id,omitempty" gorm:"column:order_id;primaryKey"`
	GrossAmount     int         `json:"gross_amount" gorm:"column:gross_amount"`
	Status          OrderStatus `json:"status" gorm:"column:status;default:PENDING_PAYMENT"`
	ShippingId      string      `json:"shipping_id,omitempty" gorm:"column:shipping_id"`
	Courier         string      `json:"courier" gorm:"column:courier"`
	RateId          int         `json:"rate_id" gorm:"column:rate_id"`
	RateName        string      `json:"rate_name" gorm:"column:rate_name"`
	RateType        string      `json:"rate_type" gorm:"column:rate_type"`
	COD             bool        `json:"cod" gorm:"column:cod"`
	UseInsurance    bool        `json:"use_insurance" gorm:"column:use_insurance"`
	PackageType     int         `json:"package_type" gorm:"column:package_type"`
	PaymentMethod   string      `json:"payment_method,omitempty" gorm:"column:payment_method"`
	SnapToken       string      `json:"snap_token,omitempty" gorm:"column:snap_token"`
	SnapRedirectURL string      `json:"snap_redirect_url,omitempty" gorm:"column:snap_redirect_url"`
	UserId          string      `json:"user_id" gorm:"column:user_id"`
	Email           string      `json:"email" gorm:"column:email"`
	Buyer           string      `json:"buyer" gorm:"column:buyer"`
	Height          int         `json:"height" gorm:"column:height"`
	Length          int         `json:"length" gorm:"column:length"`
	Width           int         `json:"width" gorm:"column:width"`
	Weight          float32     `json:"weight" gorm:"column:weight"`
	AddressOwner    string      `json:"address_owner" gorm:"column:address_owner"`
	Street          string      `json:"street" gorm:"column:street"`
	AreaId          int         `json:"area_id" gorm:"column:area_id"`
	Area            string      `json:"area" gorm:"column:area"`
	Lat             string      `json:"lat" gorm:"column:lat"`
	Lng             string      `json:"lng" gorm:"column:lng"`
	Suburb          string      `json:"suburb" gorm:"column:suburb"`
	City            string      `json:"city" gorm:"column:city"`
	Province        string      `json:"province" gorm:"column:province"`
	WhatsApp        string      `json:"whatsapp" gorm:"column:whatsapp"`
	CreatedAt       time.Time   `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt       time.Time   `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type OrderWithProducts struct {
	Order    *Order          `json:"order"`
	Products []*ProductOrder `json:"products"`
}
