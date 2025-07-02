package utils

import (
	"net/http"
	"strconv"
)

func ParsePageAndLimit(r *http.Request) (int, int) {
	//parse the query param
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 //for default
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}
	return page, limit
}
