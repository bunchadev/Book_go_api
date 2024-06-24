package server

import (
	"Book_market_api/internal/controller"
	"Book_market_api/internal/middlewares"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// api v1
	apiV1 := http.NewServeMux()
	// user
	apiV1.Handle("GET /User", middlewares.Authenticate_v1("admin", "user")(http.HandlerFunc(controller.NewUserController().GetUserPagination)))
	apiV1.HandleFunc("POST /User/signup", controller.NewUserController().CreateUser)
	apiV1.Handle("POST /User/update", middlewares.Authenticate_v1("user")(http.HandlerFunc(controller.NewUserController().UpdateUser)))
	apiV1.Handle("GET /User/delete/{id}", middlewares.Authenticate_v1("admin")(http.HandlerFunc(controller.NewUserController().DeleteUser)))
	apiV1.HandleFunc("POST /User/signin", controller.NewUserController().LoginUser)
	apiV1.Handle("GET /User/refresh_token", middlewares.Authenticate_v2(http.HandlerFunc(controller.NewUserController().GetNewToken)))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))
	return mux
}
