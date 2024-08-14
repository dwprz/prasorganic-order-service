package entity

type ProductOrder struct {
	OrderId         string `json:"order_id,omitempty" gorm:"column:order_id;primaryKey"`
	ProductId       int    `json:"product_id" gorm:"column:product_id"`
	ProductName     string `json:"product_name" gorm:"column:product_name"`
	Image           string `json:"image" gorm:"column:image"`
	Quantity        int    `json:"quantity" gorm:"column:quantity"`
	Price           int    `json:"price" gorm:"column:price"`
}