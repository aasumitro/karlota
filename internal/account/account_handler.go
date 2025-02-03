package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/request"
	"karlota.aasumitro.id/internal/utils/http/wrapper"
)

type accountHandler struct {
	service IAccountService
}

// Show godoc
// @Schemes
// @Summary     User Profile
// @Description Get User Profile
// @Tags 	    Account
// @Accept 	    json
// @Produce     json
// @Security 	ApiKeyAuth
// @Success     200     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/account [GET]
func (handler accountHandler) Show(ctx *gin.Context) {
	rawID := ctx.MustGet("id").(float64)
	data, err := handler.service.Profile(ctx, uint(rawID))
	if err != nil {
		wrapper.NewHTTPRespondWrapper(ctx, http.StatusBadRequest, err.Error())
		return
	}
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

// EditDisplayName godoc
// @Schemes
// @Summary     User Update Name
// @Description Update user display name
// @Tags 	    Account
// @Accept 	    json
// @Produce     json
// @Security 	ApiKeyAuth
// @Param 		form    body        request.UpdateAccountRequest true "form request for login"
// @Success     200     {object}    wrapper.EmptyRespond{} "OK RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/account [PATCH]
func (handler accountHandler) EditDisplayName(ctx *gin.Context) {
	rawID := ctx.MustGet("id").(float64)
	var body request.UpdateAccountRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Update the display name
	data, err := handler.service.UpdateDisplayName(ctx,
		uint(rawID), body.DisplayName)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

// EditPassword godoc
// @Schemes
// @Summary     User Update Password
// @Description Update user password
// @Tags 	    Account
// @Accept 	    json
// @Produce     json
// @Security 	ApiKeyAuth
// @Param 		form    body        request.UpdatePasswordRequest true "form request for login"
// @Success     200     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/account/password [PATCH]
func (handler accountHandler) EditPassword(ctx *gin.Context) {
	rawID := ctx.MustGet("id").(float64)
	var body request.UpdatePasswordRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Update the display name
	data, err := handler.service.UpdatePassword(ctx,
		uint(rawID), body.OldPassword, body.NewPassword)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(ctx, http.StatusOK, data)
}

func NewAccountHandler(
	router gin.IRoutes,
	service IAccountService,
) {
	handler := &accountHandler{service: service}
	router.GET(common.EmptyPath, handler.Show)
	router.PATCH(common.EmptyPath, handler.EditDisplayName)
	router.PATCH("/password", handler.EditPassword)
}
