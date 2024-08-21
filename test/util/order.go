package util

import (
	"time"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderTest struct {
	db *gorm.DB
}

func NewOrderTest(db *gorm.DB) *OrderTest {
	return &OrderTest{
		db: db,
	}
}

func (o *OrderTest) Create() *entity.OrderWithProducts {
	orderId, _ := gonanoid.New()
	//userId, _ := gonanoid.New()

	order := &entity.Order{
		OrderId:         orderId,
		GrossAmount:     500000,
		Status:          entity.PENDING_PAYMENT,
		ShippingId:      "ship-456",
		Courier:         "JNE",
		RateId:          1,
		RateName:        "Express",
		RateType:        "Overnight",
		COD:             false,
		UseInsurance:    true,
		PackageType:     2,
		PaymentMethod:   "Credit Card",
		SnapToken:       "snap-token-789",
		SnapRedirectURL: "https://payment.gateway/redirect",
		UserId:          "user_1",
		Email:           "user@example.com",
		Buyer:           "John Doe",
		Height:          10,
		Length:          20,
		Width:           15,
		Weight:          2.5,
		AddressOwner:    "John Doe",
		Street:          "123 Main St",
		AreaId:          1234,
		Area:            "Central",
		Lat:             "-6.200000",
		Lng:             "106.816666",
		Suburb:          "Jakarta Pusat",
		City:            "Jakarta",
		Province:        "DKI Jakarta",
		WhatsApp:        "081234567890",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := o.db.Create(order).Error; err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "util.NewOrderTest/Create", "section": "db.Create"}).Error(err)
	}

	products := []*entity.ProductOrder{
		{
			OrderId:     orderId,
			ProductId:   1,
			ProductName: "Product A",
			Image:       "https://example.com/image-a.jpg",
			Quantity:    2,
			Price:       150000,
		},
		{
			OrderId:     orderId,
			ProductId:   2,
			ProductName: "Product B",
			Image:       "https://example.com/image-b.jpg",
			Quantity:    1,
			Price:       200000,
		},
	}

	if err := o.db.Create(products).Error; err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "util.NewOrderTest/Create", "section": "db.Create"}).Error(err)
	}

	return &entity.OrderWithProducts{
		Order:    order,
		Products: products,
	}
}

func (o *OrderTest) Delete() {
	if err := o.db.Exec("DELETE FROM product_orders").Error; err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "util.NewOrderTest/Delete", "section": "db.Exec"}).Error(err.Error())
	}

	if err := o.db.Exec("DELETE FROM orders").Error; err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "util.NewOrderTest/Delete", "section": "db.Exec"}).Error(err.Error())
	}
}
