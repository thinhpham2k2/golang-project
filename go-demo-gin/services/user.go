package services

import (
	"context"
	"errors"
	"go-demo-gin/models"
	"go-demo-gin/pkg"
	userRequest "go-demo-gin/requests/user"
	userResponse "go-demo-gin/responses/user"
	"go-demo-gin/utils"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, pag *pkg.Pagination, search string) ([]models.User, int64, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserService struct {
	db       *gorm.DB
	userRepo UserRepository
}

func NewUserService(db *gorm.DB, ur UserRepository) *UserService { // "constructor"
	return &UserService{db: db, userRepo: ur}
}

func (s *UserService) CreateUser(ctx context.Context, in *userRequest.UserCreate) (*userResponse.UserDetail, int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the create user service", nil)

	// Lấy localizer cho i18n
	localizer := utils.LocalizerFrom(ctx)

	// Mapper
	var user models.User
	copier.Copy(&user, in)

	// Transaction boundary
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Tạo user
		ctxTx := utils.WithTx(ctx, tx)
		if err := s.userRepo.Create(ctxTx, &user); err != nil {
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

func (s *UserService) GetUserList(ctx context.Context, pag *pkg.Pagination, search string) (*pkg.Pagination, int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the get list of users service", nil)
	// Query
	users, total, err := s.userRepo.List(ctx, pag, search)
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

func (s *UserService) GetUserById(ctx context.Context, idStr string) (*userResponse.UserDetail, int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the get user by id service", nil)

	// Lấy localizer cho i18n
	localizer := utils.LocalizerFrom(ctx)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
	}

	u, err := s.userRepo.FindByID(ctx, uint(id))
	if err != nil {
		return nil, http.StatusNotFound, utils.LoadI18nMessage(localizer, utils.NOT_FOUND, nil)
	}

	var detail userResponse.UserDetail
	_ = copier.Copy(&detail, u)
	return &detail, http.StatusOK, ""
}

func (s *UserService) UpdateUser(ctx context.Context, in *userRequest.UserUpdate, idStr string) (*userResponse.UserDetail, int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the update user service", nil)

	// Lấy localizer cho i18n
	localizer := utils.LocalizerFrom(ctx)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
	}

	var out *userResponse.UserDetail

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Load hiện trạng
		ctxTx := utils.WithTx(ctx, tx)
		u, err := s.userRepo.FindByID(ctxTx, uint(id))
		if err != nil {
			return err
		}

		copier.Copy(&u, in)

		// 2) Update bằng Updates(struct) (bỏ qua zero-value)
		if err := s.userRepo.Update(ctxTx, u); err != nil {
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

func (s *UserService) DeleteUser(ctx context.Context, idStr string) (int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the delete user service", nil)

	// Lấy localizer cho i18n
	localizer := utils.LocalizerFrom(ctx)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_VALUE, nil)
	}

	// Transaction boundary
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Nếu có quan hệ phụ: xoá trước trong cùng tx (ví dụ)
		// if err := s.userRoleRepo.DeleteByUserID(c.Request.Context(), tx, uint(id)); err != nil { return err }
		// if err := s.noteRepo.DeleteByOwner(c.Request.Context(), tx, uint(id)); err != nil { return err }

		// Xoá chính user
		ctxTx := utils.WithTx(ctx, tx)
		return s.userRepo.Delete(ctxTx, uint(id))
	}); err != nil {
		// Phân loại lỗi: không tìm thấy vs lỗi khác
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, utils.LoadI18nMessage(localizer, utils.NOT_FOUND, nil)
		}
		return http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.DELETE_FAIL, nil)
	}

	return http.StatusNoContent, ""
}
