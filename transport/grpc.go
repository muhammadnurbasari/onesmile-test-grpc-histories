package transport

import (
	"context"

	logkit "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"google.golang.org/protobuf/types/known/emptypb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/endpointHistory"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/models"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
)

type grpcServer struct {
	createHistory grpctransport.Handler
	histories     grpctransport.Handler
}

func NewGrpcServer(endpoints endpointHistory.SetHistory, logger logkit.Logger) generate.TransactionsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		createHistory: grpctransport.NewServer(
			endpoints.CreateHistoryEndpoint,
			decodeGRPCCreateHistoryRequest,
			encodeGRPCCreateHistoryResponse,
			options...,
		),
		histories: grpctransport.NewServer(
			endpoints.HistoryEndpoint,
			decodeGRPCHistoryRequest,
			encodeGRPCHistoryResponse,
			options...,
		),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *generate.Transaction) (*emptypb.Empty, error) {
	_, _, err := s.createHistory.ServeGRPC(ctx, req)

	if err != nil {
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

func (s *grpcServer) Histories(ctx context.Context, req *emptypb.Empty) (*generate.HistoryList, error) {
	_, resp, err := s.histories.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp.(*generate.HistoryList), nil
}

func decodeGRPCCreateHistoryRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*generate.Transaction)
	var items []*models.Item

	for _, item := range req.Items {
		each := models.Item{
			Name:      item.Name,
			Quantity:  item.Quantity,
			SubTotal:  item.SubTotal,
			HistoryId: item.HistoryId,
		}
		items = append(items, &each)
	}

	return models.CreateRequest{
		Items:      items,
		GrandTotal: req.GrandTotal,
		CreditCard: req.CreditCard,
	}, nil
}

func encodeGRPCCreateHistoryResponse(_ context.Context, response interface{}) (interface{}, error) {
	return nil, nil
}

func decodeGRPCHistoryRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return nil, nil
}

func encodeGRPCHistoryResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*models.HistoryList)

	if resp == nil {
		return &generate.HistoryList{}, nil
	}

	var histories []*generate.History

	for _, v := range resp.List {

		var items []*generate.Item

		for _, item := range v.Items {
			eachItem := generate.Item{
				Name:      item.Name,
				Quantity:  item.Quantity,
				SubTotal:  item.SubTotal,
				HistoryId: item.HistoryId,
			}

			items = append(items, &eachItem)
		}

		each := generate.History{
			Id:         int32(v.Id),
			GrandTotal: v.GrandTotal,
			CreditCard: v.CreditCard,
			Items:      items,
		}

		histories = append(histories, &each)
	}

	return &generate.HistoryList{List: histories}, nil
}
