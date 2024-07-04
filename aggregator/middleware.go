package main

import (
	"time"

	"github.com/64bitAryan/go-microservice/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
			obuID    int
		)
		if inv != nil {

			distance = inv.TotalDistance
			amount = inv.TotalAmount
			obuID = inv.OBUID
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"distance": distance,
			"amount":   amount,
		}).Info("Calculate Invoice")
	}(time.Now())
	inv, err = l.next.CalculateInvoice(obuID)
	return

}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Aggregate distance")
	}(time.Now())
	err = l.next.AggregateDistance(distance)
	return
}
