package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/sirupsen/logrus"
)

type store struct {
	Name        string
	PhoneNumber string
	Address     string
	AreaId      int
	Latitude    string
	Longitude   string
}

type currentApp struct {
	RestfulAddress string
	GrpcPort       string
}

type frontEnd struct {
	BaseUrl string
}

type apiGateway struct {
	BaseUrl           string
	BasicAuth         string
	BasicAuthUsername string
	BasicAuthPassword string
}

type kafka struct {
	Addr1 string
	Addr2 string
	Addr3 string
}

type postgres struct {
	Url      string
	Dsn      string
	User     string
	Password string
}

type redis struct {
	AddrNode1 string
	AddrNode2 string
	AddrNode3 string
	AddrNode4 string
	AddrNode5 string
	AddrNode6 string
	Password  string
}

type jwt struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type midtrans struct {
	BaseUrl   string
	ServerKey string
}

type shipper struct {
	BaseUrl string
	ApiKey  string
}

type shipping struct {
	Coverage    string
	PaymentType string
}

type Config struct {
	Store      *store
	CurrentApp *currentApp
	FrontEnd   *frontEnd
	Kafka      *kafka
	Postgres   *postgres
	Redis      *redis
	ApiGateway *apiGateway
	Jwt        *jwt
	Midtrans   *midtrans
	Shipper    *shipper
	Shipping   *shipping
}

var Conf *Config

// *config ini hanya berisi env variable
func init() {
	appStatus := os.Getenv("PRASORGANIC_APP_STATUS")

	if appStatus == "DEVELOPMENT" {

		Conf = setUpForDevelopment()
		return
	}

	Conf = setUpForNonDevelopment(appStatus)
}

func loadRSAPrivateKey(privateKey string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		log.Logger.WithFields(logrus.Fields{
			"location": "config.loadRSAPrivateKey",
			"section":  "pem.Decode",
		}).Fatal("failed to parse pem block containing the key")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"location": "config.loadRSAPrivateKey",
			"section":  "x509.ParsePKCS1PrivateKey",
		}).Fatalf("failed to parse rsa private key: %v", err)
	}

	return rsaPrivateKey
}

func loadRSAPublicKey(publicKey string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		log.Logger.WithFields(logrus.Fields{
			"location": "config.loadRSAPrivateKey",
			"section":  "pem.Decode",
		}).Fatal("failed to parse pem block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"location": "loadRSAPublicKey",
			"section":  "x509.ParsePKCS1PublicKey",
		}).Fatalf("failed to parse rsa public key: %v", err)
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Logger.WithFields(logrus.Fields{
			"location": "loadRSAPublicKey",
			"section":  "type assertion",
		}).Fatal("failed to assert type to *rsa.PublicKey")
	}

	return rsaPublicKey
}
