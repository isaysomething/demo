package pagination

import (
	"math"
	"net/url"
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

func (p *Pagination) PageCount() uint64 {
	return uint64(math.Ceil(float64(p.Total) / float64(p.Limit)))
}

func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

func (p *Pagination) PrevURL() *url.URL {
	u := &url.URL{}
	query := u.Query()
	query.Add("page", strconv.FormatUint(p.Page-1, 10))
	query.Add("limit", strconv.FormatUint(p.Limit, 10))
	u.RawQuery = query.Encode()
	return u
}

func (p *Pagination) HasNext() bool {
	return p.Page < p.PageCount()
}

func (p *Pagination) NextURL() *url.URL {
	u := &url.URL{}
	query := u.Query()
	query.Add("page", strconv.FormatUint(p.Page+1, 10))
	query.Add("limit", strconv.FormatUint(p.Limit, 10))
	u.RawQuery = query.Encode()
	return u
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
