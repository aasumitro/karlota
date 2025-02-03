package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"karlota.aasumitro.id/internal/model/request"
	"karlota.aasumitro.id/internal/utils/http/wrapper"
)

type authHandler struct {
	service IAuthService
}

// Login godoc
// @Schemes
// @Summary     Log User In
// @Description Login with email and password
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Param 		form    body        request.LoginRequest true "form request for login"
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/login [POST]
func (handler authHandler) Login(ctx *gin.Context) {
	var body request.LoginRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Call the service
	data, err := handler.service.SignIn(ctx, body)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(
		ctx, http.StatusCreated, data)
}

// ForgotPassword godoc
// @Schemes
// @Summary     Forgot Password
// @Description Request tog get reset password link
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Param 		form    body        request.ForgotPasswordRequest true "form request to get forgot password link"
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/forgot-password [POST]
func (handler authHandler) ForgotPassword(ctx *gin.Context) {
	var body request.ForgotPasswordRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Call the service
	data, err := handler.service.RequestForgotLink(ctx, body)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(
		ctx, http.StatusCreated, data)
}

// ResetPassword godoc
// @Schemes
// @Summary     Set New Password
// @Description Set User new password after get reset password link
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Param 		form    body        request.ResetPasswordRequest true "form request to set new user password"
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/reset-password [POST]
func (handler authHandler) ResetPassword(ctx *gin.Context) {
	var body request.ResetPasswordRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Call the service
	if err := handler.service.SetNewPassword(ctx, body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err)
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(
		ctx, http.StatusCreated, "password change success")
}

// Register godoc
// @Schemes
// @Summary     Create new Account
// @Description Create new user account and get access to the app
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Param 		form    body        request.RegisterRequest true "form request to create an account"
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/register [POST]
func (handler authHandler) Register(ctx *gin.Context) {
	var body request.RegisterRequest
	// Bind the request.Body to the struct
	if err := ctx.ShouldBind(&body); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate the request.Body
	if err := body.Validate(); err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnprocessableEntity, err)
		return
	}
	// Call the service
	data, err := handler.service.SignUp(ctx, body)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(
		ctx, http.StatusCreated, data)
}

// RefreshToken godoc
// @Schemes
// @Summary     Refresh User Session
// @Description Do refresh token and get new access token
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Param       X-REFRESH-TOKEN header string true "Refresh Token"
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/refresh-token [POST]
func (handler authHandler) RefreshToken(ctx *gin.Context) {
	// get refresh token from header
	refreshToken := ctx.GetHeader("X-REFRESH-TOKEN")
	if refreshToken == "" {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnauthorized, "refresh token is required")
		return
	}
	// validate and generate new token
	data, err := handler.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusUnauthorized, err.Error())
		return
	}
	// Return the response
	wrapper.NewHTTPRespondWrapper(
		ctx, http.StatusCreated, data)
}

// VerifyEmail godoc
// @Schemes
// @Summary     Verify User Email
// @Description Verify user email via given link
// @Tags 	    Auth
// @Accept 	    json
// @Produce     json
// @Success     201     {object}    wrapper.EmptyRespond{} "CREATED RESPOND"
// @Failure     400     {object}    wrapper.EmptyRespond{} "BAD REQUEST RESPOND"
// @Failure     422     {object}    wrapper.EmptyRespond{} "UNPROCESSABLE ENTITY REQUEST RESPOND"
// @Failure     500     {object}    wrapper.EmptyRespond{} "INTERNAL SERVER ERROR RESPOND"
// @Router      /api/v1/auth/verify-email [GET]
func (handler authHandler) VerifyEmail(ctx *gin.Context) {
	_ = handler.service.VerifyEmail(ctx, ctx.Query("exchange_token"))
	ctx.Redirect(http.StatusTemporaryRedirect, "/fe")
}

func NewAuthHandler(
	router *gin.RouterGroup,
	service IAuthService,
) {
	handler := &authHandler{service: service}
	router.POST("/login", handler.Login)
	router.POST("/forgot-password", handler.ForgotPassword)
	router.POST("/reset-password", handler.ResetPassword)
	router.POST("/register", handler.Register)
	router.POST("/refresh-token", handler.RefreshToken)
	router.GET("/verify-email", handler.VerifyEmail)
}
