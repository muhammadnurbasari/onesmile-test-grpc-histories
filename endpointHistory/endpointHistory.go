package endpointHistory

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/models"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/serviceHistory"
)

type SetHistory struct {
	CreateHistoryEndpoint endpoint.Endpoint
	HistoryEndpoint       endpoint.Endpoint
}

func NewEndpointHistory(svc serviceHistory.ServiceHistory) SetHistory {
	var createHistoryEndpoint endpoint.Endpoint
	{
		createHistoryEndpoint = MakeCreateHistoryEndpoint(svc)
	}

	var historyEndpoint endpoint.Endpoint
	{
		historyEndpoint = MakeHistoryEndpoint(svc)
	}

	return SetHistory{
		CreateHistoryEndpoint: createHistoryEndpoint,
		HistoryEndpoint:       historyEndpoint,
	}
}

func (s *SetHistory) Create(ctx context.Context, param *models.CreateRequest) error {
	_, err := s.CreateHistoryEndpoint(ctx, param)

	if err != nil {
		return err
	}

	return nil
}

func (s *SetHistory) Histories(ctx context.Context) (*models.HistoryList, error) {
	resp, err := s.HistoryEndpoint(ctx, nil)

	if err != nil {
		return nil, err
	}

	response := resp.(models.HistoryList)

	return &response, nil
}

func MakeCreateHistoryEndpoint(svc serviceHistory.ServiceHistory) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.CreateRequest)
		err = svc.Create(ctx, &req)
		return nil, err
	}
}

func MakeHistoryEndpoint(svc serviceHistory.ServiceHistory) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := svc.Histories(ctx)
		return resp, err
	}
}
