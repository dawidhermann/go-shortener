// App validator
package validate

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate structure using tags
func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err != nil {
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		var fieldErrors ValidationFieldErrors
		for _, verr := range verrors {
			field := ValidationFieldError{
				Field: verr.Field(),
				Error: verr.Error(),
			}
			fieldErrors = append(fieldErrors, field)
		}
		return fieldErrors
	}
	return nil
}
