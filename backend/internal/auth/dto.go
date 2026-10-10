package auth

import "time"

type UserResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Balance     int64  `json:"balance"`
	AccountType string `json:"account_type"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type registerRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccountType string `json:"account_type"`
}
