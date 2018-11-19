package model

import "time"

type Activity struct {
	ID           int64     `json:"id" xorm:"id notnull pk autoincr"`
	Title        string    `json:"title" xorm:"title varchar(255) notnull"`
	Country      string    `json:"country" xorm:"country varchar(40)"`
	Province     string    `json:"province" xorm:"province varchar(40)"`
	City         string    `json:"city" xorm:"city varchar(40)"`
	DetailURL    string    `json:"detail_url" xorm:"detail_url varchar(255)"`
	PublicityIMG string    `json:"publicity_img" xorm:"publicity_img varchar(255)"`
	CreatedAt    time.Time `json:"created" xorm:"created"`
	DeletedAt    time.Time `json:"-" xorm:"deleted"`
}

func (a *Activity) TableName() string {
	return "activity"
}

func init() {
	register(&Activity{})
}

func (a *Activity) CreateActivity(detail *ActiveDetail) error {
	session := NewSession()
	defer session.Close()

	err := session.Begin()

	_, err = session.Insert(a)
	if err != nil {
		session.Rollback()
		return err
	}

	detail.ID = a.ID
	_, err = session.Insert(detail)
	if err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

func (a *Activity) List(from, count int) (int64, []*Activity) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(a)
	if n == 0 {
		return 0, list
	}
	if from == 0 {
		engine.Desc("id").Limit(count, 0).Find(&list)
	} else {
		engine.Where("id < ?", from).Desc("id").Limit(count, 0).Find(&list)
	}
	return n, list
}
