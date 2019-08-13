package request

import (
	"fmt"
	"strconv"

	"github.com/PhantomX7/go-pos/models"
)

// request related struct

type PaginationConfig interface {
	Limit() int
	Offset() int
	Order() string
	SearchClause() map[string][]string
}

type UserCreateRequest struct {
	Username string `form:"username" binding:"required,unique=users.username"`
	Password string `form:"password" binding:"required"`
	RoleId   int    `form:"role_id" binding:"required,exist=roles.id"`
}

func (request UserCreateRequest) ToUserModel() models.User {
	return models.User{
		Username: request.Username,
		Password: request.Password,
		RoleId:   int64(request.RoleId),
	}
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
	searchClause map[string][]string
}

func NewUserPaginationConfig(conditions map[string][]string) UserPaginationConfig {
	userPaginationConfig := UserPaginationConfig{
		limit:        buildLimit(conditions),
		offset:       buildOffset(conditions),
		order:        buildOrder(conditions),
		searchClause: buildSearchClause(conditions),
	}

	return userPaginationConfig
}

func (a UserPaginationConfig) Limit() (res int) {
	return a.limit
}

func (a UserPaginationConfig) Order() string {
	return a.order
}

func (a UserPaginationConfig) Offset() (res int) {
	return a.offset
}

func (a UserPaginationConfig) SearchClause() (res map[string][]string) {
	return a.searchClause
}

func buildLimit(conditions map[string][]string) (res int) {
	if len(conditions["limit"]) > 0 {
		res, _ = strconv.Atoi(conditions["limit"][0])
	}

	return
}

func buildOffset(conditions map[string][]string) (res int) {
	if len(conditions["offset"]) > 0 {
		res, _ = strconv.Atoi(conditions["offset"][0])
	}
	return
}

func buildOrder(conditions map[string][]string) string {
	var order string
	if len(conditions["sort"]) > 0 {
		order = conditions["sort"][0]
	}

	orderCol, orderDir := "", ""

	if len(order) > 2 {
		if order[0:1] == "-" {
			orderDir = "desc"
			order = order[1:]
		} else {
			orderDir = "asc"
		}

		if order == "created_at" || order == "order" {
			orderCol = fmt.Sprintf("`%s`", order)
		}
		return orderCol + " " + orderDir
	}

	return ""
}

func buildSearchClause(conditions map[string][]string) (res map[string][]string) {
	res = make(map[string][]string)
	if len(conditions["status"]) > 0 {
		//status := conditions["status"][0]
		//if status == "active" {
		//	res["status"] = mysql.ActiveBannerStatusCondition
		//
		//} else if status == "inactive" {
		//	res["status"] = mysql.InactiveBannerStatusCondition
		//}
	}

	if len(conditions["type"]) > 0 {
		//typeCode := models.BannerTypeString(conditions["type"][0]).ToCode()
		//if typeCode < models.BannerType(len(models.BannerTypeArray)) {
		//	res["type"] = models.BannerTypeString(conditions["type"][0]).ToCode()
		//}
	}

	if len(conditions["title"]) > 0 {
		res["title"] = conditions["title"]
	}

	return
}
