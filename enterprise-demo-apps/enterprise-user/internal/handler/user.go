package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-user/internal/domain"
	"github.com/igunnawan31/enterprise-user/internal/repository"
	"github.com/igunnawan31/enterprise-user/internal/service"
	"github.com/igunnawan31/enterprise-user/internal/utils"
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
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Invalid Input")
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password, domain.UserRole(req.Role), req.Phone, domain.Active)
	if err != nil {
		// Mengembalikan status 400 Bad Request jika gagal (misal: email duplikat)
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Duplicate Email or Invalid Input")
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "User successfully created", user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := h.userService.GetAllUsers(page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", "Internal Server Error")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", gin.H{
		"users": users,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_data":   total,
		},
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "400", "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Invalid Input")
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
			utils.ErrorResponse(c, http.StatusNotFound, "400", "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", "Internal Server Error")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users updated successfully", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "400", "User not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", "Internal Server Error")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Users deleted successfully", nil)
}
