package wrapper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmptyRespond struct{}

type CommonRespond struct {
	Data   any `json:"data,omitempty"`
	Errors any `json:"errors,omitempty"`
}

type SuccessWithPaginationRespond struct {
	Count    int    `json:"total_data"`
	Total    int    `json:"total_page"`
	Current  int    `json:"current_page"`
	Next     Paging `json:"next"`
	Previous Paging `json:"previous"`
	Data     any    `json:"data,omitempty"`
	Errors   any    `json:"errors,omitempty"`
}

type Paging struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type ErrorRespond struct {
	Errors any `json:"errors,omitempty"`
}

func NewHTTPRespondWrapper(context *gin.Context, code int, data interface{}, args ...any) {
	if code == http.StatusOK || code == http.StatusCreated {
		if len(args) > 5 {
			context.JSON(code, SuccessWithPaginationRespond{
				Total:    args[0].(int),
				Current:  args[1].(int),
				Next:     args[2].(Paging),
				Previous: args[3].(Paging),
				Count:    args[4].(int),
				Data:     data,
			})

			return
		}

		context.JSON(code, CommonRespond{Data: data})

		return
	}

	if code == http.StatusUnprocessableEntity {
		context.JSON(code, CommonRespond{Errors: data})
		return
	}

	msg := func() string {
		switch {
		case data != nil:
			return data.(string)
		case code == http.StatusBadRequest:
			return "something went wrong with the request"
		default:
			return "something went wrong with the server"
		}
	}()
	context.JSON(code, ErrorRespond{Errors: msg})
}
