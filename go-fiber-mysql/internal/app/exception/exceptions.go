package exception

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
	"gorm.io/gorm"
)

type BaseException struct {
	StatusCode int                 `json:"-"`
	Message    string              `json:"message"`
	Errors     map[string][]string `json:"errors"`
}

func (be *BaseException) Error() string {
	return be.Message
}

func NewHttpException(statusCode int, message ...string) *BaseException {
	if statusCode < 100 || statusCode > 599 {
		statusCode = fiber.StatusInternalServerError
	}

	defaultMessage := utils.StatusMessage(statusCode)

	if len(message) > 1 && message[0] != "" {
		defaultMessage = message[0]
	}

	return &BaseException{
		StatusCode: statusCode,
		Message:    defaultMessage,
	}
}

func NewBadRequestException(message ...string) *BaseException {
	statusCode := fiber.StatusBadRequest
	defaultMessage := utils.StatusMessage(statusCode)

	if len(message) > 0 && message[0] != "" {
		defaultMessage = message[0]
	}

	return &BaseException{
		StatusCode: statusCode,
		Message:    defaultMessage,
	}
}

func NewNotFoundException(message ...string) *BaseException {
	statusCode := fiber.StatusNotFound
	defaultMessage := utils.StatusMessage(statusCode)

	if len(message) > 0 && message[0] != "" {
		defaultMessage = message[0]
	}

	return &BaseException{
		StatusCode: statusCode,
		Message:    defaultMessage,
	}
}

func NewInternalServerException(message ...string) *BaseException {
	statusCode := fiber.StatusInternalServerError
	defaultMessage := utils.StatusMessage(statusCode)

	if len(message) > 0 && message[0] != "" {
		defaultMessage = message[0]
	}

	return &BaseException{
		StatusCode: statusCode,
		Message:    defaultMessage,
	}
}

func NewUnauthorizedException() *BaseException {
	statusCode := fiber.StatusUnauthorized
	return &BaseException{
		StatusCode: statusCode,
		Message:    utils.StatusMessage(statusCode),
	}
}

func NewForbiddenException() *BaseException {
	statusCode := fiber.StatusForbidden
	return &BaseException{
		StatusCode: statusCode,
		Message:    utils.StatusMessage(statusCode),
	}
}

type ValidationException struct {
	*BaseException
}

func NewValidationException(errors map[string][]string) *ValidationException {
	statusCode := fiber.StatusUnprocessableEntity
	return &ValidationException{
		BaseException: &BaseException{
			StatusCode: statusCode,
			Message:    utils.StatusMessage(statusCode),
			Errors:     errors,
		},
	}
}

func NewFileRequiredException(keys ...string) *ValidationException {
	errs := make(map[string][]string)
	for _, key := range keys {
		errs[key] = []string{key + " is required"}
	}

	return NewValidationException(errs)
}

// database exception
var gormErrors = map[error]int{
	gorm.ErrRecordNotFound:                404,
	gorm.ErrInvalidTransaction:            500,
	gorm.ErrNotImplemented:                501,
	gorm.ErrMissingWhereClause:            400,
	gorm.ErrUnsupportedRelation:           400,
	gorm.ErrPrimaryKeyRequired:            400,
	gorm.ErrModelValueRequired:            400,
	gorm.ErrModelAccessibleFieldsRequired: 400,
	gorm.ErrSubQueryRequired:              400,
	gorm.ErrInvalidData:                   400,
	gorm.ErrUnsupportedDriver:             500,
	gorm.ErrRegistered:                    400,
	gorm.ErrInvalidField:                  400,
	gorm.ErrEmptySlice:                    400,
	gorm.ErrDryRunModeUnsupported:         500,
	gorm.ErrInvalidDB:                     500,
	gorm.ErrInvalidValue:                  400,
	gorm.ErrInvalidValueOfLength:          400,
	gorm.ErrPreloadNotAllowed:             400,
	gorm.ErrDuplicatedKey:                 409,
	gorm.ErrForeignKeyViolated:            409,
	gorm.ErrCheckConstraintViolated:       400,
}

func NewDatabaseException(err error) *BaseException {
	for gormError, statusCode := range gormErrors {
		if errors.Is(err, gormError) {
			if statusCode == 404 {
				return NewNotFoundException()
			}

			if statusCode < 500 {
				return &BaseException{
					StatusCode: statusCode,
					Message:    utils.StatusMessage(statusCode),
				}
			}

			return &BaseException{
				StatusCode: statusCode,
				Message:    utils.StatusMessage(statusCode),
			}
		}
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return &BaseException{StatusCode: 409, Message: utils.StatusMessage(409)}
	}

	return &BaseException{StatusCode: 500, Message: utils.StatusMessage(500)}

	// var mysqlErr *mysql.MySQLError
	// if errors.As(err, &mysqlErr) {
	// 	switch mysqlErr.Number {
	// 	// case 1146: // Table doesn't exist
	// 	// 	return &BaseException{StatusCode: 500, Message: mysqlErr.Message}
	// 	case 1062: // Duplicate entry
	// 		return &BaseException{StatusCode: 409, Message: mysqlErr.Message}
	// 	default:
	// 		return &BaseException{StatusCode: 500, Message: mysqlErr.Error()}
	// 	}
	// }

	// return &BaseException{StatusCode: 500, Message: err.Error()}
}
