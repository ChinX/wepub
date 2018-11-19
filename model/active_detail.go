package model

import "time"

type ActiveDetail struct {
	ID             int64     `json:"id" xorm:"id notnull pk"`
	Price          int       `json:"price" xorm:"price notnull default 0"`
	Final          int       `json:"final" xorm:"final notnull default 0"`
	Quantity       int       `json:"quantity" xorm:"quantity notnull default 0"`
	Total          int64     `json:"total" xorm:"Total notnull default 0"`
	Completed      int64     `json:"completed" xorm:"completed notnull default 0"`
	DailyTotal     int64     `json:"daily_total" xorm:"daily_total notnull default 0"`
	DailyCompleted int64     `json:"daily_completed" xorm:"daily_completed notnull default 0"`
	ExpireDate     time.Time `json:"expire_date" xorm:"expire notnull"`
}

func (a *ActiveDetail) TableName() string {
	return "active_detail"
}

func init() {
	register(&ActiveDetail{})
}

func (a *ActiveDetail) List(conditions []interface{}) []interface{} {
	return In(*a, "id", conditions)
}
