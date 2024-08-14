package client

import (
	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/interface/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// this main grpc client
type Grpc struct {
	Product     delivery.ProductGrpc
	productConn *grpc.ClientConn
}

func NewGrpc(pgd delivery.ProductGrpc, productConn *grpc.ClientConn) *Grpc {
	return &Grpc{
		Product:     pgd,
		productConn: productConn,
	}
}

func (g *Grpc) Close() {
	if err := g.productConn.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.Grpc/Close", "section": "productConn.Close"}).Error(err)
	}
}
