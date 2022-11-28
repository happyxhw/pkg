package query

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var (
	re = regexp.MustCompile("[^a-zA-Z0-9_-]") // 替换特殊字符, 防止注入
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
	OnlyCount bool   `query:"-"`
	Page      int    `query:"page" validate:"min=1"`
	Size      int    `query:"size" validate:"max=100"`
	Filter    string `query:"filter"`  // filter=xxx
	SortBy    string `query:"sort_by"` // sort_by=-last_modified,+email
}

type Opt struct {
	Fields []string
}

func Fields(fields ...string) Opt {
	return Opt{
		Fields: fields,
	}
}

type FieldFunc func(string) string

func ParseOrder(sortBy string, handle ...FieldFunc) string {
	items := strings.Split(sortBy, ",")
	orders := make([]string, 0, len(items))

	for _, item := range items {
		if len(item) < 2 {
			continue
		}
		direction, key := item[0], item[1:]
		if direction != '+' && direction != '-' {
			continue
		}
		key = re.ReplaceAllString(key, "")
		if key == "" {
			continue
		}
		if len(handle) > 0 {
			key = handle[0](key)
			if key == "" {
				continue
			}
		}
		if direction == '+' {
			orders = append(orders, key+" ASC")
		} else {
			orders = append(orders, key+" DESC")
		}
	}

	return strings.Join(orders, ",")
}

func WrapPageQuery(db *gorm.DB, pp Param, out interface{}) (*PagingResult, error) {
	if pp.OnlyCount {
		var count int64
		err := db.Count(&count).Error
		if err != nil {
			return nil, err
		}
		return &PagingResult{Total: count}, nil
	}

	total, err := findPage(db, pp, out)
	if err != nil {
		return nil, err
	}

	return &PagingResult{
		Total: total,
		Index: pp.Page,
		Size:  pp.Size,
	}, nil
}

func findPage(db *gorm.DB, pp Param, out interface{}) (int64, error) {
	// 0. query count
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	pageIndex, pageSize := pp.Page, pp.Size
	if count == 0 || (pageIndex-1)*pageSize >= int(count) {
		return count, nil
	}
	queryDB := db.Session(&gorm.Session{})
	// 1. query id
	var ids []int
	err = queryDB.Offset((pageIndex-1)*pageSize).Limit(pageSize).Pluck("id", &ids).Error
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}
	// 2. query rows
	err = queryDB.Where("id IN (?)", ids).Find(out).Error

	return count, err
}

func Take(tx *gorm.DB, opt Opt, out interface{}) error {
	if len(opt.Fields) > 0 {
		tx = tx.Select(opt.Fields)
	}
	err := tx.Take(out).Error

	return err
}
