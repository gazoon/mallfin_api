package handlers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
	"mallfin_api/db/models"
)

type mallsListForm struct {
	City  int
	Shop  int
	Query string
	Ids   []int
}

func (mp *mallsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&mp.City: binding.Field{
			Form:     "city",
			Required: true,
		},
		&mp.Shop: binding.Field{
			Form:     "shop",
			Required: true,
		},
		&mp.Query: "query",
		&mp.Ids:   "ids",
	}
}

func MallsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := mallsListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INVALID_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	log.Info("success")
}
func MallDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mallID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INVALID_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	mall := models.GetMall(mallID)
	if mall == nil {
		errorResponse(w, MALL_NOT_FOUND, "Mall with such id does not exists", http.StatusNotFound)
		return
	}
	response(w, mall)
}
