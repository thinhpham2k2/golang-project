package controllers

import (
	"context"
	"go-demo-gin/pkg"
	userRequest "go-demo-gin/requests/user"
	errorResponse "go-demo-gin/responses/error"
	userResponse "go-demo-gin/responses/user"
	"go-demo-gin/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	_ userResponse.UserList
	_ userResponse.UserDetail
	_ errorResponse.HTTPError
)

type UserService interface {
	CreateUser(ctx context.Context, in *userRequest.UserCreate) (*userResponse.UserDetail, int, string)
	GetUserList(ctx context.Context, pag *pkg.Pagination, search string) (*pkg.Pagination, int, string)
	GetUserById(ctx context.Context, id string) (*userResponse.UserDetail, int, string)
	UpdateUser(ctx context.Context, in *userRequest.UserUpdate, id string) (*userResponse.UserDetail, int, string)
	DeleteUser(ctx context.Context, id string) (int, string)
}

type UserController struct {
	v   *utils.Validator
	svc UserService
}

func NewUserController(v *utils.Validator, svc UserService) *UserController {
	return &UserController{v: v, svc: svc}
}

// UsersCreate creates a new user
//
// @Summary      Create user
// @Description  Create a new user
// @Tags         👨🏻‍💼Users
// @Security	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      userRequest.UserCreate  true  "User to create"
// @Success      201      {object}  userResponse.UserDetail
// @Failure      400      {object}  errorResponse.HTTPError
// @Failure      500      {string}  httputil.HTTPError
// @Router       /api/v1/users [post]
func (h *UserController) UsersCreate(c *gin.Context) {
	// Logging
	ctx := c.Request.Context()
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the create user controller", nil)

	// Get data off request body
	var create userRequest.UserCreate
	if err := c.ShouldBindJSON(&create); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Request binding failed: "+err.Error(), nil)
		return
	}

	// Validation
	if err := create.Validate(ctx, h.v); err != nil {
		utils.HandleValidationError(c, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Validation failed", nil)
		return
	}

	// Create user
	detail, status, err := h.svc.CreateUser(ctx, &create)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Create user failed: "+err, nil)
		return
	}

	c.JSON(status, detail)
}

// UsersIndex lists all existing users
//
// @Summary      List users
// @Description  Get list of all users
// @Tags         👨🏻‍💼Users
// @Security	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        search		query     string  false  "Search query"
// @Param        limit		query     string  false  "Number of results per page"				default(10)
// @Param        page		query     string  false  "Current page in the paginated results"	default(1)
// @Param        sort		query     string  false  "Sorting criteria for the results"			default(id desc)
// @Success      200   {array}   pkg.Pagination{result=[]userResponse.UserList}
// @Failure      400   {object}  errorResponse.HTTPError
// @Failure      500   {string}  httputil.HTTPError
// @Router       /api/v1/users [get]
func (h *UserController) UsersIndex(c *gin.Context) {
	// Logging
	ctx := c.Request.Context()
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the get list of users controller", nil)

	// Get parameters in query string
	search := c.Query("search")

	// Get pagination
	var pag pkg.Pagination
	if err := c.ShouldBindQuery(&pag); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Request binding failed: "+err.Error(), nil)
		return
	}

	// Get user list
	result, status, err := h.svc.GetUserList(ctx, &pag, search)
	if result == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Get list of users failed: "+err, nil)
		return
	}

	c.JSON(status, result)
}

// UsersShow get user detail
//
// @Summary      Get user detail
// @Description  Get user by ID
// @Tags         👨🏻‍💼Users
// @Security	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  userResponse.UserDetail
// @Failure      404  {string}  httputil.HTTPError
// @Failure      500  {string}  httputil.HTTPError
// @Router       /api/v1/users/{id} [get]
func (h *UserController) UsersShow(c *gin.Context) {
	// Logging
	ctx := c.Request.Context()
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the get user by id controller", nil)

	// Get id from url
	id := c.Param("id")

	// Gte user detail
	detail, status, err := h.svc.GetUserById(ctx, id)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Get user by id failed: "+err, nil)
		return
	}

	c.JSON(status, detail)
}

// UsersUpdate updates an existing user
//
// @Summary      Update user
// @Description  Update existing user by ID
// @Tags         👨🏻‍💼Users
// @Security	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "user ID"
// @Param        request  body      userRequest.UserUpdate  true  "Updated user data"
// @Success      200      {object}  userResponse.UserDetail
// @Failure      400      {object}  errorResponse.HTTPError
// @Failure      404      {string}  httputil.HTTPError
// @Failure      500      {string}  httputil.HTTPError
// @Router       /api/v1/users/{id} [put]
func (h *UserController) UsersUpdate(c *gin.Context) {
	// Logging
	ctx := c.Request.Context()
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the update user controller", nil)

	// Get id from url
	id := c.Param("id")

	// Get data off request body
	var update userRequest.UserUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Request binding failed: "+err.Error(), nil)
		return
	}

	// Validation
	if err := update.Validate(ctx, h.v); err != nil {
		utils.HandleValidationError(c, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Validation failed", nil)
		return
	}

	// Update user
	detail, status, err := h.svc.UpdateUser(ctx, &update, id)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Update user failed: "+err, nil)
		return
	}

	c.JSON(status, detail)
}

// UsersDelete deletes an user
//
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         👨🏻‍💼Users
// @Security	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      404  {string}  httputil.HTTPError
// @Failure      500  {string}  httputil.HTTPError
// @Router       /api/v1/users/{id} [delete]
func (h *UserController) UsersDelete(c *gin.Context) {
	// Logging
	ctx := c.Request.Context()
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the delete user controller", nil)

	// Get id from url
	id := c.Param("id")

	// Delete user
	status, err := h.svc.DeleteUser(ctx, id)
	if err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.LogCtx(ctx, logrus.ErrorLevel, "Delete user failed: "+err, nil)
		return
	}

	c.Status(status)
}
