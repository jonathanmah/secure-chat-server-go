package router

import (
	"chatapp/internal/chat"
	"chatapp/internal/handlers"
	"chatapp/internal/middleware"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// create a router for server
func NewRouter() http.Handler {
	router := chi.NewRouter()
	registerAuthRoutes(router)
	registerWsRoutes(router)
	registerHTMLRoutes(router)
	registerStaticRoutes(router)
	return router
}

// register authentication routes on router
func registerAuthRoutes(r chi.Router) {

	r.Group(func(sub chi.Router) {
		sub.Use(middleware.NoCache)
		r.Get("/login/google", handlers.RedirectOAuthHandler)
		r.Get("/auth/google", handlers.PostOAuthRedirectHandler)
		r.Get("/auth/confirm", handlers.ConfirmEmailHandler)
		r.Post("/auth/logout", handlers.LogoutHandler)
		r.Post("/auth/login", handlers.LoginHandler)
		r.Post("/auth/sign-up", handlers.SignUpHandler)
		r.Post("/auth/forgot-password", handlers.ForgotPasswordHandler)
		r.Post("/auth/reset-password", handlers.ResetPasswordHandler)
	})

	r.Post("/auth/refresh", handlers.RefreshAccessTokenHandler)

	r.With(middleware.AuthenticateAccessToken).Get("/auth/user-info", handlers.GetUserInfoHandler)
	r.With(middleware.AuthenticateAccessToken).Post("/auth/update-username", handlers.UpdateUsernameHandler)
}

// register websocket routes for chat messages
func registerWsRoutes(r chi.Router) {
	hub := chat.NewHub()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWsConn(hub, w, r)
	})
	go hub.Run() // have hub running on its own thread
}

// register routes for serving static frontend content
func registerStaticRoutes(r chi.Router) {
	publicDir := filepath.Join("frontend", "public") // relative to current directory
	fs := http.FileServer(http.Dir(publicDir))
	r.Handle("/static/*", http.StripPrefix("/", fs)) // relative to address the router is listening
}

// register routes to serve html
func registerHTMLRoutes(r chi.Router) {
	publicDir := filepath.Join("frontend", "public", "html")

	r.Group(func(sub chi.Router) {
		sub.Use(middleware.NoCache)
		sub.Get("/", serveFile(publicDir, "login.html"))
		sub.Get("/sign-up", serveFile(publicDir, "sign-up.html"))
		sub.Get("/forgot-password", serveFile(publicDir, "forgot-password.html"))
	})

	r.With(middleware.NoCache).Get("/lobby", serveFile(publicDir, "lobby.html"))
	r.With(middleware.AuthenticateQueryParamToken, middleware.NoCache).Get("/reset-password", serveFile(publicDir, "reset-password.html"))
}
