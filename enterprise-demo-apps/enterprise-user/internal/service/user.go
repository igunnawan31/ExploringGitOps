package service

import (
	"fmt"

	"github.com/igunnawan31/enterprise-user/internal/domain"
	"github.com/igunnawan31/enterprise-user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(name, email, password string, role domain.UserRole, phone string, status domain.UserStatus) (repository.User, error)
	GetAllUsers(page, limit int, search string) ([]repository.User, int64, error)
	GetUserByID(id string) (repository.User, error)
	UpdateUser(id string, updates map[string]interface{}) (repository.User, error)
	DeleteUser(id string) error
}

type userService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(name, email, password string, role domain.UserRole, phone string, status domain.UserStatus) (repository.User, error) {
	// Jika role tidak dipilih, makan otomatis rolenya akan menjadi customer
	if role == "" {
		role = domain.RoleCustomer
	}

	// Menambahkan hashing password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return repository.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user := repository.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedBytes),
		Role:         role,
		Phone:        phone,
		Status:       domain.Active,
	}

	err = s.repo.Save(&user)
	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(page, limit int, search string) ([]repository.User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(limit, offset, search)
}

func (s *userService) GetUserByID(id string) (repository.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id string, updates map[string]interface{}) (repository.User, error) {
	return s.repo.Update(id, updates)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
