package repositories

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/pkg/request"
	"time"

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
	SyncRoles(id uint64, roleIDs []uint) error
	GetPermissionByUserID(userID uint64) ([]string, error)

	GetRefreshToken(token string) (entities.Token, error)
	DeleteRefreshTokenByUserID(id uint64) error
	DeleteRefreshToken(refreshToken string) error
	SaveRefreshToken(refreshToken *entities.Token) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) DB() *gorm.DB {
	return r.db.Model(&entities.User{})
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) FindAll(db *gorm.DB, page, perPage int) ([]entities.User, int64, error) {
	var datas []entities.User
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

func (r *userRepository) FindOne(id uint64, includes []request.AppliedInclude) (entities.User, error) {
	var data entities.User

	db := r.DB()

	for _, include := range includes {
		db = include.Apply(db)
	}

	err := db.Take(&data, id).Error

	return data, err
}

func (r *userRepository) FindOneByEmail(email string) (entities.User, error) {
	var data entities.User

	err := r.db.Where("email=?", email).Take(&data).Error

	return data, err
}

func (r *userRepository) Create(data *entities.User) error {
	return r.db.Create(data).Error
}

func (r *userRepository) Update(id uint64, data map[string]any) error {
	result := r.DB().Where(request.ParamID+"=?", id).Updates(data)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *userRepository) Delete(id, userID uint64) error {
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

func (r *userRepository) ForceDelete(id uint64) error {
	result := r.db.Unscoped().Delete(&entities.User{}, id)

	if result.RowsAffected <= 0 {
		return exception.NewNotFoundException()
	}

	if result.Error != nil {
		return exception.NewDatabaseException(result.Error)
	}

	return nil
}

func (r *userRepository) Restore(id uint64) error {
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

func (r *userRepository) SyncRoles(id uint64, roleIDs []uint) error {
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

// for jwt
func (r *userRepository) GetRefreshToken(token string) (entities.Token, error) {
	var data entities.Token
	err := r.db.Where("token = ?", token).First(&data).Error
	return data, err
}

func (r *userRepository) DeleteRefreshTokenByUserID(id uint64) error {
	return r.db.Where("user_id = ?", id).Delete(&entities.Token{}).Error
}

func (r *userRepository) DeleteRefreshToken(refreshToken string) error {
	return r.db.Where("token = ?", refreshToken).
		Where("type = ?", enums.TokenTypeRefreshToken).
		Delete(&entities.Token{}).Error
}

func (r *userRepository) SaveRefreshToken(refreshToken *entities.Token) error {
	return r.db.Save(refreshToken).Error
}

// for jwt
