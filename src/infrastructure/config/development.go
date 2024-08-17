package config

import (
	"os"
	"strconv"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setUpForDevelopment() *Config {
	err := os.Chdir(os.Getenv("PRASORGANIC_ORDER_SERVICE_WORKSPACE"))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "os.Chdir"}).Fatal(err)
	}

	viper := viper.New()
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "viper.ReadInConfig"}).Fatal(err)
	}

	storeConf := new(store)
	storeConf.Name = viper.GetString("STORE_NAME")
	storeConf.PhoneNumber = viper.GetString("STORE_PHONE_NUMBER")
	storeConf.Address = viper.GetString("STORE_ADDRESS")

	areaId, err := strconv.Atoi(viper.GetString("STORE_AREA_ID"))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForDevelopment", "section": "strconv.Atoi"}).Fatal(err)
	}

	storeConf.AreaId = areaId
	storeConf.Latitude = viper.GetString("STORE_LATITUDE")
	storeConf.Longitude = viper.GetString("STORE_LONGITUDE")

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = viper.GetString("CURRENT_APP_RESTFUL_ADDRESS")
	currentAppConf.GrpcPort = viper.GetString("CURRENT_APP_GRPC_PORT")

	frontEndConf := new(frontEnd)
	frontEndConf.BaseUrl = viper.GetString("FRONT_END_BASE_URL")

	kafkaConf := new(kafka)
	kafkaConf.Addr1 = viper.GetString("KAFKA_ADDRESS_1")
	kafkaConf.Addr2 = viper.GetString("KAFKA_ADDRESS_2")
	kafkaConf.Addr3 = viper.GetString("KAFKA_ADDRESS_3")

	postgresConf := new(postgres)
	postgresConf.Url = viper.GetString("POSTGRES_URL")
	postgresConf.Dsn = viper.GetString("POSTGRES_DSN")
	postgresConf.User = viper.GetString("POSTGRES_USER")
	postgresConf.Password = viper.GetString("POSTGRES_PASSWORD")

	redisConf := new(redis)
	redisConf.AddrNode1 = viper.GetString("REDIS_ADDR_NODE_1")
	redisConf.AddrNode2 = viper.GetString("REDIS_ADDR_NODE_2")
	redisConf.AddrNode3 = viper.GetString("REDIS_ADDR_NODE_3")
	redisConf.AddrNode4 = viper.GetString("REDIS_ADDR_NODE_4")
	redisConf.AddrNode5 = viper.GetString("REDIS_ADDR_NODE_5")
	redisConf.AddrNode6 = viper.GetString("REDIS_ADDR_NODE_6")
	redisConf.Password = viper.GetString("REDIS_PASSWORD")

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = viper.GetString("API_GATEWAY_BASE_URL")
	apiGatewayConf.BasicAuth = viper.GetString("API_GATEWAY_BASIC_AUTH")
	apiGatewayConf.BasicAuthUsername = viper.GetString("API_GATEWAY_BASIC_AUTH_USERNAME")
	apiGatewayConf.BasicAuthPassword = viper.GetString("API_GATEWAY_BASIC_AUTH_PASSWORD")

	jwtConf := new(jwt)
	jwtConf.PrivateKey = loadRSAPrivateKey(viper.GetString("JWT_PRIVATE_KEY"))
	jwtConf.PublicKey = loadRSAPublicKey(viper.GetString("JWT_PUBLIC_KEY"))

	midtransConf := new(midtrans)
	midtransConf.BaseUrl = viper.GetString("MIDTRANS_BASE_URL")
	midtransConf.ServerKey = viper.GetString("MIDTRANS_SERVER_KEY")

	shipperConf := new(shipper)
	shipperConf.BaseUrl = viper.GetString("SHIPPER_BASE_URL")
	shipperConf.ApiKey = viper.GetString("SHIPPER_API_KEY")

	shippingConf := new(shipping)
	shippingConf.Coverage = viper.GetString("SHIPPING_COVERAGE")
	shippingConf.PaymentType = viper.GetString("SHIPPING_PAYMENT_TYPE")

	return &Config{
		Store:      storeConf,
		CurrentApp: currentAppConf,
		FrontEnd:   frontEndConf,
		Kafka:      kafkaConf,
		Postgres:   postgresConf,
		Redis:      redisConf,
		ApiGateway: apiGatewayConf,
		Jwt:        jwtConf,
		Midtrans:   midtransConf,
		Shipper:    shipperConf,
		Shipping:   shippingConf,
	}
}
