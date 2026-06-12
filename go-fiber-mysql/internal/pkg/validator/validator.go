package validator

import (
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/exception"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	// "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Validator struct {
	V  *validator.Validate
	DB *gorm.DB
}

func NewValidator(db *gorm.DB) *Validator {
	v := &Validator{
		V:  validator.New(validator.WithRequiredStructEnabled()),
		DB: db,
	}

	v.registerCustomRules()

	return v
}

func (val *Validator) registerCustomRules() {
	var err error

	if err = val.V.RegisterValidation("exists", existsValidator(val.DB)); err != nil {
		// logrus.Fatalf("exists validation error: %v", err)
	}

	if err = val.V.RegisterValidation("unique", uniqueValidator(val.DB)); err != nil {
		// logrus.Fatalf("unique validation error: %v", err)
	}

	if err := val.V.RegisterValidation("enum", enumValidator()); err != nil {
		// logrus.Fatalf("enum validation error: %v", err)
	}

	if err := val.V.RegisterValidation("length", LengthValidator()); err != nil {
		// logrus.Fatalf("length validation error: %v", err)
	}
}

func (val *Validator) Struct(s any) error {
	return val.V.Struct(s)
}

// FormatValidationErrors converts validator errors into a readable map
func (val *Validator) FormatValidationErrors(err error, obj any) map[string][]string {
	out := make(map[string][]string)
	if err == nil {
		return out
	}

	ves, ok := err.(validator.ValidationErrors)
	if !ok {
		out["error"] = []string{err.Error()}
		return out
	}

	for _, fe := range ves {
		fieldPath := jsonFieldPath(obj, fe)
		msg := val.translateError(fe)
		out[fieldPath] = append(out[fieldPath], msg)
	}

	// t := reflectType(obj)
	// for _, fe := range ves {
	// 	jsonName := jsonFieldName(t, fe.StructField())
	// 	msg := val.translateError(fe)
	// 	out[jsonName] = append(out[jsonName], msg)
	// }
	return out
}

func (val *Validator) translateError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "max":
		switch fe.Kind() {
		case reflect.String:
			return "maximum " + fe.Param() + " characters"
		default:
			return "maximum value is " + fe.Param()
		}

	case "min":
		switch fe.Kind() {
		case reflect.String:
			return "minimum " + fe.Param() + " characters"
		default:
			return "minimum value is " + fe.Param()
		}
	case "unique":
		return strings.ToLower(fe.Field()) + " already exist"
	case "exists":
		return "must exist"
	case "len":
		return "must be " + fe.Param() + " characters long"
	case "enum":
		enum, ok := fe.Value().(enums.Enum)
		if !ok {
			return "invalid value"
		}

		return "must be one of " + strings.Join(enum.GetValues(), ", ")
	case "length":
		parts := strings.Split(fe.Param(), ".")
		if len(parts) != 2 {
			return "invalid format for length validator"
		}

		return "must be " + parts[0] + " to " + parts[1] + " characters long"
	default:
		return fe.Error()
	}
}

// Reflection helpers
func reflectType(obj any) reflect.Type {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func jsonFieldPath(obj any, fe validator.FieldError) string {
	ns := fe.Namespace() // Users[0].ID

	parts := strings.Split(ns, ".")
	t := reflectType(obj)

	var path []string

	for i, part := range parts {
		// array index, contoh: Users[0]
		if strings.Contains(part, "[") {
			field := part[:strings.Index(part, "[")]
			index := part[strings.Index(part, "["):]

			if f, ok := t.FieldByName(field); ok {
				jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
				path = append(path, jsonTag+index)
				// path = append(path, jsonTag+"."+strconv.Itoa(i))
				t = f.Type.Elem()
			}
			continue
		}

		if i == 0 {
			continue // skip root struct
		}

		if f, ok := t.FieldByName(part); ok {
			jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
			path = append(path, jsonTag)
			t = f.Type
		}
	}

	return strings.Join(path, ".")
}

func jsonFieldName(t reflect.Type, fieldName string) string {
	jsonName := strings.ToLower(fieldName)
	if f, found := t.FieldByName(fieldName); found {
		tag := f.Tag.Get("json")
		if tag != "" {
			jsonName = strings.Split(tag, ",")[0]
		}
	}
	return jsonName
}

func ValidateBody[T any](ctx fiber.Ctx, v *Validator) (*T, error) {
	var body T

	var err error
	if err = ctx.Bind().Body(&body); err != nil {
		return nil, exception.NewBadRequestException(err.Error())
	}

	if err = v.Struct(&body); err != nil {
		return nil, exception.NewValidationException(v.FormatValidationErrors(err, &body))
	}

	return &body, nil
}
