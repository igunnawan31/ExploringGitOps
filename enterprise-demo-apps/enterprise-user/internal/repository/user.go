package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/igunnawan31/enterprise-user/internal/domain"
	"gorm.io/gorm"
)

// Custom Error untuk menandai bahwa user tidak ditemukan
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepositoryInterface untuk mendefinisikan kontrak repository user
type UserRepositoryInterface interface {
	Save(user *User) error
	FindByEmail(email string) (User, error)
	FindAll(limit, offset int, search string) ([]User, int64, error)
	FindByID(id string) (User, error)
	Update(id string, updates map[string]interface{}) (User, error)
	Delete(id string) error
}

type User struct {
	ID           string            `gorm:"primaryKey;type:uuid"`
	Name         string            `gorm:"type:varchar(100);not null"`
	Email        string            `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string            `gorm:"type:varchar(100);column:password_hash;not null" json:"-"`
	Role         domain.UserRole   `gorm:"type:varchar(20);not null;default:'customer'"`
	Phone        string            `gorm:"type:varchar(20)"`
	Status       domain.UserStatus `gorm:"type:varchar(20);not null;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// func (s *UserRepository) Save(name, email string) (User, error) {
// 	user := User{
// 		ID:    uuid.NewString(),
// 		Name:  name,
// 		Email: email,
// 	}

// 	// GORM men-generate query : INSERT INTO users (id, name, email) VALUES (?, ?, ?)
// 	result := s.db.Create(&user)
// 	if result.Error != nil {
// 		return User{}, result.Error
// 	}
// 	return user, result.Error
// }

func (s *UserRepository) Save(user *User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	// GORM men-generate query : INSERT INTO users (id, name, email) VALUES (?, ?, ?)
	result := s.db.Create(user)
	return result.Error
}

func (s *UserRepository) FindByEmail(email string) (User, error) {
	var user User

	result := s.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, result.Error
	}

	return user, nil
}

func (s *UserRepository) FindAll(limit, offset int, search string) ([]User, int64, error) {
	var users []User
	var total int64

	query := s.db.Model(&User{})
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	// GORM men-generate query : SELECT * FROM users
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (s *UserRepository) FindByID(id string) (User, error) {
	var user User
	// GORM men-generate query : SELECT * FROM users WHERE id = ? LIMIT 1
	result := s.db.First(&user, "id = ?", id)

	if result.Error != nil { // Jika error tidak nil
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Jika error adalah tidak ditemukan datanya
			return User{}, ErrUserNotFound
		}

		return User{}, result.Error
	}

	return user, nil
}

func (s *UserRepository) Update(id string, updates map[string]interface{}) (User, error) {
	var user User
	// GORM men-generate query : SELECT * FROM users where id = ? LIMIT 1
	result := s.db.First(&user, "id = ?", id)
	if result.Error != nil { // Jika error tidak nil
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Jika error adalah tidak ditemukan datanya
			return User{}, ErrUserNotFound
		}

		return User{}, result.Error
	}

	if len(updates) == 0 {
		return user, nil // Tidak ada perubahan yang dilakukan
	}

	// GORM men-generate query : UPDATE users SET name = ?, email = ? WHERE id = ?
	resultUpdate := s.db.Model(&user).Updates(updates)
	if resultUpdate.Error != nil {
		return User{}, resultUpdate.Error
	}

	return user, nil
}

func (s *UserRepository) Delete(id string) error {
	// GORM men-generate query: DELETE FROM users WHERE id = ?
	result := s.db.Delete(&User{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return result.Error
}
