package module

import (
	"errors"
	"time"

	"github.com/chinx/wepub/model"
)

type Activity struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Country        string    `json:"country"`
	Province       string    `json:"province"`
	City           string    `json:"city"`
	DetailURL      string    `json:"detail_url"`
	PublicityIMG   string    `json:"publicity_img"`
	CreatedAt      time.Time `json:"created"`
	Price          int       `json:"price"`
	Final          int       `json:"final"`
	Quantity       int       `json:"quantity"`
	Total          int64     `json:"total"`
	Completed      int64     `json:"completed"`
	DailyTotal     int64     `json:"daily_total"`
	DailyCompleted int64     `json:"daily_completed"`
	ExpireDate     time.Time `json:"expire_date"`
}

func CreateActivity()  {
}

func ListActivities(from, count int) (int64, interface{}) {
	list := make([]*Activity, 0, 10)
	activity := &model.Activity{}
	total, activities := activity.List(from, count)
	switch len(activities) {
	case 0:
		return total, list
	case 1:
		detail, err := GetActiveDetail(activities[0].ID)
		if err != nil {
			return 0, list
		}
		return total, append(list, mergeActivityDetail(activities[0], detail))
	default:
		conditions := make([]interface{}, 0, len(activities))
		for i := range activities {
			conditions = append(conditions, activities[i].ID)
		}
		tb := &model.ActiveDetail{}
		results := tb.List(conditions)
		for i := range results {
			detail, ok := results[i].(*model.ActiveDetail)
			if !ok {
				return 0, list
			}
			if detail.ExpireDate.Before(time.Now()) {
				go DeleteActivity(detail.ID)
				continue
			}
			for j := range activities {
				if detail.ID == activities[j].ID {
					list = append(list, mergeActivityDetail(activities[j], detail))
				}
			}
		}
		return total, list
	}
}

func GetActivity(id int64) (*Activity, error) {
	activity := &model.Activity{ID: id}
	if ok := model.Get(activity); !ok {
		return nil, errors.New("没有相关的活动记录")
	}
	detail, err := GetActiveDetail(id)
	if err != nil {
		return nil, err
	}
	return mergeActivityDetail(activity, detail), nil
}

func GetActiveDetail(id int64) (*model.ActiveDetail, error) {
	detail := &model.ActiveDetail{ID: id}
	if ok := model.Get(detail); !ok {
		return nil, errors.New("没有相关的活动记录")
	}
	if detail.ExpireDate.Before(time.Now()) {
		go DeleteActivity(detail.ID)
		return nil, errors.New("活动已失效")
	}
	return detail, nil
}

func DeleteActivity(id int64) (*model.Activity, error) {
	activity := &model.Activity{ID: id}
	if ok := model.Delete(activity); !ok {
		return nil, errors.New("没有相关的活动记录")
	}
	return activity, nil
}

func mergeActivityDetail(ma *model.Activity, detail *model.ActiveDetail) *Activity {
	return &Activity{
		ID:             ma.ID,
		Title:          ma.Title,
		Country:        ma.Country,
		Province:       ma.Province,
		City:           ma.City,
		DetailURL:      ma.DetailURL,
		PublicityIMG:   ma.PublicityIMG,
		CreatedAt:      ma.CreatedAt,
		ExpireDate:     detail.ExpireDate,
		Price:          detail.Price,
		Final:          detail.Final,
		Quantity:       detail.Quantity,
		Total:          detail.Total,
		Completed:      detail.Completed,
		DailyTotal:     detail.DailyTotal,
		DailyCompleted: detail.DailyCompleted,
	}
}
