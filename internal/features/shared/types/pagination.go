package types

import (
	"errors"
)

type Pagination struct {
	Limit PositiveInt[int]
	Page  PositiveInt[int]
}

func (p *Pagination) Offset() int {
	return (p.Limit.Value() * p.Page.Value()) - p.Limit.Value()
}

func NewPagination(limit int, page int) (Pagination, error) {
	parsedLimit, err := NewPositiveInt(limit)
	if err != nil {
		return Pagination{}, errors.New("limit must be positive")
	}

	parsedPage, err := NewPositiveInt(page)
	if err != nil {
		return Pagination{}, errors.New("page must be positive")
	}

	return Pagination{
		Limit: parsedLimit,
		Page:  parsedPage,
	}, nil
}
