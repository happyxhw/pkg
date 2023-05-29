package query

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var (
	re = regexp.MustCompile("[^a-zA-Z0-9_-]") // 替换特殊字符, 防止注入
)

type FieldFunc func(string) string
type FilterFunc func(*gorm.DB, string, string) bool

func ParseOrder(sortBy string, handles ...FieldFunc) string {
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
		// 校验排序key的合法性
		if len(handles) > 0 {
			if key = handles[0](key); key == "" {
				continue
			}
		} else if key = re.ReplaceAllString(key, ""); key == "" { // 替换特殊字符
			continue
		}
		if direction == '+' {
			orders = append(orders, key+" ASC")
		} else {
			orders = append(orders, key+" DESC")
		}
	}
	if len(orders) == 0 {
		return "id DESC"
	}

	return strings.Join(orders, ",")
}

var logicalOp = map[string]string{
	"and": "and",
	"or":  "or",
}

var compOp = map[string]string{
	"eq":     "=",
	"ne":     "!=",
	"le":     "<=",
	"lt":     "<",
	"ge":     ">=",
	"gt":     ">",
	"in":     "in",
	"nin":    "not in",
	"like":   "like",
	"starts": "starts",
	"ends":   "ends",
}

var compRangeOp = map[string]bool{
	"in":  true,
	"nin": true,
}

// or(and(eq(name,h),le(age,10),or(eq(name,y),le(age,11))),or(in(status,0,1,2),nin(role,1,2,3)))
func ParseFilter(tx *gorm.DB, filter string, handles ...FilterFunc) (string, []interface{}, error) { //nolint:gocyclo,gocritic
	var stack []string
	var params []interface{}
	var err error
	var gOp string
	for i := 0; i < len(filter); {
		c := filter[i]
		if c == '(' {
			var j int
			for j = i - 1; j >= 0; j-- {
				if filter[j] == ',' || filter[j] == '(' {
					break
				}
			}
			op := strings.ToLower(strings.TrimSpace(filter[j+1 : i]))
			if logicalOp[op] != "" {
				stack = append(stack, op)
				i++
				continue
			}
			if compOp[op] != "" {
				for j = i; j < len(filter); j++ {
					if filter[j] != ')' {
						continue
					}
					gOp, params, err = parseComp(tx, op, filter[i+1:j], params, handles...)
					if err != nil {
						return "", nil, err
					}
					stack = append(stack, gOp)
					break
				}
				i = j + 1
				continue
			} else {
				return "", nil, errors.New("incorrect filter syntax: " + filter[0:i])
			}
		} else if c == ')' {
			if logicalOp[stack[len(stack)-1]] != "" {
				return "", nil, errors.New("incorrect filter syntax: " + filter[0:i])
			}
			var t []string
			for j := len(stack) - 1; j >= 0; j-- {
				n := len(stack)
				op := stack[n-1]
				stack = stack[0 : n-1]
				if logicalOp[op] != "" {
					stack = append(stack, "("+strings.Join(t, " "+logicalOp[op]+" ")+")")
					break
				} else {
					t = append(t, op)
				}
			}
		}
		i++
	}
	if len(stack) != 1 {
		return "", nil, errors.New("incorrect filter syntax: " + filter)
	}
	for i, j := 0, len(params)-1; i < j; i, j = i+1, j-1 {
		params[i], params[j] = params[j], params[i]
	}
	return stack[0], params, nil
}

// 转化为 gorm 语句
//
//nolint:gocritic
func parseComp(tx *gorm.DB, op, filter string, params []interface{},
	handles ...FilterFunc) (string, []interface{}, error) {
	t := strings.Split(filter, ",")
	if len(t) < 2 {
		return "", params, errors.New("incorrect filter syntax: " + filter)
	}
	col := strings.TrimSpace(t[0])
	if len(handles) > 0 {
		if !handles[0](tx, col, op) {
			return "", nil, fmt.Errorf("(%s, %s) not supported", col, op)
		}
	}
	switch {
	case compRangeOp[op]:
		var p []interface{}
		for i := 1; i < len(t); i++ {
			p = append(p, strings.TrimSpace(t[i]))
		}
		params = append(params, p)
		return fmt.Sprintf("%s %s ?", col, compOp[op]), params, nil
	case compOp[op] == "like":
		params = append(params, "%"+strings.TrimSpace(t[1])+"%")
		return fmt.Sprintf("%s %s ?", col, compOp[op]), params, nil
	case compOp[op] == "starts":
		params = append(params, strings.TrimSpace(t[1])+"%")
		return fmt.Sprintf("%s like ?", col), params, nil
	case compOp[op] == "ends":
		params = append(params, "%"+strings.TrimSpace(t[1]))
		return fmt.Sprintf("%s like?", col), params, nil
	default:
		params = append(params, strings.TrimSpace(t[1]))
		return fmt.Sprintf("%s %s ?", col, compOp[op]), params, nil
	}
}
