package utils

import (
	"fmt"
	"net/url"
)

func ValidateURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)

	if err != nil {
		fmt.Printf("%s is not a valid URL", urlToCheck)
		return false
	}
	return true
}
