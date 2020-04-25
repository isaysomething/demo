package pagination

import (
	"strconv"

	"github.com/clevergo/clevergo"
)

var (
	PageParam         = "page"
	LimitParam        = "limit"
	MaxLimit   uint64 = 1000
)

type Pagination struct {
	Page  uint64      `json:"page"`
	Limit uint64      `json:"limit"`
	Total uint64      `json:"total"`
	Items interface{} `json:"items"`
}

func New(page, limit uint64) *Pagination {
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

func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.Limit
}

func parsePage(ctx *clevergo.Context, defaultValue uint64) uint64 {
	v, err := strconv.ParseUint(ctx.QueryParam(PageParam), 10, 64)
	if err == nil && v > 0 {
		return v
	}

	return defaultValue
}

func parseLimit(ctx *clevergo.Context, defaultValue uint64) uint64 {
	v, err := strconv.ParseUint(ctx.QueryParam(LimitParam), 10, 64)
	if err == nil && v > 0 {
		if v > MaxLimit {
			v = MaxLimit
		}
		return v
	}

	return defaultValue
}
