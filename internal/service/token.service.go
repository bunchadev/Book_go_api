package service

import "Book_market_api/internal/repo"

type TokenService struct {
	tokenRepo *repo.TokenRepo
}

func NewTokenService() *TokenService {
	return &TokenService{
		tokenRepo: repo.NewTokeRepo(),
	}
}

func (ts *TokenService) DeleteToken_v2(userId string) error {
	return ts.tokenRepo.DeleteToken_v2(userId)
}

func (ts *TokenService) TokenRetrieval(userId string) error {
	return ts.tokenRepo.TokenRetrieval(userId)
}
