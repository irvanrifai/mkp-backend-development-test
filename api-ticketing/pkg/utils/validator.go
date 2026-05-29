package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register fungsi untuk mengambil nama field dari tag JSON
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate berfungsi mengecek struct dan mengembalikan map error yang rapi
func Validate(s interface{}) map[string]string {
	errors := make(map[string]string)
	err := validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = getCustomMessage(err)
		}
	}
	return errors
}

// getCustomMessage adalah tempat kita mengatur pesan error (mirip Laravel/Joi)
func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return "Minimal karakter adalah " + fe.Param()
	case "max":
		return "Maksimal karakter adalah " + fe.Param()
	}
	return fe.Error() // Default error
}
