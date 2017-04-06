package handlers

import (
	"encoding/json"
	"net/http"

	"reflect"
	"strconv"

	"context"
	"mallfin_api/db"
	"mallfin_api/logging"
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

func writeJSON(ctx context.Context, w http.ResponseWriter, resp interface{}, status int) {
	logger := logging.FromContext(ctx)
	b, err := json.Marshal(resp)
	if err != nil {
		logger.WithField("resp", resp).Errorf("Cannot serialize response to json: %s", err)
		internalErrorResponse(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		logger.Errorf("Cannot write json response: %s", err)
	}

}

func errorResponse(ctx context.Context, w http.ResponseWriter, errorCode, details string, status int) {
	errObj := ErrorData{Code: errorCode, Details: details, Status: status}
	resp := ErrorResponse{Error: &errObj}
	writeJSON(ctx, w, resp, status)
}

func notFoundResponse(ctx context.Context, w http.ResponseWriter, errorCode string) {
	errorResponse(ctx, w, errorCode, errorCode, http.StatusNotFound)
}

func internalErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func response(ctx context.Context, w http.ResponseWriter, data interface{}) {
	resp := SuccessResponse{Data: data}
	writeJSON(ctx, w, resp, http.StatusOK)
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

func paginateResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, resultsList interface{}, totalCount int, limit, offset *int) {
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
	response(ctx, w, data)
}

func totalCountFromResults(resultsLen int, limit, offset *int) (int, bool) {
	if (limit == nil || *limit == 0) && (offset == nil || *offset == 0 || resultsLen != 0) {
		totalCount := resultsLen
		if offset != nil {
			totalCount += *offset
		}
		return totalCount, true
	}
	return 0, false
}

func checkCity(ctx context.Context, w http.ResponseWriter, cityID *int) bool {
	if cityID != nil {
		logger := logging.FromContext(ctx)
		logger.Info("Check city in db")
		exists, err := db.IsCityExists(*cityID)
		if err != nil {
			logger.Errorf("Cannot check city: %s", err)
			internalErrorResponse(w)
			return false
		}
		if !exists {
			logger.WithField("city_id", *cityID).Warn("City does not exist")
			notFoundResponse(ctx, w, CITY_NOT_FOUND)
			return false
		}
	}
	return true
}

func checkShop(ctx context.Context, w http.ResponseWriter, shopID int) bool {
	logger := logging.FromContext(ctx)
	exists, err := db.IsShopExists(shopID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		notFoundResponse(ctx, w, SHOP_NOT_FOUND)
		return false
	}
	return true
}

func checkSubwayStation(ctx context.Context, w http.ResponseWriter, stationID int) bool {
	logger := logging.FromContext(ctx)
	exists, err := db.IsSubwayStationExists(stationID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		notFoundResponse(ctx, w, SUBWAY_STATION_NOT_FOUND)
		return false
	}
	return true
}

func checkCategory(ctx context.Context, w http.ResponseWriter, categoryID int) bool {
	logger := logging.FromContext(ctx)
	exists, err := db.IsCategoryExists(categoryID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		notFoundResponse(ctx, w, CATEGORY_NOT_FOUND)
		return false
	}
	return true
}

func checkMall(ctx context.Context, w http.ResponseWriter, mallID int) bool {
	logger := logging.FromContext(ctx)
	exists, err := db.IsMallExists(mallID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return false
	}
	if !exists {
		notFoundResponse(ctx, w, MALL_NOT_FOUND)
		return false
	}
	return true
}
