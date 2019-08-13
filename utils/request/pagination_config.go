package request

import (
	"fmt"
	"github.com/jinzhu/now"
	"strconv"
	"strings"
)

const (
	IdType     string = "ID"
	NumberType string = "NUMBER"
	StringType string = "STRING"
	DateType   string = "DATE"
)

type PaginationConfig interface {
	Limit() int
	Offset() int
	Order() string
	SearchClause() SearchStruct
}

type SearchStruct struct {
	Query string
	Args  []interface{}
}

func BuildLimit(conditions map[string][]string) int {
	res := 20
	if len(conditions["limit"]) > 0 {
		res, _ = strconv.Atoi(conditions["limit"][0])
	}
	return res
}

func BuildOffset(conditions map[string][]string) int {
	res := 0
	if len(conditions["offset"]) > 0 {
		res, _ = strconv.Atoi(conditions["offset"][0])
	}
	return res
}

func BuildOrder(conditions map[string][]string) string {
	var orders string
	if len(conditions["sort"]) > 0 {
		orders = strings.Join(conditions["sort"], ",")
		return orders
	}
	return ""
}

func BuildSearchClause(conditions map[string][]string, filterable map[string]string) SearchStruct {
	var query string = "true "
	var args []interface{}

	for name, value := range filterable {
		if len(conditions[name]) > 0 {
			switch value {
			case IdType:
				query += fmt.Sprint("AND ", name, " IN (?) ")
				args = append(args, conditions[name])
			case StringType:
				query += fmt.Sprint("AND ", name, " LIKE %?% ")
				args = append(args, conditions[name][0])
			case NumberType:
				query += fmt.Sprint("AND ", name, " BETWEEN ? AND ? ")
				minmax := strings.Split(conditions[name][0], ",")
				args = append(args, minmax[0], minmax[1])
			case DateType:
				query += fmt.Sprint("AND ", name, " BETWEEN ? AND ? ")
				minmax := strings.Split(conditions[name][0], ",")
				min, _ := now.Parse(minmax[0])
				max, _ := now.Parse(minmax[1])
				args = append(args, now.New(min).BeginningOfDay(), now.New(max).EndOfDay())
			}
		}
	}
	return SearchStruct{
		query, args,
	}
}
