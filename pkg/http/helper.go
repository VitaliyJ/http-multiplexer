package http

import "net/http"

func fillByOkCodes(urls []string) map[string]int {
	m := make(map[string]int, len(urls))
	for i := range urls {
		m[urls[i]] = http.StatusOK
	}

	return m
}
