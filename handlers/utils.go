package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type JSONObject map[string]interface{}
type ErrorObject struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Status  int    `json:"status_code"`
}
type ErrorResponse struct {
	Error *ErrorObject `json:"error"`
}
type SuccessResponse struct {
	Data interface{} `json:"data"`
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
	errObj := ErrorObject{Code: errorCode, Details: details, Status: status}
	resp := ErrorResponse{Error: &errObj}
	writeJSON(w, resp, status)
}
func response(w http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{Data: data}
	writeJSON(w, resp, http.StatusOK)
}
