package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

var Validate = validator.New()

func BindAndValidateReqBody[T any] (c *gin.Context) (*T, error) {
	var dto T
	if err := c.BindJSON(&dto); err != nil {
		return nil, err
	}

	validationErr := Validate.Struct(dto)
	if validationErr != nil {
		return nil, validationErr
	}

	return &dto, nil
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", fieldName)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", fieldName)
			default:
				message = fmt.Sprintf("%s is invalid", fieldName)
			}

			errors = append(errors, ValidationError{
				Field: fieldName,
				Error: message,
			})
		}
	} else{	//i.e not a validator error, json binding error
		errors = append(errors, ValidationError{
			Field: "json binding",
			Error: err.Error(),
		})
	}

	return errors
}
