package models

import (
	"github.com/pkg/errors"
	pg "gopkg.in/pg.v5"
	"mallfin_api/db"
	"mallfin_api/utils"
)

func GetCategoryDetails(categoryID int, cityID *int) (*Category, error) {
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

func GetCategories(cityID *int, sortKey *string) ([]*Category, error) {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
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

func GetCategoriesByIDs(categoryIDs []int, cityID *int) ([]*Category, int, error) {
	categoryIDsArray := pg.Array(categoryIDs)
	queryName := utils.CurrentFuncName()
	categories, err := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ANY (?0)
	`), categoryIDsArray)
	if err != nil {
		return nil, 0, err
	}
	totalCount := len(categories)
	return categories, totalCount, nil
}

func GetCategoriesByShop(shopID int, cityID *int, sortKey *string) ([]*Category, error) {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
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

func categoriesQuery(queryName string, queryBasis baseQuery, args ...interface{}) ([]*Category, error) {
	client := db.GetClient()
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
	categories := make([]*Category, len(rows))
	for i, row := range rows {
		categories[i] = &Category{
			ID:         row.CategoryID,
			Name:       row.CategoryName,
			Logo:       Logo{Small: row.CategoryLogoSmall, Large: row.CategoryLogoLarge},
			ShopsCount: row.ShopsCount,
		}
	}
	return categories, nil
}
