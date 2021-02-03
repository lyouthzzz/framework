package gormutil

import (
	"context"
	"gorm.io/gorm"
	"math"
)

type Param struct {
	DB       *gorm.DB
	PageNum  int
	PageSize int
	OrderBy  []string
	Debug    bool
}

type Paginator struct {
	TotalCount int64       `json:"total_count"`
	TotalPage  int         `json:"total_page"`
	Items      interface{} `json:"items"`
	Offset     int         `json:"offset"`
	PageNum    int         `json:"page_num"`
	PageSize   int         `json:"page_size"`
	PagePrev   int         `json:"page_prev"`
	PageNext   int         `json:"page_next"`
}

func Paging(ctx context.Context, param *Param, items interface{}) (*Paginator, error) {
	var db = param.DB

	if param.Debug {
		db = db.Debug()
	}

	var (
		paginator = &Paginator{PageNum: param.PageNum, PageSize: param.PageSize}
		listTx    = db.WithContext(ctx)
		countDone = make(chan error, 1)
		err       error
	)

	go count(ctx, db, items, countDone, &paginator.TotalCount)

	if paginator.PageNum < 1 {
		paginator.PageNum = 1
	}
	if paginator.PageSize == 0 {
		paginator.PageSize = 10
	}

	if len(param.OrderBy) > 0 {
		for _, o := range param.OrderBy {
			listTx = listTx.Order(o)
		}
	}

	if paginator.PageNum == 1 {
		paginator.Offset = 0
	} else {
		paginator.Offset = (paginator.PageNum - 1) * paginator.PageSize
	}

	err = listTx.Limit(paginator.PageSize).Offset(paginator.Offset).Find(items).Error
	if err != nil {
		return nil, err
	}

	if err = <-countDone; err != nil {
		return nil, err
	}

	paginator.Items = items
	paginator.TotalPage = int(math.Ceil(float64(paginator.TotalCount) / float64(paginator.PageSize)))

	if paginator.PageNum > 1 {
		paginator.PagePrev = paginator.PageNum - 1
	} else {
		paginator.PagePrev = paginator.PageNum
	}
	if paginator.PageNum == paginator.TotalPage {
		paginator.PageNext = paginator.PageNum
	} else {
		paginator.PageNext = paginator.PageNum + 1
	}
	return paginator, err
}

func count(ctx context.Context, db *gorm.DB, anyType interface{}, done chan error, count *int64) {
	done <- db.WithContext(ctx).Model(anyType).Count(count).Error
}
