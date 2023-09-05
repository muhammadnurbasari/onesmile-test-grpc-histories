package models

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

type CreateRequest struct {
	Items      []*Item
	GrandTotal int64
	CreditCard string
}

type Item struct {
	Name      string
	Quantity  int32
	SubTotal  int64
	HistoryId int64
}

type HistoryList struct {
	List []*History
}

type History struct {
	Id         uint64
	Items      []*Item
	GrandTotal int64
	CreditCard string
}
