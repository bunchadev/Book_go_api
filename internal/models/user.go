package models

type UserPagination struct {
	Id            string  `json:"id" binding:"required"`
	Username      string  `json:"username" binding:"required"`
	Email         string  `json:"email" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	Balance       float64 `json:"balance"`
	Login_enabled bool    `json:"login_enabled"`
	Depot_limit   int     `json:"depot_limit"`
	Auth_method   string  `json:"auth_method"`
	Create_at     string  `json:"create_at"`
	Role          string  `json:"role" binding:"required"`
}

type UserCreate struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Auth_method string `json:"auth_method"`
	Role        string `json:"role" binding:"required"`
}

type UserResponse struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Balance     string `json:"balance"`
	Depot_limit int    `json:"depot_limit"`
	Auth_method string `json:"auth_method"`
	Role        string `json:"role"`
}

type TokenResponse struct {
	Access_token  string       `json:"access_token"`
	Refresh_token string       `json:"refresh_token"`
	Expires_in    int          `json:"expires_in"`
	User          UserResponse `json:"user"`
}

type UpdateUser struct {
	Id            string  `json:"id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Balance       float64 `json:"balance"`
	Login_enabled bool    `json:"login_enabled"`
	Depot_limit   int     `json:"depot_limit"`
	Role          string  `json:"role"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SocialMedia struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}
