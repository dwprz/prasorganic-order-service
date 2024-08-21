package queue

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return &ClientMock{
		Mock: mock.Mock{},
	}
}

func (s *ClientMock) Create(typename string, queue string, payload []byte, delay time.Duration) {}

func (s *ClientMock) Close() {}
