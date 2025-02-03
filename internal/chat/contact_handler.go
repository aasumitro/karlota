package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/internal/utils/http/middleware"
	"karlota.aasumitro.id/internal/utils/http/wrapper"
)

type contactHandler struct {
	service IContactService
}

func (h *contactHandler) fetch(ctx *gin.Context) {
	email := ctx.MustGet("email").(string)
	data, err := h.service.List(ctx, email)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func NewContactHandler(
	router gin.IRoutes,
	service IContactService,
) {
	handler := &contactHandler{service: service}
	router.GET("/contacts", middleware.Auth(), handler.fetch)
}
