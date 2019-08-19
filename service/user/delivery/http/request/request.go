package request

import (
	"github.com/PhantomX7/go-pos/utils/request_util"
)

// request related struct
type UserCreateRequest struct {
	Username string `form:"username" binding:"required,unique=users.username"`
	Password string `form:"password" binding:"required"`
	RoleId   int    `form:"role_id" binding:"required,exist=roles.id"`
}

type UserUpdateRequest struct {
	Username string `form:"username" binding:"omitempty,unique=users.username"`
	Password string `form:"password"`
	RoleId   int    `form:"role_id" binding:"omitempty,exist=roles.id"`
}

type UserPaginationConfig struct {
	limit        int
	offset       int
	order        string
	searchClause request_util.SearchStruct
}

func NewUserPaginationConfig(conditions map[string][]string) UserPaginationConfig {
	filterable := map[string]string{
		"id":         request_util.IdType,
		"username":   request_util.StringType,
		"role_id":    request_util.IdType,
		"created_at": request_util.DateType,
	}

	userPaginationConfig := UserPaginationConfig{
		limit:        request_util.BuildLimit(conditions),
		offset:       request_util.BuildOffset(conditions),
		order:        request_util.BuildOrder(conditions),
		searchClause: request_util.BuildSearchClause(conditions, filterable),
	}

	return userPaginationConfig
}

func (u UserPaginationConfig) Limit() (res int) {
	return u.limit
}

func (u UserPaginationConfig) Order() string {
	return u.order
}

func (u UserPaginationConfig) Offset() (res int) {
	return u.offset
}

func (u UserPaginationConfig) SearchClause() (res request_util.SearchStruct) {
	return u.searchClause
}
