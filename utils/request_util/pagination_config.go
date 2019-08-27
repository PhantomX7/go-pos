package request_util

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
	BoolType   string = "BOOL"
	DateType   string = "DATE"
)

type PaginationConfig interface {
	Limit() int
	Offset() int
	Order() string
	SearchClause() SearchStruct
}

type SearchStruct struct {
	// default
	Query string
	Args  []interface{}
}

type Config struct {
	limit        int
	offset       int
	order        string
	searchClause SearchStruct
}

func (c Config) Limit() (res int) {
	return c.limit
}

func (c Config) Order() string {
	return c.order
}

func (c Config) Offset() (res int) {
	return c.offset
}

func (c Config) SearchClause() (res SearchStruct) {
	return c.searchClause
}

func GeneratePaginationConfig(limit int, offset int, order string, searchStruct SearchStruct) PaginationConfig {
	paginationConfig := Config{
		limit:        limit,
		offset:       offset,
		order:        order,
		searchClause: searchStruct,
	}

	return paginationConfig
}

func GenerateDefaultSearchStruct() SearchStruct {
	return SearchStruct{
		Query: "true",
		Args:  []interface{}{},
	}
}

func GenerateDefaultPaginationConfig() PaginationConfig {
	return GeneratePaginationConfig(0, 0, "", GenerateDefaultSearchStruct())
}

func BuildLimit(conditions map[string][]string) int {
	res := 20
	if len(conditions["limit"]) > 0 {
		res, _ = strconv.Atoi(conditions["limit"][0])
		if res > 100 {
			res = 100
		}
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
				query += fmt.Sprint("AND ", name, " LIKE ? ")
				args = append(args, fmt.Sprint("%", conditions[name][0], "%"))
			case BoolType:
				boolean := 0
				if conditions[name][0] == "true" {
					boolean = 1
				}
				query += fmt.Sprint("AND ", name, " = ? ")
				args = append(args, boolean)
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
