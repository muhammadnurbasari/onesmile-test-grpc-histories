package serviceHistory

import (
	"context"
	"errors"
	"fmt"

	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ServiceHistory interface {
	Create(ctx context.Context, param *models.CreateRequest) error
	Histories(context.Context) (*models.HistoryList, error)
}

func NewServiceHistory(Connection *gorm.DB) ServiceHistory {
	var s ServiceHistory
	{
		s = NewBasicServiceHistory(Connection)
	}

	return s
}

func NewBasicServiceHistory(Connection *gorm.DB) ServiceHistory {
	return &basicServiceHistory{Connection}
}

type basicServiceHistory struct {
	Connection *gorm.DB
}

func (s *basicServiceHistory) Create(ctx context.Context, param *models.CreateRequest) error {
	tx := s.Connection.Begin()

	histories := models.Histories{
		CreditCard: param.CreditCard,
		GrandTotal: uint64(param.GrandTotal),
	}

	result := tx.Create(&histories)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	id := histories.Id

	var historyDetails []*models.HistoryDetails

	for _, v := range param.Items {
		each := models.HistoryDetails{
			HistoryId: id,
			Name:      v.Name,
			Quantity:  uint64(v.Quantity),
			SubTotal:  uint64(v.SubTotal),
		}
		historyDetails = append(historyDetails, &each)
	}

	result = tx.Create(historyDetails)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	log.Info().Msg("success : transaction has been created")
	return nil
}

func (s *basicServiceHistory) Histories(ctx context.Context) (*models.HistoryList, error) {
	rowsHistories, err := s.Connection.Select("id, grand_total, credit_card").Table("histories").Rows()

	if err != nil {
		return nil, err
	}

	var histories []*models.History
	var ids []uint64
	for rowsHistories.Next() {
		var (
			each models.History
			id   uint64
		)
		err = rowsHistories.Scan(&each.Id, &each.GrandTotal, &each.CreditCard)

		if err != nil {
			return nil, errors.New("test rows histories : " + err.Error())
		}

		histories = append(histories, &each)

		id = uint64(each.Id)

		ids = append(ids, id)

	}

	rowsDetails, err := s.Connection.Select("name, quantity, sub_total, history_id").Table("history_details").Where("history_id IN ?", ids).Rows()
	if err != nil {
		return nil, err
	}

	var items []*models.Item
	for rowsDetails.Next() {
		var each *models.Item

		err = rowsDetails.Scan(&each.Name, &each.Quantity, &each.SubTotal, &each.HistoryId)

		if err != nil {
			return nil, errors.New("test rows histories : " + err.Error())
		}

		items = append(items, each)

	}

	fmt.Println(histories)

	for key, history := range histories {
		var itemsFix []*models.Item
		for _, item := range items {
			var each models.Item
			if int64(history.Id) == item.HistoryId {
				each.Name = item.Name
				each.SubTotal = item.SubTotal
				each.Quantity = item.Quantity
				each.HistoryId = item.HistoryId
				itemsFix = append(itemsFix, &each)
			}
		}

		// history.Items = items
		histories[key].Items = itemsFix
	}

	var result models.HistoryList
	result.List = histories

	log.Info().Msg("success : get transaction history successfully")
	return &result, nil
}
