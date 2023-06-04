package validate

import "encoding/json"

// Error that represents information about validation error
type ValidationFieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Define new type to group ValidationError
type ValidationFieldErrors []ValidationFieldError

// Implement error interface
func (vfe ValidationFieldErrors) Error() string {
	data, err := json.Marshal(vfe)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
