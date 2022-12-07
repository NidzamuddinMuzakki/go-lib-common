package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// Response using jsend format
// ref: https://github.com/omniti-labs/jsend
type Response struct {
	Status       any    `json:"status"`
	Message      string `json:"message"`
	Data         any    `json:"data,omitempty"`
	Limit        uint   `json:"limit,omitempty"`
	TotalRecords uint64 `json:"totalRecords,omitempty"`
	CurrentPage  uint   `json:"currentPage,omitempty"`
	NextPage     uint   `json:"nextPage,omitempty"`
	PreviousPage uint   `json:"previousPage,omitempty"`
	TotalPages   uint   `json:"totalPages,omitempty"`
}

// RouteNotFound handle when user is hitting non-exist endpoint.
// It will imediately return error 404 not found.
func RouteNotFound(e *gin.Engine) {
	e.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, Response{
			Message: http.StatusText(http.StatusNotFound),
			Status:  StatusFail,
		})
	})
}
