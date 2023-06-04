package validators

import "net/url"

func ValidateUrl(rawUrl string) error {
	_, err := url.ParseRequestURI(rawUrl)
	return err
}
