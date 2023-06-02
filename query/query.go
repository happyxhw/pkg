package query

import (
	"gorm.io/gorm"
)

type ListResult struct {
	List       interface{}   `json:"list"`
	Pagination *PagingResult `json:"pagination,omitempty"`
}

type PagingResult struct {
	Total int64 `json:"total"`
	Index int   `json:"index"`
	Size  int   `json:"size"`
}

type Param struct {
	Page   int    `query:"page"`
	Size   int    `query:"size" validate:"max=100"`
	Filter string `query:"filter"`
	SortBy string `query:"sort_by"` // sort_by=-last_modified,+email
}

type Opt struct {
	Fields []string
}

func Fields(fields ...string) Opt {
	return Opt{
		Fields: fields,
	}
}

func WrapPageQuery(db *gorm.DB, out interface{}, pp *Param, orderFn FieldFunc, opts ...Opt) (*PagingResult, error) {
	if pp.Page == 0 {
		pp.Page = 1
	}
	if pp.Size == 0 {
		pp.Size = 20
	}

	total, err := findPage(db, out, pp, orderFn, opts...)
	if err != nil {
		return nil, err
	}

	return &PagingResult{
		Total: total,
		Index: pp.Page,
		Size:  pp.Size,
	}, nil
}

func findPage(tx *gorm.DB, out interface{}, pp *Param, orderFn FieldFunc, opts ...Opt) (int64, error) {
	opt := GetOpt(opts...)
	if len(opt.Fields) > 0 {
		tx = tx.Select(opt.Fields)
	}
	// 0. query count
	page, size := pp.Page, pp.Size
	query := tx.Session(&gorm.Session{})
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 || (page-1)*size >= int(count) {
		return count, nil
	}

	// 1. query id
	var ids []int64
	if pp.SortBy != "" {
		if pp.SortBy = ParseOrder(pp.SortBy, orderFn); pp.SortBy != "" {
			query = query.Order(pp.SortBy)
		}
	}
	err = query.Offset((page-1)*size).Limit(size).Pluck("id", &ids).Error
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}

	// 2. query rows
	query = tx.Where("id IN (?)", ids)
	if pp.SortBy != "" {
		query = query.Order(pp.SortBy)
	}
	err = query.Find(out).Error

	return count, err
}

func Take(tx *gorm.DB, out interface{}, opts ...Opt) error {
	opt := GetOpt(opts...)
	if len(opt.Fields) > 0 {
		tx = tx.Select(opt.Fields)
	}
	err := tx.Take(out).Error

	return err
}

func GetOpt(opts ...Opt) Opt {
	if len(opts) > 0 {
		return opts[0]
	}
	return Opt{}
}
