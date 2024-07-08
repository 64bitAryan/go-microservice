package client

import (
	"context"

	"github.com/64bitAryan/go-microservice/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
