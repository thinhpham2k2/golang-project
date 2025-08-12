package services

import (
	"errors"
	"go-demo-gin/models"
	"go-demo-gin/pkg"
	"go-demo-gin/repo"
	userRequest "go-demo-gin/requests/user"
	userResponse "go-demo-gin/responses/user"
	"go-demo-gin/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	db       *gorm.DB
	userRepo repo.UserRepo
}

func NewUserService(db *gorm.DB) *UserService { // "constructor"
	return &UserService{db: db, userRepo: repo.NewGormUserRepo(db)}
}

func (s *UserService) CreateUser(c *gin.Context, in *userRequest.UserCreate) (*userResponse.UserDetail, int, string) {
	// Logging
	utils.Log(c, logrus.InfoLevel, "Entering the create user service")

	// Lấy localizer cho i18n
	localizer := utils.LoadVariablesInContext(c)

	// Mapper
	var user models.User
	copier.Copy(&user, &in)

	// Transaction boundary
	if err := s.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// 1) Tạo user
		if err := s.userRepo.Create(c.Request.Context(), tx, &user); err != nil {
			return err // => auto ROLLBACK
		}
		// 2) (Ví dụ) gán role mặc định/ghi audit... (nếu thêm bước, vẫn trong tx)
		return nil // => COMMIT
	}); err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.CREATE_FAIL, nil)
	}

	// Mapper
	var detail userResponse.UserDetail
	copier.Copy(&detail, &user)

	return &detail, http.StatusCreated, ""
}

func (s *UserService) GetUserList(c *gin.Context, pag *pkg.Pagination, search string) (*pkg.Pagination, int, string) {
	// Logging
	utils.Log(c, logrus.InfoLevel, "Entering the get list of users service")

	// Query
	users, total, err := s.userRepo.List(c.Request.Context(), nil, pag, search)
	if err != nil {
		return nil, http.StatusBadRequest, err.Error()
	}

	// Mapper
	var list []userResponse.UserList
	copier.Copy(&list, &users)

	// Assign results to Pagination Struct
	pag.TotalRows = total
	pag.Result = list

	return pag, http.StatusOK, ""
}

func (s *UserService) GetUserById(c *gin.Context, idStr string) (*userResponse.UserDetail, int, string) {
	utils.Log(c, logrus.InfoLevel, "Entering the get user by id service")
	loc := utils.LoadVariablesInContext(c)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(loc, utils.INVALID_VALUE, nil)
	}

	u, err := s.userRepo.FindByID(c.Request.Context(), nil, uint(id))
	if err != nil {
		return nil, http.StatusNotFound, utils.LoadI18nMessage(loc, utils.NOT_FOUND, nil)
	}

	var detail userResponse.UserDetail
	_ = copier.Copy(&detail, u)
	return &detail, http.StatusOK, ""
}

func (s *UserService) UpdateUser(c *gin.Context, in *userRequest.UserUpdate, idStr string) (*userResponse.UserDetail, int, string) {
	utils.Log(c, logrus.InfoLevel, "Entering the update user service")
	localizer := utils.LoadVariablesInContext(c)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
	}

	var out *userResponse.UserDetail

	err = s.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// 1) Load hiện trạng
		u, err := s.userRepo.FindByID(c.Request.Context(), tx, uint(id))
		if err != nil {
			return err
		}

		copier.Copy(&u, in)

		// 2) Update bằng Updates(struct) (bỏ qua zero-value)
		if err := s.userRepo.Update(c.Request.Context(), tx, u); err != nil {
			return err
		}

		// 3) (khuyến nghị) reload để lấy DB-managed fields (UpdatedAt, v.v.)
		if err := tx.First(u, u.ID).Error; err != nil {
			return err
		}

		var d userResponse.UserDetail
		copier.Copy(&d, u)
		out = &d
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, utils.LoadI18nMessage(localizer, utils.NOT_FOUND, nil)
		}
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.UPDATE_FAIL, nil)
	}

	return out, http.StatusOK, ""
}

func (s *UserService) DeleteUser(c *gin.Context, idStr string) (int, string) {
	utils.Log(c, logrus.InfoLevel, "Entering the delete user service")
	localizer := utils.LoadVariablesInContext(c)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
	}

	// Transaction boundary
	if err := s.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// Nếu có quan hệ phụ: xoá trước trong cùng tx (ví dụ)
		// if err := s.userRoleRepo.DeleteByUserID(c.Request.Context(), tx, uint(id)); err != nil { return err }
		// if err := s.noteRepo.DeleteByOwner(c.Request.Context(), tx, uint(id)); err != nil { return err }

		// Xoá chính user
		return s.userRepo.Delete(c.Request.Context(), tx, uint(id))
	}); err != nil {
		// Phân loại lỗi: không tìm thấy vs lỗi khác
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, utils.LoadI18nMessage(localizer, utils.NOT_FOUND, nil)
		}
		return http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.DELETE_FAIL, nil)
	}

	return http.StatusNoContent, ""
}
