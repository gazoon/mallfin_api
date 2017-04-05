package handlers

import (
	"net/http"

	"mallfin_api/db"
	"mallfin_api/models"
	"mallfin_api/serializers"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
)

func CategoriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	formData := categoriesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	if !checkCity(w, formData.City, "log prefix") {
		return
	}
	var categories []*models.Category
	if formData.Shop != nil {
		shopID := *formData.Shop
		if !checkShop(w,shopID, "log prefix") {
			return
		}
	} else {
		var err error
		categories, err = db.GetCategories(sorting)
		if err != nil {
			log.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeCategories(categories)
	response(w, serialized)
}

func CategoryDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	formData := categoryDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	categoryID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
	if !checkCity(w, cityID, "log prefix") {
		return
	}
	category, err := db.GetCategoryDetails(categoryID)
	if err != nil {
		log.Error(err)
		internalErrorResponse(w)
		return
	}
	if category == nil {
		notFoundResponse(w, CATEGORY_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeCategory(category)
	response(w, serialized)
}
