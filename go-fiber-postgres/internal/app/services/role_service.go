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

type RoleService interface {
	FindAll(ctx context.Context, request *request.PaginationRequest) ([]entities.Role, *helpers.Meta, error)
	FindOne(ctx context.Context, id uint64, includes []request.AppliedInclude) (entities.Role, error)
	Create(ctx context.Context, request *requests.RoleCreateRequest) error
	Update(ctx context.Context, id uint64, request *requests.RoleUpdateRequest) error
	Delete(ctx context.Context, id uint64) error
	ForceDelete(ctx context.Context, id uint64) error
	Restore(ctx context.Context, id uint64) error
	SyncPermissions(ctx context.Context, id uint64, request *requests.RoleSyncPermissionsRequest) error
}

type roleService struct {
	repository repositories.RoleRepository
}

func NewRoleService(repository repositories.RoleRepository) RoleService {
	return &roleService{
		repository: repository,
	}
}

func (s *roleService) FindAll(ctx context.Context, req *request.PaginationRequest) ([]entities.Role, *helpers.Meta, error) {
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

func (s *roleService) FindOne(ctx context.Context, id uint64, includes []request.AppliedInclude) (entities.Role, error) {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return entities.Role{}, err
	}

	data, err := s.repository.FindOne(id, includes)
	if err != nil {
		return entities.Role{}, exception.NewDatabaseException(err)
	}

	return data, nil
}

func (s *roleService) Create(ctx context.Context, request *requests.RoleCreateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	data := entities.Role{
		Name:               request.Name,
		AuditCreatedEntity: entities.AuditCreatedEntity{CreatedByIDEntity: entities.CreatedByIDEntity{CreatedByID: &userLoggedIn.ID}},
	}

	return s.repository.Create(&data, request.PermissionIDs)
}

func (s *roleService) Update(ctx context.Context, id uint64, req *requests.RoleUpdateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	updates := entities.UpdatedFields{
		entities.UpdatedAt:   time.Now(),
		entities.UpdatedByID: userLoggedIn.ID,
	}

	helpers.AddToMapIfNotNil(updates, "name", req.Name)

	return s.repository.Update(id, updates, req.PermissionIDs)
}

func (s *roleService) Delete(ctx context.Context, id uint64) error {
	user, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Delete(id, user.ID)
}

func (s *roleService) ForceDelete(ctx context.Context, id uint64) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.ForceDelete(id)
}

func (s *roleService) Restore(ctx context.Context, id uint64) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Restore(id)
}

func (s *roleService) SyncPermissions(ctx context.Context, id uint64, req *requests.RoleSyncPermissionsRequest) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.repository.SyncPermissions(id, req.PermissionIDs)
	if err != nil {
		return exception.NewDatabaseException(err)
	}

	return nil
}
