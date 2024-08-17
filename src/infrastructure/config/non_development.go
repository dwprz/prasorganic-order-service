package config

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

func setUpForNonDevelopment(appStatus string) *Config {
	defaultConf := vault.DefaultConfig()
	defaultConf.Address = os.Getenv("PRASORGANIC_CONFIG_ADDRESS")

	client, err := vault.NewClient(defaultConf)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "vault.NewClient"}).Fatal(err)
	}

	client.SetToken(os.Getenv("PRASORGANIC_CONFIG_TOKEN"))

	mountPath := "prasorganic-secrets" + "-" + strings.ToLower(appStatus)

	storeSecrets, err := client.KVv2(mountPath).Get(context.Background(), "store")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	orderServiceSecrets, err := client.KVv2(mountPath).Get(context.Background(), "order-service")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	frontEndSecrets, err := client.KVv2(mountPath).Get(context.Background(), "front-end")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	kafkaSecrets, err := client.KVv2(mountPath).Get(context.Background(), "kafka")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	apiGatewaySecrets, err := client.KVv2(mountPath).Get(context.Background(), "api-gateway")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	jwtSecrets, err := client.KVv2(mountPath).Get(context.Background(), "jwt")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	midtransSecrets, err := client.KVv2(mountPath).Get(context.Background(), "midtrans")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	shipperSecrets, err := client.KVv2(mountPath).Get(context.Background(), "shipper")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	shippingSecrets, err := client.KVv2(mountPath).Get(context.Background(), "shipping-service")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "KVv2.Get"}).Fatal(err)
	}

	storeConf := new(store)
	storeConf.Name = storeSecrets.Data["NAME"].(string)
	storeConf.PhoneNumber = storeSecrets.Data["PHONE_NUMBER"].(string)
	storeConf.Address = storeSecrets.Data["ADDRESS"].(string)
	storeConf.AreaId = storeSecrets.Data["AREA_ID"].(int)
	storeConf.Latitude = storeSecrets.Data["LATITUDE"].(string)
	storeConf.Longitude = storeSecrets.Data["LONGITUDE"].(string)

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddress = orderServiceSecrets.Data["RESTFUL_ADDRESS"].(string)
	currentAppConf.GrpcPort = orderServiceSecrets.Data["GRPC_PORT"].(string)

	frontEndConf := new(frontEnd)
	frontEndConf.BaseUrl = frontEndSecrets.Data["BASE_URL"].(string)

	kafkaConf := new(kafka)
	kafkaConf.Addr1 = kafkaSecrets.Data["ADDRESS_1"].(string)
	kafkaConf.Addr2 = kafkaSecrets.Data["ADDRESS_2"].(string)
	kafkaConf.Addr3 = kafkaSecrets.Data["ADDRESS_3"].(string)

	postgresConf := new(postgres)
	postgresConf.Url = orderServiceSecrets.Data["POSTGRES_URL"].(string)
	postgresConf.Dsn = orderServiceSecrets.Data["POSTGRES_DSN"].(string)
	postgresConf.User = orderServiceSecrets.Data["POSTGRES_USER"].(string)
	postgresConf.Password = orderServiceSecrets.Data["POSTGRES_PASSWORD"].(string)

	redisConf := new(redis)
	redisConf.AddrNode1 = orderServiceSecrets.Data["REDIS_ADDR_NODE_1"].(string)
	redisConf.AddrNode2 = orderServiceSecrets.Data["REDIS_ADDR_NODE_2"].(string)
	redisConf.AddrNode3 = orderServiceSecrets.Data["REDIS_ADDR_NODE_3"].(string)
	redisConf.AddrNode4 = orderServiceSecrets.Data["REDIS_ADDR_NODE_4"].(string)
	redisConf.AddrNode5 = orderServiceSecrets.Data["REDIS_ADDR_NODE_5"].(string)
	redisConf.AddrNode6 = orderServiceSecrets.Data["REDIS_ADDR_NODE_6"].(string)
	redisConf.Password = orderServiceSecrets.Data["REDIS_PASSWORD"].(string)

	apiGatewayConf := new(apiGateway)
	apiGatewayConf.BaseUrl = apiGatewaySecrets.Data["BASE_URL"].(string)
	apiGatewayConf.BasicAuth = apiGatewaySecrets.Data["BASIC_AUTH"].(string)
	apiGatewayConf.BasicAuthUsername = apiGatewaySecrets.Data["BASIC_AUTH_PASSWORD"].(string)
	apiGatewayConf.BasicAuthPassword = apiGatewaySecrets.Data["BASIC_AUTH_USERNAME"].(string)

	jwtConf := new(jwt)

	jwtPrivateKey := jwtSecrets.Data["PRIVATE_KEY"].(string)
	base64Byte, err := base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPrivateKey = string(base64Byte)

	jwtPublicKey := jwtSecrets.Data["Public_KEY"].(string)
	base64Byte, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "config.setUpForNonDevelopment", "section": "base64.StdEncoding.DecodeString"}).Fatal(err)
	}
	jwtPublicKey = string(base64Byte)

	jwtConf.PrivateKey = loadRSAPrivateKey(jwtPrivateKey)
	jwtConf.PublicKey = loadRSAPublicKey(jwtPublicKey)

	midtransConf := new(midtrans)
	midtransConf.BaseUrl = midtransSecrets.Data["BASE_URL"].(string)
	midtransConf.ServerKey = midtransSecrets.Data["SERVER_KEY"].(string)

	shipperConf := new(shipper)
	shipperConf.BaseUrl = shipperSecrets.Data["BASE_URL"].(string)
	shipperConf.ApiKey = shipperSecrets.Data["API_KEY"].(string)

	shippingConf := new(shipping)
	shippingConf.Coverage = shippingSecrets.Data["COVERAGE"].(string)
	shippingConf.PaymentType = shippingSecrets.Data["PAYMENT_TYPE"].(string)

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
