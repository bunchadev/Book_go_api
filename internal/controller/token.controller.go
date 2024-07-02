package controller

import (
	"Book_market_api/internal/service"
	"Book_market_api/response"
	"net/http"
)

type TokenController struct {
	tokenService *service.TokenService
}

func NewTokenController() *TokenController {
	return &TokenController{
		tokenService: service.NewTokenService(),
	}
}

func (tr *TokenController) DeleteToken_v2(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	if userId == "" {
		response.ErrorResponse(w, 301)
		return
	}
	err := tr.tokenService.DeleteToken_v2(userId)
	if err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, nil)
}

func (tr *TokenController) TokenRetrieval(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	if userId == "" {
		response.ErrorResponse(w, 301)
		return
	}
	err := tr.tokenService.TokenRetrieval(userId)
	if err != nil {
		response.ErrorResponse(w, 302)
		return
	}
	response.SuccessResponse(w, 300, nil)
}
