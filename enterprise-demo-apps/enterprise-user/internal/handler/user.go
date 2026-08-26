package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-user/internal/domain"
	"github.com/igunnawan31/enterprise-user/internal/repository"
	"github.com/igunnawan31/enterprise-user/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// DTO Request Handler
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password_policy"`
	Role     string `json:"role" binding:"omitempty,oneof=admin cashier customer manager"`
	Phone    string `json:"phone" binding:"omitempty"`
}

type UpdateUserRequest struct {
	Name   *string            `json:"name" binding:"omitempty,min=3"`
	Email  *string            `json:"email" binding:"omitempty,email"`
	Role   *domain.UserRole   `json:"role" binding:"omitempty,oneof=admin cashier customer manager"`
	Phone  *string            `json:"phone" binding:"omitempty"`
	Status *domain.UserStatus `json:"status" binding:"omitempty,oneof=active inactive"`
}

// RegisterRoutes links the HTTP methods and endpoints to the handler methods
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", h.Create)
	rg.GET("/users", h.GetAll)
	rg.GET("/users/:id", h.GetByID)
	rg.PATCH("/users/:id", h.Update)
	rg.DELETE("/users/:id", h.Delete)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password, domain.UserRole(req.Role), req.Phone, domain.Active)
	if err != nil {
		// Mengembalikan status 400 Bad Request jika gagal (misal: email duplikat)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": user})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	user, err := h.userService.UpdateUser(id, updates)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user successfully deleted"})
}
