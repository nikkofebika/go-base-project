package repositories

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/pkg/request"
	"time"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db    *gorm.DB
	model T
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	var model T
	return &BaseRepository[T]{db: db, model: model}
}

func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db.Model(&r.model)
}

func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: tx, model: r.model}
}

func (r *BaseRepository[T]) FindAll(db *gorm.DB, page, perPage int) ([]T, int64, error) {
	var datas []T
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

func (r *BaseRepository[T]) FindOne(id uint64, includes []request.AppliedInclude) (T, error) {
	var data T

	db := r.DB()
	for _, include := range includes {
		db = include.Apply(db)
	}

	err := db.Take(&data, id).Error

	return data, err
}

func (r *BaseRepository[T]) Create(data *T) error {
	return r.db.Create(data).Error
}

func (r *BaseRepository[T]) Update(id uint64, data map[string]any) error {
	result := r.DB().Where(request.ParamID+"=?", id).Updates(data)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *BaseRepository[T]) Delete(id, userID uint64) error {
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

func (r *BaseRepository[T]) ForceDelete(id uint64) error {
	result := r.db.Unscoped().Delete(new(T), id)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *BaseRepository[T]) Restore(id uint64) error {
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
