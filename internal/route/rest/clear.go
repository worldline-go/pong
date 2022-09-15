package rest

import "strings"

func cleanMethod(method string) string {
	method = strings.TrimSpace(method)

	if method == "" {
		return "GET"
	}

	return strings.TrimSpace(strings.ToUpper(method))
}
