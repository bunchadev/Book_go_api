package controller

import (
	"Book_market_api/internal/middlewares"
	"Book_market_api/internal/models"
	"Book_market_api/internal/service"
	"Book_market_api/response"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (uc *UserController) GetUserPagination(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	order := query.Get("order")
	field := query.Get("field")
	username := query.Get("username")
	email := query.Get("email")
	validField := map[string]bool{
		"username": true,
		"email":    true,
		"asc":      true,
		"desc":     true,
	}
	if order != "" && field != "" {
		if !validField[order] || !validField[field] {
			response.ErrorResponse(w, 301)
			return
		}
	} else {
		order = "asc"
		field = "id"
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		response.ErrorResponse(w, 301)
		return
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		response.ErrorResponse(w, 301)
		return
	}
	users, err := uc.userService.GetUserPagination(page, limit, order, field, username, email)
	if err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, users)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *models.UserCreate
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.ErrorResponse(w, 301)
		return
	}
	defer r.Body.Close()
	if check := uc.userService.CheckUserName(user.Username); check {
		response.ErrorResponse(w, 303)
		return
	}
	err = uc.userService.CreateUser(user)
	if err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, nil)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user *models.UpdateUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.ErrorResponse(w, 301)
		return
	}
	defer r.Body.Close()
	if err := uc.userService.UpdateUser(user); err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, user)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	if userId == "" {
		response.ErrorResponse(w, 301)
		return
	}
	if err := uc.userService.DeleteUser(userId); err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, nil)
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin *models.LoginUser
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		response.ErrorResponse(w, 301)
		return
	}
	defer r.Body.Close()
	user, err := uc.userService.LoginUser(userLogin)
	if err != nil {
		response.ErrorResponse(w, 304)
		return
	}
	response.SuccessResponse(w, 300, user)
}

func (uc *UserController) GetNewToken(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.ContextUserID_v2).(string)
	hour := r.Context().Value(middlewares.ContextHour).(time.Duration)
	if userId == "" {
		response.ErrorResponse(w, 301)
		return
	}
	user, err := uc.userService.GetNewToken(userId, hour)
	if err != nil {
		response.ErrorResponse(w, 305)
		return
	}
	response.SuccessResponse(w, 300, user)
}

func (uc *UserController) LoginSocialMedia(w http.ResponseWriter, r *http.Request) {
	var user *models.SocialMedia
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.ErrorResponse(w, 301)
		return
	}
	defer r.Body.Close()
	token, err := uc.userService.LoginSocialMedia(user)
	if err != nil {
		response.ErrorResponse(w, 304)
		return
	}
	response.SuccessResponse(w, 300, token)
}
