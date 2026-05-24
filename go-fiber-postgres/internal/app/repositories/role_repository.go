package repositories

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/pkg/request"
	"time"

	"gorm.io/gorm"
)

type RoleRepository interface {
	DB() *gorm.DB
	WithTx(db *gorm.DB) RoleRepository
	FindAll(db *gorm.DB, page, perPage int) ([]entities.Role, int64, error)
	FindOne(id uint64, includes []request.AppliedInclude) (entities.Role, error)
	Create(data *entities.Role, permissionIDs []uint) error
	Update(id uint64, data map[string]any, permissionIDs []uint) error
	Delete(id, userID uint64) error
	ForceDelete(id uint64) error
	Restore(id uint64) error
	SyncPermissions(id uint64, permissionIDs []uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) DB() *gorm.DB {
	return r.db.Model(&entities.Role{})
}

func (r *roleRepository) WithTx(tx *gorm.DB) RoleRepository {
	return &roleRepository{db: tx}
}

func (r *roleRepository) FindAll(db *gorm.DB, page, perPage int) ([]entities.Role, int64, error) {
	var datas []entities.Role
	var total int64

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := db.Limit(perPage).Offset(offset).Find(&datas).Error; err != nil {
		return nil, 0, err
	}

	return datas, total, nil
}

func (r *roleRepository) FindOne(id uint64, includes []request.AppliedInclude) (entities.Role, error) {
	var data entities.Role

	db := r.DB()
	for _, include := range includes {
		db = include.Apply(db)
	}

	err := db.Take(&data, id).Error
	return data, err
}

func (r *roleRepository) Create(data *entities.Role, permissionIDs []uint) error {
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

func (r *roleRepository) Update(id uint64, data map[string]any, permissionIDs []uint) error {
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

func (r *roleRepository) Delete(id, userID uint64) error {
	result := r.DB().Where(request.ParamID+"=?", id).Updates(map[string]any{
		entities.DeletedAt:   gorm.DeletedAt{Time: time.Now(), Valid: true},
		entities.DeletedByID: userID,
	})

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *roleRepository) ForceDelete(id uint64) error {
	result := r.db.Unscoped().Delete(&entities.Role{}, id)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *roleRepository) Restore(id uint64) error {
	result := r.DB().Where(request.ParamID+"=?", id).UpdateColumns(map[string]any{
		entities.DeletedByID: nil,
	})

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *roleRepository) SyncPermissions(id uint64, permissionIDs []uint) error {
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
