package repositories

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/pkg/request"
	"time"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	DB() *gorm.DB
	WithTx(db *gorm.DB) PermissionRepository
	FindAll(db *gorm.DB, page, perPage int) ([]entities.Permission, int64, error)
	FindOne(id uint64, includes []request.AppliedInclude) (entities.Permission, error)
	Create(data *entities.Permission) error
	Update(id uint64, data map[string]any) error
	Delete(id, userID uint64) error
	ForceDelete(id uint64) error
	Restore(id uint64) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) DB() *gorm.DB {
	return r.db.Model(&entities.Permission{})
}

func (r *permissionRepository) WithTx(tx *gorm.DB) PermissionRepository {
	return &permissionRepository{db: tx}
}

func (r *permissionRepository) FindAll(db *gorm.DB, page, perPage int) ([]entities.Permission, int64, error) {
	var datas []entities.Permission
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

func (r *permissionRepository) FindOne(id uint64, includes []request.AppliedInclude) (entities.Permission, error) {
	var data entities.Permission

	db := r.DB()
	for _, include := range includes {
		db = include.Apply(db)
	}

	err := db.Take(&data, id).Error
	return data, err
}

func (r *permissionRepository) Create(data *entities.Permission) error {
	return r.db.Create(data).Error
}

func (r *permissionRepository) Update(id uint64, data map[string]any) error {
	result := r.DB().Where(request.ParamID+"=?", id).Updates(data)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *permissionRepository) Delete(id, userID uint64) error {
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

func (r *permissionRepository) ForceDelete(id uint64) error {
	result := r.db.Unscoped().Delete(&entities.Permission{}, id)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *permissionRepository) Restore(id uint64) error {
	result := r.DB().Where(request.ParamID+"=?", id).UpdateColumns(map[string]any{
		entities.DeletedAt:   nil,
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
