package handler

import (
	"net/http"

	"github.com/esmailemami/chess/shared/consts"
)

type Response interface {
	Status() int
	Message() string
	Result() any
}

type JSONResponse[T any] struct {
	Msg        string `json:"message"`
	StatusCode int    `json:"status"`
	Data       *T     `json:"data"`
}

func (r *JSONResponse[T]) Status() int {
	return r.StatusCode
}

func (r *JSONResponse[T]) Message() string {
	return r.Msg
}

func (r *JSONResponse[T]) Result() any {
	return r.Data
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

func OK[T any](data *T, msg ...string) Response {
	r := &JSONResponse[T]{
		Msg:        consts.OperationDone,
		StatusCode: http.StatusOK,
		Data:       data,
	}

	if len(msg) > 0 {
		r.Msg = msg[0]
	}

	return r
}

func OKBool(msg ...string) Response {
	ok := true
	return OK[bool](&ok)
}

func ListOK[T any](page, limit int, total int64, data []T, msg ...string) Response {
	r := &JSONResponse[ListResponse[T]]{
		Msg:        consts.OperationDone,
		StatusCode: http.StatusOK,
		Data:       NewListResponse[T](page, limit, total, data),
	}

	if len(msg) > 0 {
		r.Msg = msg[0]
	}

	return r
}
