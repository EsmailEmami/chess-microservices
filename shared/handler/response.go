package handler

import (
	"net/http"

	"github.com/esmailemami/chess/shared/consts"
)

type Response[T any] struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    *T     `json:"data"`
}

type ListResponse[T any] struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	Limit    int64 `json:"limit"`
	LastPage int64 `json:"last_page"`
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	Data     []T   `json:"data"`
}

func NewListResponse[T any](page, limit int, total int64, data []T) *ListResponse[T] {
	response := new(ListResponse[T])
	response.Page = int64(page)
	response.Limit = int64(limit)
	response.From = ((response.Page - 1) * response.Limit) + 1
	response.To = response.From + response.Limit - 1
	response.Total = total
	response.Data = data

	// calculate last page
	lp := float64(total) / float64(limit)
	lastPage := int64(lp)
	if lp > float64(lastPage) {
		lastPage++
	}
	response.LastPage = lastPage

	return response

}

func OK[T any](data *T, msg ...string) *Response[T] {
	r := &Response[T]{
		Message: consts.OperationDone,
		Status:  http.StatusOK,
		Data:    data,
	}

	if len(msg) > 0 {
		r.Message = msg[0]
	}

	return r
}

func OKBool(msg ...string) *Response[bool] {
	ok := true
	return OK[bool](&ok)
}

func ListOK[T any](page, limit int, total int64, data []T, msg ...string) *Response[ListResponse[T]] {
	r := &Response[ListResponse[T]]{
		Message: consts.OperationDone,
		Status:  http.StatusOK,
		Data:    NewListResponse[T](page, limit, total, data),
	}

	if len(msg) > 0 {
		r.Message = msg[0]
	}

	return r
}
