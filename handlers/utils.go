package handlers

import (
	"encoding/json"
	"net/http"

	"reflect"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

const (
	INCORRECT_REQUEST_DATA   = "INCORRECT_REQUEST_DATA"
	INTERNAL_ERROR           = "INTERNAL_ERROR"
	MALL_NOT_FOUND           = "MALL_NOT_FOUND"
	CITY_NOT_FOUND           = "CITY_NOT_FOUND"
	SUBWAY_STATION_NOT_FOUND = "SUBWAY_STATION_NOT_FOUND"
	SHOP_NOT_FOUND           = "SHOP_NOT_FOUND"
	CATEGORY_NOT_FOUND       = "CATEGORY_NOT_FOUND"
)
const DoesNotExistMsg = "%s with such id does not exists."

type JSONObject map[string]interface{}

type ErrorData struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Status  int    `json:"status_code"`
}

type ErrorResponse struct {
	Error *ErrorData `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type PaginationData struct {
	Count      int         `json:"count"`
	TotalCount int         `json:"total_count"`
	Next       *string     `json:"next"`
	Prev       *string     `json:"prev"`
	Results    interface{} `json:"results"`
}

func writeJSON(w http.ResponseWriter, resp interface{}, status int) {
	b, err := json.Marshal(resp)
	if err != nil {
		log.Panicf("Cannot serialize response to json: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func errorResponse(w http.ResponseWriter, errorCode, details string, status int) {
	errObj := ErrorData{Code: errorCode, Details: details, Status: status}
	resp := ErrorResponse{Error: &errObj}
	writeJSON(w, resp, status)
}
func internalErrorResponse(w http.ResponseWriter) {
	errorResponse(w, INTERNAL_ERROR, "An internal server error occurred, please try again later.", http.StatusInternalServerError)
}
func response(w http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{Data: data}
	writeJSON(w, resp, http.StatusOK)
}

func nextPage(totalCount, limit, offset int) (int, int, bool) {
	if limit+offset >= totalCount {
		return 0, 0, false
	}
	nextOffset := limit + offset
	nextLimit := limit
	if nextOffset+nextLimit > totalCount {
		nextLimit = totalCount - nextOffset
	}
	return nextLimit, nextOffset, true
}

func prevPage(totalCount, limit, offset int) (int, int, bool) {
	if offset == 0 {
		return 0, 0, false
	}
	var prevOffset int
	var prevLimit int
	if offset < limit {
		prevOffset = 0
		prevLimit = offset
	} else {
		prevOffset = offset - limit
		prevLimit = limit
	}
	return prevLimit, prevOffset, true

}

func pageURL(r *http.Request, limit, offset int) string {
	url := r.URL
	params := url.Query()
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))
	url.RawQuery = params.Encode()
	return url.String()
}

func paginateResponse(w http.ResponseWriter, r *http.Request, resultsList interface{}, totalCount int, limit, offset *int) {
	limitValue := totalCount
	if limit != nil {
		limitValue = *limit
	}
	offsetValue := 0
	if offset != nil {
		offsetValue = *offset
	}
	var nextPageURL *string = nil
	if nextLimit, nextOffset, ok := nextPage(totalCount, limitValue, offsetValue); ok {
		url := pageURL(r, nextLimit, nextOffset)
		nextPageURL = &url
	}
	var prevPageURL *string = nil
	if prevLimit, prevOffset, ok := prevPage(totalCount, limitValue, offsetValue); ok {
		url := pageURL(r, prevLimit, prevOffset)
		prevPageURL = &url
	}
	data := &PaginationData{
		TotalCount: totalCount,
		Count:      reflect.ValueOf(resultsList).Len(),
		Results:    resultsList,
		Next:       nextPageURL,
		Prev:       prevPageURL,
	}
	response(w, data)
}
