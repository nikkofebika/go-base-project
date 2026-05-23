package helpers

import (
	"math"

	"github.com/gofiber/fiber/v3"
)

type Response[T any] struct {
	Data    *T                  `json:"data,omitempty"`
	Meta    *Meta               `json:"meta,omitempty"`
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	From        int   `json:"from"`
	PerPage     int   `json:"per_page"`
	To          int   `json:"to"`
	Total       int64 `json:"total"` // use int64 because return db.Count() from gorm is int64
	TotalPages  int   `json:"total_pages"`
}

func NewMeta(page, perPage int, total int64) *Meta {
	offset := (page - 1) * perPage

	return &Meta{
		CurrentPage: page,
		From:        min(offset+1, int(total)),
		PerPage:     perPage,
		To:          min(offset+perPage, int(total)),
		Total:       total,
		TotalPages:  int(math.Ceil(float64(total) / float64(perPage))),
	}
}

func NewResponse[T any](c fiber.Ctx, data T) error {
	return c.Status(fiber.StatusOK).JSON(Response[T]{
		Data: &data,
	})
}

func NewResponsePagination[T any](c fiber.Ctx, data []T, meta *Meta) error {
	return c.JSON(Response[[]T]{
		Data: &data,
		Meta: meta,
	})
}

func NewResponseMessage(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response[string]{
		Message: message,
	})
}

func NewResponseErrors(c fiber.Ctx, status int, errs map[string][]string) error {
	// If err implements the error interface, use its Error() message.
	// var errMsg any = err

	// if e, ok := err.(error); ok {
	// 	errMsg = e.Error()
	// }

	// // Log original error for debugging.
	// fmt.Println("NewResponseError NewResponseError", err)

	return c.Status(status).JSON(Response[any]{
		Errors: errs,
	})
}

func NewResponseCreated(c fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(Response[string]{
		Message: "Data created successfully",
	})
}

// NewResponseCreatedWithData returns created status with data and message
func NewResponseCreatedWithData[T any](c fiber.Ctx, data T, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response[T]{
		Data:    &data,
		Message: message,
	})
}

func NewResponseUpdated(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response[string]{
		Message: "Data updated successfully",
	})
}

// NewResponseUpdatedWithData returns updated status with data and message
func NewResponseUpdatedWithData[T any](c fiber.Ctx, data T, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response[T]{
		Data:    &data,
		Message: message,
	})
}

func NewResponseDeleted(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response[string]{
		Message: "Data deleted successfully",
	})
}

func NewResponseForceDeleted(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response[string]{
		Message: "Data force deleted successfully",
	})
}

func NewResponseRestored(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response[string]{
		Message: "Data restored successfully",
	})
}
