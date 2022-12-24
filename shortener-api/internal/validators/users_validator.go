package validators

import "net/mail"

func ValidateUserName(username string) bool {
	return validateStringLength(username, 25, 3)
}

func ValidateUserPassword(password string) bool {
	return validateStringLength(password, 50, 5)
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validateStringLength(value string, maxLength int, minLength int) bool {
	valueLength := len(value)
	return valueLength <= maxLength && valueLength >= minLength
}
