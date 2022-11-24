// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

var (
	DefaultPageSize    = 20
	MaxDefaultPageSize = 100
)

type Pagination struct {
	Page     int
	PageSize int
}

// LimitOffset returns LIMIT and OFFSET parameter for SQL.
// The first page is page 0.
func (p Pagination) LimitOffset() (limit, offset int) {
	page := p.Page
	pageSize := p.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	return pageSize, (page - 1) * pageSize
}

type Cursor struct {
	Value    interface{}
	PageSize int
}

func (p Cursor) Limit() int {
	pageSize := p.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxDefaultPageSize {
		pageSize = MaxDefaultPageSize
	}
	return pageSize
}
