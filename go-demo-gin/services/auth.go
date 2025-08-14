package services

import (
	"context"
	"errors"
	authenRequest "go-demo-gin/requests/authen"
	"go-demo-gin/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthConfig struct {
	JWTKey    []byte
	Issuer    string
	AccessTTL time.Duration
}

type AuthService struct {
	db       *gorm.DB
	cfg      AuthConfig
	userRepo UserRepository
}

func NewAuthService(db *gorm.DB, cfg AuthConfig, ur UserRepository) *AuthService {
	return &AuthService{
		db:       db,
		cfg:      cfg,
		userRepo: ur,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, in *authenRequest.LoginForm) (*string, int, string) {
	// Logging
	utils.LogCtx(ctx, logrus.InfoLevel, "Entering the login service", nil)

	// Lấy localizer cho i18n
	localizer := utils.LocalizerFrom(ctx)

	// Look up requested user
	// lấy user qua repo (context-aware)
	ctxTx := utils.WithTx(ctx, nil)
	user, err := s.userRepo.FindByUsername(ctxTx, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusUnauthorized, utils.LoadI18nMessage(localizer, utils.INVALID_USERNAME_PASSWORD, nil)
		}
		utils.LogCtx(ctx, logrus.InfoLevel, "DB error on login", nil)
		return nil, http.StatusInternalServerError, utils.LoadI18nMessage(localizer, utils.INTERNAL_ERROR, nil)
	}

	// Compare sent in pass with saved user pass hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.INVALID_USERNAME_PASSWORD, nil)
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"id":  user.ID,
		"exp": time.Now().Add(s.cfg.AccessTTL).Unix(),
		"iss": s.cfg.Issuer,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, http.StatusBadRequest, utils.LoadI18nMessage(localizer, utils.FAIL_CREATE_TOKEN, nil)
	}

	return &tokenString, http.StatusOK, ""
}
