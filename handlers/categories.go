package handlers

import (
	"net/http"

	"mallfin_api/db"
	"mallfin_api/models"
	"mallfin_api/serializers"

	"mallfin_api/logging"

	"github.com/gazoon/binding"
	"github.com/gazoon/httprouter"
)

func CategoriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := categoriesListForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	sorting := formData.Sort
	if !checkCity(ctx, w, formData.City, "log prefix") {
		return
	}
	var categories []*models.Category
	if formData.Shop != nil {
		shopID := *formData.Shop
		if !checkShop(ctx, w, shopID, "log prefix") {
			return
		}
	} else {
		var err error
		categories, err = db.GetCategories(sorting)
		if err != nil {
			logger.Error(err)
			internalErrorResponse(w)
			return
		}
	}
	serialized := serializers.SerializeCategories(categories)
	response(ctx, w, serialized)
}

func CategoryDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	formData := categoryDetailsForm{}
	errs := binding.Form(r, &formData)
	if errs != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, errs.Error(), http.StatusBadRequest)
		return
	}
	categoryID, err := ps.ByNameInt("id")
	if err != nil {
		errorResponse(ctx, w, INCORRECT_REQUEST_DATA, err.Error(), http.StatusBadRequest)
		return
	}
	cityID := formData.City
	if !checkCity(ctx, w, cityID, "log prefix") {
		return
	}
	category, err := db.GetCategoryDetails(categoryID)
	if err != nil {
		logger.Error(err)
		internalErrorResponse(w)
		return
	}
	if category == nil {
		notFoundResponse(ctx, w, CATEGORY_NOT_FOUND)
		return
	}
	serialized := serializers.SerializeCategory(category)
	response(ctx, w, serialized)
}
