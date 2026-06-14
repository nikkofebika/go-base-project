package repositories

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/pkg/request"

	"gorm.io/gorm"
)

type UserRepository interface {
	DB() *gorm.DB
	WithTx(db *gorm.DB) UserRepository
	FindAll(db *gorm.DB, page, perPage int) ([]entities.User, int64, error)
	FindOne(id uint64, includes []request.AppliedInclude) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	Create(data *entities.User) error
	Update(id uint64, data map[string]any) error
	Delete(id, userID uint64) error
	ForceDelete(id uint64) error
	Restore(id uint64) error
	SyncRoles(id uint64, roleIDs []uint64) error
	GetPermissionByUserID(userID uint64) ([]string, error)
}

type userRepository struct {
	*BaseRepository[entities.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[entities.User](db),
	}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{BaseRepository: r.BaseRepository.WithTx(tx)}
}

func (r *userRepository) FindOneByEmail(email string) (entities.User, error) {
	var data entities.User
	err := r.db.Where("email=?", email).Take(&data).Error
	return data, err
}

func (r *userRepository) SyncRoles(id uint64, roleIDs []uint64) error {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}

	var roles []entities.Role
	if len(roleIDs) > 0 {
		if err := r.db.Find(&roles, roleIDs).Error; err != nil {
			return err
		}
	}

	return r.db.Model(&user).Association("Roles").Replace(roles)
}

func (r *userRepository) GetPermissionByUserID(userID uint64) ([]string, error) {
	var names []string

	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Distinct("permissions.name").
		Pluck("permissions.name", &names).Error

	return names, err
}
