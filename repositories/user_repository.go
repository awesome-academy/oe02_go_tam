package repositories

import (
	"errors"
	"gorm.io/gorm"
	"oe02_go_tam/models"
	"strings"
)

type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	Update(user *models.User) error
	GetAllUsers() ([]models.User, error)
	FindAllWithSearch(search string, page, limit int) ([]models.User, int64, error)
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *userRepository) FindAllWithSearch(search string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

func (r *userRepository) Delete(id uint) error {
	var booking models.Booking
	if err := r.db.Select("id").Where("user_id = ?", id).Limit(1).First(&booking).Error; err == nil {
		return errors.New("cannot delete user: has related bookings")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check comments
	var comment models.Comment
	if err := r.db.Select("id").Where("user_id = ?", id).Limit(1).First(&comment).Error; err == nil {
		return errors.New("cannot delete user: has related comments")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check reviews
	var review models.Review
	if err := r.db.Select("id").Where("user_id = ?", id).Limit(1).First(&review).Error; err == nil {
		return errors.New("cannot delete user: has related reviews")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check likes
	var like models.Like
	if err := r.db.Select("id").Where("user_id = ?", id).Limit(1).First(&like).Error; err == nil {
		return errors.New("cannot delete user: has liked reviews")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Delete(&models.User{}, id).Error
}
