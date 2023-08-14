package histories

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TransactionsServer struct {
	Connection *gorm.DB
}

type Histories struct {
	Id         uint64
	CreditCard string
	GrandTotal uint64
}

type HistoryDetails struct {
	HistoryId uint64
	Name      string
	Quantity  uint64
	SubTotal  uint64
}

func (t TransactionsServer) Create(ctx context.Context, param *generate.Transaction) (*empty.Empty, error) {
	tx := t.Connection.Begin()

	histories := Histories{
		CreditCard: param.CreditCard,
		GrandTotal: uint64(param.GrandTotal),
	}

	result := tx.Create(&histories)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	id := histories.Id

	var historyDetails []*HistoryDetails

	for _, v := range param.Items {
		each := HistoryDetails{
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
		return nil, result.Error
	}

	tx.Commit()

	log.Info().Msg("success : transaction has been created")
	return new(empty.Empty), nil

}

func (t TransactionsServer) Histories(context.Context, *empty.Empty) (*generate.HistoryList, error) {
	rowsHistories, err := t.Connection.Select("id, grand_total, credit_card").Table("histories").Rows()

	if err != nil {
		return nil, err
	}

	var histories []*generate.History
	var ids []uint64
	for rowsHistories.Next() {
		var (
			each generate.History
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

	// var whereIn = "("

	// for key, id := range ids {

	// 	if key == len(ids)-1 {
	// 		whereIn += fmt.Sprintf("%d", id.Id)
	// 	} else {

	// 		whereIn += fmt.Sprintf("%d,", id.Id)
	// 	}
	// }
	// whereIn += ")"

	rowsDetails, err := t.Connection.Select("name, quantity, sub_total, history_id").Table("history_details").Where("history_id IN ?", ids).Rows()
	if err != nil {
		return nil, err
	}

	var items []generate.Item
	for rowsDetails.Next() {
		var each generate.Item

		err = rowsDetails.Scan(&each.Name, &each.Quantity, &each.SubTotal, &each.HistoryId)

		if err != nil {
			return nil, errors.New("test rows histories : " + err.Error())
		}

		items = append(items, each)

	}

	fmt.Println(histories)

	for key, history := range histories {
		var itemsFix []*generate.Item
		for _, item := range items {
			var each generate.Item
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

	var result generate.HistoryList
	result.List = histories

	log.Info().Msg("success : get transaction history successfully")
	return &result, nil
}
