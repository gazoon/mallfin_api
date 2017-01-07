package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"reflect"
)

const (
	INVALID_REQUEST_DATA     = "INVALID_REQUEST_DATA"
	MALL_NOT_FOUND           = "MALL_NOT_FOUND"
	CITY_NOT_FOUND           = "CITY_NOT_FOUND"
	SUBWAY_STATION_NOT_FOUND = "SUBWAY_STATION_NOT_FOUND"
	SHOP_NOT_FOUND           = "SHOP_NOT_FOUND"
	CATEGORY_NOT_FOUND       = "CATEGORY_NOT_FOUND"
)
const DOES_NOT_EXISTS_MSG = "%s with such id does not exists."

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
type ListData struct {
	Count      int         `json:"count"`
	TotalCount int         `json:"total_count"`
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
func response(w http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{Data: data}
	writeJSON(w, resp, http.StatusOK)
}
func objectResponse(w http.ResponseWriter, object interface{}) {
	data := object
	response(w, data)
}
func listResponse(w http.ResponseWriter, resultsList interface{}, totalCount int) {
	data := &ListData{
		TotalCount: totalCount,
		Count:      reflect.ValueOf(resultsList).Len(),
		Results:    resultsList,
	}
	response(w, data)
}
