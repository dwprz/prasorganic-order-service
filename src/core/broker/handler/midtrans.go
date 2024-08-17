package handler

import (
	"context"
	"encoding/json"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/interface/service"
	"github.com/dwprz/prasorganic-order-service/src/model/entity"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type MidtransKafka struct {
	txService service.Transaction
}

func NewMidtransKafka(ts service.Transaction) *MidtransKafka {
	return &MidtransKafka{
		txService: ts,
	}
}

func (m *MidtransKafka) ProcessMessage(ctx context.Context, msg kafka.Message) {
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {

		tx := new(entity.Transaction)
		if err := json.Unmarshal(msg.Value, &tx); err != nil {
			log.Logger.WithFields(logrus.Fields{"location": "handler.MidtransKafka", "section": "json.Unmarshal"})
			continue
		}

		if err := m.txService.HandleNotif(ctx, tx); err != nil {
			log.Logger.WithFields(logrus.Fields{"location": "handler.MidtransKafka", "section": "txService.HandleNotif"})
			continue
		}

		break
	}
}
