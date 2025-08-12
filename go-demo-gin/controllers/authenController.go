package controllers

import (
	authenRequest "go-demo-gin/requests/authen"
	errorResponse "go-demo-gin/responses/error"
	"go-demo-gin/services"
	"go-demo-gin/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	_ errorResponse.HTTPError
)

type TokenResponse struct {
	Token *string `json:"token"`
}

type AuthenController struct {
	svc *services.AuthenService
}

func NewAuthenController(svc *services.AuthenService) *AuthenController {
	return &AuthenController{svc: svc}
}

// Login login to system
//
// @Summary      Login
// @Description  Login to system
// @Tags         🔐Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      authenRequest.LoginForm  true  "Login form"
// @Success      200      {object}  TokenResponse
// @Failure      400      {object}  errorResponse.HTTPError
// @Failure      500      {string}  httputil.HTTPError
// @Router       /api/v1/authen/login [post]
func (h *AuthenController) Login(c *gin.Context) {
	// Logging
	utils.Log(c, logrus.InfoLevel, "Entering the login controller")

	// Get the username/password off req body
	var authen authenRequest.LoginForm
	if err := c.ShouldBind(&authen); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Request binding failed: "+err.Error())
		return
	}

	// Check user infor & generate jwt token
	token, status, err := h.svc.Authenticate(c, &authen)
	if token == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Authentication failed: "+err)
		return
	}

	// Logging
	utils.Log(c, logrus.ErrorLevel, "Authentication successful for user: "+authen.Username)
	// Send it back
	c.JSON(status, TokenResponse{
		Token: token,
	})
}
