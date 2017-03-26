package models

import (
	pg "gopkg.in/pg.v5"
	"mallfin_api/db"
	"mallfin_api/utils"
)

func GetCategoryDetails(categoryID int, cityID *int) *Category {
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ?0
	LIMIT 1
	`), categoryID)
	if len(categories) == 0 {
		return nil
	}
	return categories[0]
}

func GetCategories(cityID *int, sortKey *string) []*Category {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	ORDER BY {order}
	`))
	return categories
}

func GetCategoriesByIDs(categoryIDs []int, cityID *int) ([]*Category, int) {
	categoryIDsArray := pg.Array(categoryIDs)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, baseQuery(`
	SELECT {columns}
	FROM category c
	WHERE c.category_id = ANY (?0)
	`), categoryIDsArray)
	totalCount := len(categories)
	return categories, totalCount
}

func GetCategoriesByShop(shopID int, cityID *int, sortKey *string) []*Category {
	orderBy := CATEGORIES_SORT_KEYS.CorrespondingOrderBy(sortKey)
	queryName := utils.CurrentFuncName()
	categories := categoriesQuery(queryName, orderBy.CompileBaseQuery(`
	SELECT {columns}
	FROM category c
	  JOIN shop_category sc ON c.category_id = sc.category_id
	WHERE sc.shop_id = ?0
	ORDER BY {order}
	`), shopID)
	return categories
}

func categoriesQuery(queryName string, queryBasis baseQuery, args ...interface{}) []*Category {
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
	locLog := moduleLog.WithField("query", queryName)
	if err != nil {
		locLog.Panicf("Cannot get categories rows: %s", err)
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
	return categories
}
