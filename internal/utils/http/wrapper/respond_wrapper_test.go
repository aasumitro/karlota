package wrapper_test

import (
	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/internal/utils/http/wrapper"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHttpRespond(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		data     interface{}
		args     []any
		expected interface{}
	}{
		{
			name:     "success with no pagination",
			code:     http.StatusOK,
			data:     []string{"foo", "bar"},
			args:     nil,
			expected: wrapper.CommonRespond{Data: []string{"foo", "bar"}},
		},
		{
			name: "success with pagination",
			code: http.StatusOK,
			data: []string{"foo", "bar"},
			args: []any{
				2,
				1,
				wrapper.Paging{URL: "https://example.com/next", Path: "/next"},
				wrapper.Paging{URL: "https://example.com/prev", Path: "/prev"},
				100,
			},
			expected: wrapper.SuccessWithPaginationRespond{
				Total:    2,
				Current:  1,
				Next:     wrapper.Paging{URL: "https://example.com/next", Path: "/next"},
				Previous: wrapper.Paging{URL: "https://example.com/prev", Path: "/prev"},
				Data:     []string{"foo", "bar"},
			},
		},
		{
			name:     "error with data",
			code:     http.StatusBadRequest,
			data:     "invalid request",
			expected: wrapper.ErrorRespond{Errors: "invalid request"},
		},
		{
			name:     "error with no data",
			code:     http.StatusBadRequest,
			expected: wrapper.ErrorRespond{Errors: "something went wrong with the request"},
		},
		{
			name:     "error with no data and server error code",
			code:     http.StatusInternalServerError,
			expected: wrapper.ErrorRespond{Errors: "something went wrong with the server"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(writer)
			wrapper.NewHTTPRespondWrapper(c, test.code, test.data, test.args...)
		})
	}
}
