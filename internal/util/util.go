package util

import (
	"fmt"
	"net/http"
)

func SliceContainsE[T comparable](arr []T, a T) bool {
	for _, x := range arr {
		if x == a {
			return true
		}
	}
	return false
}

func HTTPErrWrap(w http.ResponseWriter, err string, code int, actErr *string) {
	if actErr != nil {
		err = fmt.Sprintf("%s: %s\n", err, *actErr)
	}
	w.WriteHeader(code)
	w.Write([]byte(err))

	return
}
