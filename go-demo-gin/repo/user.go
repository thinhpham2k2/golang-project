package repo

import (
	"context"

	"go-demo-gin/models"
	"go-demo-gin/pkg"
	"go-demo-gin/utils"

	"gorm.io/gorm"
)

type GormUserRepo struct{ db *gorm.DB }

func NewGormUserRepo(db *gorm.DB) *GormUserRepo { return &GormUserRepo{db: db} }

// Lấy DB/Tx từ context nếu có, ngược lại dùng db gốc
func (r *GormUserRepo) dbFrom(ctx context.Context) *gorm.DB {
	if tx, ok := utils.TxFrom(ctx); ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *GormUserRepo) Create(ctx context.Context, u *models.User) error {
	return r.dbFrom(ctx).WithContext(ctx).Create(u).Error
}

func (r *GormUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.dbFrom(ctx).WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) Update(ctx context.Context, u *models.User) error {
	return r.dbFrom(ctx).WithContext(ctx).Updates(u).Error
}

func (r *GormUserRepo) Delete(ctx context.Context, id uint) error {
	return r.dbFrom(ctx).WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *GormUserRepo) List(ctx context.Context, pag *pkg.Pagination, search string) ([]models.User, int64, error) {
	q := r.dbFrom(ctx).WithContext(ctx).Model(&models.User{})
	if search != "" {
		// Lưu ý: ILIKE là của Postgres; nếu test bằng SQLite thì đổi sang LOWER(...) LIKE ...
		q = q.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
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

func (r *GormUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	if err := r.dbFrom(ctx).WithContext(ctx).
		Where("username = ?", username).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
