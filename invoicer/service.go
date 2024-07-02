package main

import "github.com/64bitAryan/go-microservice/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
}
