package consumer

import (
	"context"
	"strings"

	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/dwprz/prasorganic-order-service/src/core/broker/handler"
	"github.com/dwprz/prasorganic-order-service/src/infrastructure/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type MidtransKafka struct {
	reader  *kafka.Reader
	handler *handler.MidtransKafka
}

func NewMidtransKafka(mh *handler.MidtransKafka) *MidtransKafka {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Conf.Kafka.Addr1, config.Conf.Kafka.Addr2, config.Conf.Kafka.Addr3},
		GroupID: "order-consumer",
		Topic:   "midtrans",
		Logger: kafka.LoggerFunc(func(s string, i ...interface{}) {
			if strings.Contains(s, "no messages received from kafka within the allocated time for partition") {
				return
			}

			log.Logger.Infof(s, i...)
		}),
	})

	return &MidtransKafka{
		reader:  r,
		handler: mh,
	}
}

func (m *MidtransKafka) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := m.reader.ReadMessage(ctx)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{"location": "consumer.MidtransKafka/Consume", "section": "reader.ReadMessage"}).Error(err)
				continue
			}

			go func(msg kafka.Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Logger.WithFields(logrus.Fields{"location": "consumer.MidtransKafka/Consume", "section": "ProcessMessage"}).Errorf("Recovered from panic: %v", r)
					}
				}()

				m.handler.ProcessMessage(ctx, msg)
			}(msg)

		}
	}
}

func (m *MidtransKafka) Close() {
	if err := m.reader.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "consumer.MidtransKafka/Close", "section": "reader.Close"}).Error(err)
	}
}
