package repo

import (
	"context"

	"go-demo-gin/models"
	"go-demo-gin/pkg"
	"go-demo-gin/utils"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, tx *gorm.DB, u *models.User) error
	FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error)
	Update(ctx context.Context, tx *gorm.DB, u *models.User) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	List(ctx context.Context, tx *gorm.DB, pag *pkg.Pagination, search string) ([]models.User, int64, error)
	FindByUsername(ctx context.Context, tx *gorm.DB, username string) (*models.User, error)
}

type GormUserRepo struct{ db *gorm.DB }

func NewGormUserRepo(db *gorm.DB) *GormUserRepo { return &GormUserRepo{db: db} }

func (r *GormUserRepo) use(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *GormUserRepo) Create(ctx context.Context, tx *gorm.DB, u *models.User) error {
	return r.use(tx).WithContext(ctx).Create(u).Error
}

func (r *GormUserRepo) FindByID(ctx context.Context, tx *gorm.DB, id uint) (*models.User, error) {
	var u models.User
	if err := r.use(tx).WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) Update(ctx context.Context, tx *gorm.DB, u *models.User) error {
	return r.use(tx).WithContext(ctx).Updates(u).Error
}

func (r *GormUserRepo) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return r.use(tx).WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *GormUserRepo) List(ctx context.Context, tx *gorm.DB, pag *pkg.Pagination, search string) ([]models.User, int64, error) {
	q := r.use(tx).WithContext(ctx).Model(&models.User{})
	if search != "" {
		q = q.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%") // Postgres: ILIKE
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []models.User
	if err := q.Scopes(utils.Paginate(pag, q)).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepo) FindByUsername(ctx context.Context, tx *gorm.DB, username string) (*models.User, error) {
	var u models.User
	if err := r.use(tx).WithContext(ctx).
		Where("username = ?", username).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
