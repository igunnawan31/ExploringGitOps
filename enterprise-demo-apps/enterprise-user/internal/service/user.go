package service

import (
	"github.com/igunnawan31/enterprise-user/internal/domain"
	"github.com/igunnawan31/enterprise-user/internal/repository"
)

type UserService interface {
	CreateUser(name, email, password string, role domain.UserRole, phone string, status domain.UserStatus) (repository.User, error)
	GetAllUsers() ([]repository.User, error)
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
	if role == "" {
		role = domain.RoleCustomer
	}

	user := repository.User{
		Name:         name,
		Email:        email,
		PasswordHash: password,
		Role:         role,
		Phone:        phone,
		Status:       domain.Active,
	}

	err := s.repo.Save(&user)
	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

func (s *userService) GetAllUsers() ([]repository.User, error) {
	return s.repo.FindAll()
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
