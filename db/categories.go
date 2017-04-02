package db

import (
	"mallfin_api/models"
	"mallfin_api/utils"

	"github.com/pkg/errors"
	pg "gopkg.in/pg.v5"
)

func GetCategoryDetails(categoryID int) (*models.Category, error) {
	queryName := utils.CurrentFuncName()
	categories, err := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ?0
	LIMIT 1
	`), categoryID)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, nil
	}
	return categories[0], nil
}

func GetCategories(sorting models.Sorting) ([]*models.Category, error) {
	orderBy := categoryOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	categories, err := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	ORDER BY {order}
	`))
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoriesByIDs(categoryIDs []int) ([]*models.Category, error) {
	categoryIDsArray := pg.Array(categoryIDs)
	queryName := utils.CurrentFuncName()
	categories, err := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ANY (?0)
	`), categoryIDsArray)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoriesByShop(shopID int, sorting models.Sorting) ([]*models.Category, error) {
	orderBy := categoryOrderBy(sorting)
	queryName := utils.CurrentFuncName()
	categories, err := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	  JOIN shop_category sc ON c.category_id = sc.category_id
	WHERE sc.shop_id = ?0
	ORDER BY {order}
	`), shopID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func categoriesQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*models.Category, error) {
	client := GetClient()
	query := queryBasis.withColumns(`
	  c.category_id,
	  c.category_name,
	  c.category_logo_small,
	  c.category_logo_large,
	  c.shops_count
	`)
	var rows []*struct {
		CategoryID        int
		CategoryName      string
		CategoryLogoSmall string
		CategoryLogoLarge string
		ShopsCount        int
	}
	_, err := client.Query(&rows, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, queryName)
	}
	categories := make([]*models.Category, len(rows))
	for i, row := range rows {
		categories[i] = &models.Category{
			ID:         row.CategoryID,
			Name:       row.CategoryName,
			Logo:       models.Logo{Small: row.CategoryLogoSmall, Large: row.CategoryLogoLarge},
			ShopsCount: row.ShopsCount,
		}
	}
	return categories, nil
}
