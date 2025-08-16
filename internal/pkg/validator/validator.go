package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

// Init initializes the validator instance
func Init() {
	validate = validator.New()
}

// ValidateStruct validates any struct with tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
