package repositories

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/pkg/request"

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
	*BaseRepository[entities.Permission]
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		BaseRepository: NewBaseRepository[entities.Permission](db),
	}
}

func (r *permissionRepository) WithTx(tx *gorm.DB) PermissionRepository {
	return &permissionRepository{BaseRepository: r.BaseRepository.WithTx(tx)}
}
