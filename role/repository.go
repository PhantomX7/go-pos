package role

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request"
)

type RoleRepository interface {
	Insert(role *models.Role) error
	Update(role *models.Role) error
	FindAll(config request.PaginationConfig) ([]models.Role, error)
	FindByID(roleID int64) (models.Role, error)
	FindByName(roleName string) (models.Role, error)
}
