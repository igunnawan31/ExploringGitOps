package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/igunnawan31/enterprise-auth/internal/broker"
	"github.com/igunnawan31/enterprise-auth/internal/repository"
	"github.com/igunnawan31/enterprise-auth/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Function to create a new Google OAuth2 Configuration
func (s *AuthService) GoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// Function to manual register user to the database and publish an event to Kafka
func (s *AuthService) RegisterUserConfig(email, name, password string) error {
	err := broker.PublishUserRegistrationEvent(broker.UserRegistrationEvent{
		Email:        email,
		Name:         name,
		PasswordHash: password,
		Provider:     "manual",
	})

	return err
}

// Function to handle Google OAuth2 Callback
func (s *AuthService) HandleGoogleCallback(code string) (string, string, string, error) {
	googleConfig := s.GoogleOAuthConfig()

	// Exchange the authorization code for an access token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", "", "", errors.New("failed to exchange token with Google")
	}

	// Use the access token to create an HTTP client for making requests to Google APIs (fetching user profile information such as email and name)
	client := googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", "", errors.New("failed to get user info from Google")
	}
	defer resp.Body.Close()

	// Decode the response body to extract user information
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", "", "", errors.New("failed to parse google user info")
	}

	// Sent event to kafka topic "user-registered" using the PublishUserLoginEvent function from the broker package
	_ = broker.PublishUserRegistrationEvent(broker.UserRegistrationEvent{
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		Provider: "google",
	})

	// JWT Token Internal for user login
	jwtToken, err := utils.GenerateTokenJWT(userInfo.Email, userInfo.Name)
	if err != nil {
		return "", "", "", errors.New("failed to generate jwt token")
	}

	// Return the user information as a JSON responsse
	return userInfo.Email, userInfo.Name, jwtToken, nil
}

func (s *AuthService) LoginUserConfig(email, password string) (string, string, error) {
	userInternalURL := os.Getenv("ENTERPRISE_USER_URL")
	if userInternalURL == "" {
		userInternalURL = "http://enterprise-user:8080/api/v1/internal/verify"
	}

	reqBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	resp, err := http.Post(userInternalURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", "", errors.New("failed to connect to user service for verification")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("invalid email or password")
	}

	var apiResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", "", errors.New("failed to parse verification response")
	}

	// Ambil nama dari dalam struct Data
	jwtToken, err := utils.GenerateTokenJWT(email, apiResponse.Data.Name)
	if err != nil {
		return "", "", errors.New("failed to generate jwt token")
	}

	return apiResponse.Data.Name, jwtToken, nil
}
