package validator

import (
	"log"

	validator "gopkg.in/go-playground/validator.v9"
)

type ErrorListValidation struct {
	Validation []ErrorListItem `json:"validation"`
	Entity     string          `json:"entity"`
}

type ErrorListItem struct {
	Field string `json:"field"`
	Key   string `json:"key"`
}

func IsRequestValid[T any](s T) (bool, *ErrorListValidation) {
	validate := validator.New()
	var errorList ErrorListValidation

	err := validate.Struct(s)
	if err == nil {
		return true, nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		log.Printf("an error occurred tue to %+v", err)
	}

	errEntity := err.(validator.ValidationErrors)

	if errEntity != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errorList.Validation = append(errorList.Validation, ErrorListItem{
				Field: e.StructField(),
				Key:   e.Tag(),
			})
		}

		if len(errorList.Validation) > 0 {
			return false, &errorList
		}
	}

	return true, nil
}
