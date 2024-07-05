package main

import (
	"context"

	"github.com/64bitAryan/go-microservice/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

// HTTP server works with JSON only
// for GRPC we are creating a new GRPC server
func NewAggregatorGRPCService(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

/*
	Transport Layer
		JSON -> types.Distance
		GRPC -> types.Aggregator -> it converts to -> types.Distance
		WEBPack -> types.WEBPack -> it converts to -> types.Distance
*/

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
