package services

import (
	"context"
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/app/helpers"
	ctxHelper "go-fiber-mysql/internal/app/helpers/context"
	"go-fiber-mysql/internal/app/repositories"
	"go-fiber-mysql/internal/app/requests"
	"go-fiber-mysql/internal/pkg/request"
	"time"
)

type PermissionService interface {
	FindAll(ctx context.Context, request *request.PaginationRequest) ([]entities.Permission, *helpers.Meta, error)
	FindOne(ctx context.Context, id uint, includes []request.AppliedInclude) (entities.Permission, error)
	Create(ctx context.Context, request *requests.PermissionCreateRequest) error
	Update(ctx context.Context, id uint, request *requests.PermissionUpdateRequest) error
	Delete(ctx context.Context, id uint) error
	ForceDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
}

type permissionService struct {
	repository repositories.PermissionRepository
}

func NewPermissionService(repository repositories.PermissionRepository) PermissionService {
	return &permissionService{
		repository: repository,
	}
}

func (s *permissionService) FindAll(ctx context.Context, req *request.PaginationRequest) ([]entities.Permission, *helpers.Meta, error) {
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

func (s *permissionService) FindOne(ctx context.Context, id uint, includes []request.AppliedInclude) (entities.Permission, error) {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return entities.Permission{}, err
	}

	data, err := s.repository.FindOne(id, includes)
	if err != nil {
		return entities.Permission{}, exception.NewDatabaseException(err)
	}

	return data, nil
}

func (s *permissionService) Create(ctx context.Context, request *requests.PermissionCreateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	data := entities.Permission{
		Name:               request.Name,
		Slug:               request.Slug,
		AuditCreatedEntity: entities.AuditCreatedEntity{CreatedByID: &userLoggedIn.ID},
	}

	return s.repository.Create(&data)
}

func (s *permissionService) Update(ctx context.Context, id uint, req *requests.PermissionUpdateRequest) error {
	userLoggedIn, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	updates := entities.UpdatedFields{
		entities.UpdatedAt:   time.Now(),
		entities.UpdatedByID: userLoggedIn.ID,
	}

	helpers.AddToMapIfNotNil(updates, "name", req.Name)
	helpers.AddToMapIfNotNil(updates, "slug", req.Slug)

	return s.repository.Update(id, updates)
}

func (s *permissionService) Delete(ctx context.Context, id uint) error {
	user, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Delete(id, user.ID)
}

func (s *permissionService) ForceDelete(ctx context.Context, id uint) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.ForceDelete(id)
}

func (s *permissionService) Restore(ctx context.Context, id uint) error {
	_, err := ctxHelper.ExtractUserFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repository.Restore(id)
}
