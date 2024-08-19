package cbreaker

import (
	"net/http"
	"time"

	"github.com/dwprz/prasorganic-order-service/src/common/errors"
	"github.com/dwprz/prasorganic-order-service/src/common/log"
	"github.com/sony/gobreaker/v2"
)

var Shipper *gobreaker.CircuitBreaker[any]

func init() {
	settings := gobreaker.Settings{
		Name:        "shipper-restful",
		MaxRequests: 3,
		Interval:    1 * time.Minute,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {

			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.8
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}

			if errRes, ok := err.(*errors.Response); ok {
				successCode := []int{
					http.StatusOK,
					http.StatusCreated,
					http.StatusAccepted,
					http.StatusNoContent,
					http.StatusNotFound,
				}

				for _, code := range successCode {
					if errRes.HttpCode == code {
						return true
					}
				}
			}

			return false
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Logger.Infof("circuit breaker %v from status %v to %v", name, from, to)
		},
	}

	Shipper = gobreaker.NewCircuitBreaker[any](settings)
}
