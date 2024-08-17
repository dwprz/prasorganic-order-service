package queue

import "time"

type Client interface {
	Create(typename string, queue string, payload []byte, delay time.Duration)
	Close()
}
