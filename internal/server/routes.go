package server

import (
	"Book_market_api/internal/controller"
	"Book_market_api/internal/middlewares"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	//
	r.Route("/api/v1", func(r chi.Router) {
		// User
		r.With(middlewares.Authenticate_v1("admin", "user")).
			Get("/User", controller.NewUserController().GetUserPagination)
		r.Post("/User/signup", controller.NewUserController().CreateUser)
		r.With(middlewares.Authenticate_v1("user")).
			Post("/User/update", controller.NewUserController().UpdateUser)
		r.With(middlewares.Authenticate_v1("user")).
			Get("/User/delete/{id}", controller.NewUserController().DeleteUser)
		r.Post("/User/signin", controller.NewUserController().LoginUser)
		r.With(middlewares.Authenticate_v2).
			Get("/User/refresh_token", controller.NewUserController().GetNewToken)
		r.Post("/User/social_media", controller.NewUserController().LoginSocialMedia)
		// Token
		r.Get("/Token/delete/{id}", controller.NewTokenController().DeleteToken_v2)
		r.Get("/Token/revoke/{id}", controller.NewTokenController().TokenRetrieval)
	})

	return r
}
