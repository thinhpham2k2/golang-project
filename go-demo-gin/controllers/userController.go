package controllers

import (
	"go-demo-gin/pkg"
	userRequest "go-demo-gin/requests/user"
	errorResponse "go-demo-gin/responses/error"
	userResponse "go-demo-gin/responses/user"
	"go-demo-gin/services"
	"go-demo-gin/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	_ userResponse.UserList
	_ userResponse.UserDetail
	_ errorResponse.HTTPError
)

type UserController struct {
	v   *utils.Validator
	svc *services.UserService
}

func NewUserController(v *utils.Validator, svc *services.UserService) *UserController {
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
	utils.Log(c, logrus.InfoLevel, "Entering the create user controller")

	// Get data off request body
	var create userRequest.UserCreate
	if err := c.ShouldBindJSON(&create); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Request binding failed: "+err.Error())
		return
	}

	// Validation
	if err := create.Validate(c, h.v); err != nil {
		utils.HandleValidationError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Validation failed")
		return
	}

	// Create user
	detail, status, err := h.svc.CreateUser(c, &create)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Create user failed: "+err)
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
	utils.Log(c, logrus.InfoLevel, "Entering the get list of users controller")

	// Get parameters in query string
	search := c.Query("search")

	// Get pagination
	var pag pkg.Pagination
	if err := c.ShouldBindQuery(&pag); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Request binding failed: "+err.Error())
		return
	}

	// Get user list
	result, status, err := h.svc.GetUserList(c, &pag, search)
	if result == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Get list of users failed: "+err)
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
	utils.Log(c, logrus.InfoLevel, "Entering the get user by id controller")

	// Get id from url
	id := c.Param("id")

	// Gte user detail
	detail, status, err := h.svc.GetUserById(c, id)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Get user by id failed: "+err)
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
	utils.Log(c, logrus.InfoLevel, "Entering the update user controller")

	// Get id from url
	id := c.Param("id")

	// Get data off request body
	var update userRequest.UserUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		utils.HandleBindError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Request binding failed: "+err.Error())
		return
	}

	// Validation
	if err := update.Validate(c, h.v); err != nil {
		utils.HandleValidationError(c, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Validation failed")
		return
	}

	// Update user
	detail, status, err := h.svc.UpdateUser(c, &update, id)
	if detail == nil || err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Update user failed: "+err)
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
	utils.Log(c, logrus.InfoLevel, "Entering the delete user controller")

	// Get id from url
	id := c.Param("id")

	// Delete user
	status, err := h.svc.DeleteUser(c, id)
	if err != "" {
		utils.HandleServiceError(c, status, err)
		// Logging
		utils.Log(c, logrus.ErrorLevel, "Delete user failed: "+err)
		return
	}

	c.Status(status)
}
