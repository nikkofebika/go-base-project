package services

import (
	"context"
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/app/helpers"
	ctxHelper "go-fiber-postgres/internal/app/helpers/context"
	"go-fiber-postgres/internal/app/repositories"
	"go-fiber-postgres/internal/app/requests"
	"go-fiber-postgres/internal/pkg/request"
	"time"
)

type UserService interface {
	FindAll(ctx context.Context, request *request.PaginationRequest) ([]entities.User, *helpers.Meta, error)
	FindOne(ctx context.Context, id uint64, includes []request.AppliedInclude) (entities.User, error)
	Create(ctx context.Context, request *requests.UserCreateRequest) error
	Update(ctx context.Context, id uint64, request *requests.UserUpdateRequest) error
	Delete(ctx context.Context, id uint64) error
	ForceDelete(ctx context.Context, id uint64) error
	Restore(ctx context.Context, id uint64) error
	SyncRoles(ctx context.Context, id uint64, request *requests.UserSyncRolesRequest) error
	GetPermissionByUserID(ctx context.Context, userID uint64) ([]string, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) FindAll(ctx context.Context, req *request.PaginationRequest) ([]entities.User, *helpers.Meta, error) {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	db := s.repository.DB()

	for _, filter := range req.Filter {
		db = filter.Apply(db, filter.Value)
	}

	for _, include := range req.Includes {
		db = include.Apply(db)
	}

	db = request.ApplySorts(db, req.Sorts)

	datas, total, err := s.repository.FindAll(db, req.Page, req.PerPage)
	if err != nil {
		return nil, nil, exception.NewDatabaseException(err)
	}

	meta := helpers.NewMeta(req.Page, req.PerPage, total)

	return datas, meta, nil
}

func (s *userService) FindOne(ctx context.Context, id uint64, includes []request.AppliedInclude) (entities.User, error) {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return entities.User{}, err
	}

	data, err := s.repository.FindOne(id, includes)
	if err != nil {
		return entities.User{}, exception.NewDatabaseException(err)
	}

	return data, nil
}

func (s *userService) Create(ctx context.Context, request *requests.UserCreateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	data := entities.User{
		Name:               request.Name,
		Email:              request.Email,
		Password:           request.Password,
		AuditCreatedEntity: entities.AuditCreatedEntity{CreatedByIDEntity: entities.CreatedByIDEntity{CreatedByID: &userLoggedIn.ID}},
	}

	return s.repository.Create(&data)
}

func (s *userService) Update(ctx context.Context, id uint64, req *requests.UserUpdateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := s.FindOne(ctx, id, nil)
	if err != nil {
		return err
	}

	var errors = make(map[string][]string)
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		existEmail, err := s.repository.FindOneByEmail(*req.Email)
		if err == nil && existEmail.ID != user.ID {
			errors["email"] = append(errors["email"], "Email already exist")
		}
	}

	updates := entities.UpdatedFields{
		entities.UpdatedAt:   time.Now(),
		entities.UpdatedByID: userLoggedIn.ID,
	}

	helpers.AddToMapIfNotNil(updates, "name", req.Name)
	helpers.AddToMapIfNotNil(updates, "email", req.Email)

	return s.repository.Update(id, updates)
}

func (s *userService) Delete(ctx context.Context, id uint64) error {
	user, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Delete(id, user.ID)
}

func (s *userService) ForceDelete(ctx context.Context, id uint64) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.ForceDelete(id)
}

func (s *userService) Restore(ctx context.Context, id uint64) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Restore(id)
}

func (s *userService) SyncRoles(ctx context.Context, id uint64, req *requests.UserSyncRolesRequest) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.repository.SyncRoles(id, req.RoleIDs)
	if err != nil {
		return exception.NewDatabaseException(err)
	}

	return nil
}

func (s *userService) GetPermissionByUserID(ctx context.Context, userID uint64) ([]string, error) {
	// No extraction check here because it might be called by middleware before context is fully set
	// or specifically to check permissions for a given ID.
	return s.repository.GetPermissionByUserID(userID)
}
