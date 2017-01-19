package models

import (
	"testing"
)

func Benchmark3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		results := []int{2, 34, 2, 5, 54, 46, 64, 4, 6}
		var limit *uint = nil
		var tmp uint = 2
		offset := &tmp
		queryName := "asdfsdfsdfsd"
		arg := 228
		_ = computeTotalCount3(len(results), limit, offset, queryName, `
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1)
		`, arg)
	}
}
func Benchmark2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		results := []int{2, 34, 2, 5, 54, 46, 64, 4, 6}
		var limit *uint = nil
		var tmp uint = 2
		offset := &tmp
		queryName := "asdfsdfsdfsd"
		arg := 228
		f := lazyCountQuery(queryName, `
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1)
		`, arg)
		_ = computeTotalCount2(len(results), limit, offset, f)
	}
}
func Benchmark1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		results := []int{2, 34, 2, 5, 54, 46, 64, 4, 6}
		var limit *uint = nil
		var tmp uint = 2
		offset := &tmp
		queryName := "asdfsdfsdfsd"
		arg := 228
		f := lazyCountQuery(queryName, `
		SELECT count(DISTINCT m.id) total_count
		FROM mall m
		JOIN mall_shop ms ON m.id = ms.mall_id
		WHERE ms.shop_id = ANY ($1)
		`, arg)
		_ = computeTotalCount(results, limit, offset, f)
	}
}
