package pagination

import (
	"strconv"

	"github.com/clevergo/clevergo"
)

var (
	PageParam  = "page"
	LimitParam = "limit"
	MaxLimit   = 1000
)

type Pagination struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

func New(page, limit int) *Pagination {
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

func NewFromContext(ctx *clevergo.Context) *Pagination {
	return New(
		parsePage(ctx, 1),
		parseLimit(ctx, 20),
	)
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func parsePage(ctx *clevergo.Context, defaultValue int) int {
	v, err := strconv.Atoi(ctx.QueryParam(PageParam))
	if err == nil && v > 0 {
		return v
	}

	return defaultValue
}

func parseLimit(ctx *clevergo.Context, defaultValue int) int {
	v, err := strconv.Atoi(ctx.QueryParam(LimitParam))
	if err == nil && v > 0 {
		if v > MaxLimit {
			v = MaxLimit
		}
		return v
	}

	return defaultValue
}
