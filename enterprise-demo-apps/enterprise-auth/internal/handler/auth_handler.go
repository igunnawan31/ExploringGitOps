package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-auth/internal/service"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Function to directing user to Google Login Page
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	googleConfig := h.authService.GoogleOAuthConfig()

	// Generate Auth URL for Google Login using state-string-random to prevent CSRF attacks
	url := googleConfig.AuthCodeURL("state-string-random", oauth2.AccessTypeOffline)

	// Redirect user to Google Login Page
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	email, name, jwtToken, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully logged in with Google",
		"email":        email,
		"name":         name,
		"access_token": jwtToken,
	})
}

func (h *AuthHandler) RegisterManual(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	err := h.authService.RegisterUserConfig(input.Email, input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration event sent successfully! Please login.",
		"email":   input.Email,
		"name":    input.Name,
	})
}

func (h *AuthHandler) LoginManual(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	name, jwtToken, err := h.authService.LoginUserConfig(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully logged in!",
		"email":        input.Email,
		"name":         name,
		"access_token": jwtToken,
	})
}
