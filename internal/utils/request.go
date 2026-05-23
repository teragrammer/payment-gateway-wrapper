package utils

import "net/http"

func ExtractPaginationFromRequest(r *http.Request) (int, int) {
	pageStr := StringToInt(r.URL.Query().Get("page"), 1)
	limitStr := StringToInt(r.URL.Query().Get("limit"), 10)
	return pageStr, limitStr
}
