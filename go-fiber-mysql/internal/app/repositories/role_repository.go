package repositories

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/pkg/request"

	"gorm.io/gorm"
)

type RoleRepository interface {
	DB() *gorm.DB
	WithTx(db *gorm.DB) RoleRepository
	FindAll(db *gorm.DB, page, perPage int) ([]entities.Role, int64, error)
	FindOne(id uint64, includes []request.AppliedInclude) (entities.Role, error)
	Create(data *entities.Role, permissionIDs []uint64) error
	Update(id uint64, data map[string]any, permissionIDs []uint64) error
	Delete(id, userID uint64) error
	ForceDelete(id uint64) error
	Restore(id uint64) error
	SyncPermissions(id uint64, permissionIDs []uint64) error
}

type roleRepository struct {
	*BaseRepository[entities.Role]
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: NewBaseRepository[entities.Role](db),
	}
}

func (r *roleRepository) WithTx(tx *gorm.DB) RoleRepository {
	return &roleRepository{BaseRepository: r.BaseRepository.WithTx(tx)}
}

func (r *roleRepository) Create(data *entities.Role, permissionIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}

		if len(permissionIDs) > 0 {
			var permissions []entities.Permission
			if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(data).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roleRepository) Update(id uint64, data map[string]any, permissionIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entities.Role{}).Where(request.ParamID+"=?", id).Updates(data)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected <= 0 {
			return exception.NewNotFoundException()
		}

		if permissionIDs != nil {
			var role entities.Role
			if err := tx.First(&role, id).Error; err != nil {
				return err
			}

			var permissions []entities.Permission
			if len(permissionIDs) > 0 {
				if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roleRepository) SyncPermissions(id uint64, permissionIDs []uint64) error {
	var role entities.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return err
	}

	var permissions []entities.Permission
	if len(permissionIDs) > 0 {
		if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}
	}

	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}
