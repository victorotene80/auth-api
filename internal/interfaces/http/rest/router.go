package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/handler"
)

type Router struct {
	mux *chi.Mux

	CreateUserHandler *handler.CreateUserHandler
	LoginHandler      *handler.LoginHandler
	LogoutHandler     *handler.LogoutHandler

	// TODO: Add other handlers like RefreshHandler, ChangePasswordHandler, etc.
}

func NewRouter(
	createUserHandler *handler.CreateUserHandler,
	loginHandler *handler.LoginHandler,
	logoutHandler *handler.LogoutHandler,
) *Router {
	return &Router{
		mux:               chi.NewRouter(),
		CreateUserHandler: createUserHandler,
		LoginHandler:      loginHandler,
		LogoutHandler:     logoutHandler,
	}
}

func (rt *Router) Setup() http.Handler {
	rt.mux.Use(middleware.Logger)    // logs requests
	rt.mux.Use(middleware.Recoverer) // recovers panics
	rt.mux.Use(middleware.RequestID) // injects request IDs
	rt.mux.Use(middleware.RealIP)    // extracts real IP
	rt.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"}, //, "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	rt.mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})


	rt.mux.Route("/api/v1", func(r chi.Router) {

		r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
			rt.CreateUserHandler.CreateUser(w, r)
		})
		r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
			rt.LoginHandler.Login(w, r)
		})

		r.Group(func(r chi.Router) {
			// TODO: add authentication middleware here
			// r.Use(AuthMiddleware)

			// Session routes
			//r.Post("/auth/logout", rt.LogoutHandler.Logout)
			r.Post("/auth/refresh", nil)         // TODO
			r.Post("/auth/change-password", nil) // TODO

			// User routes
			r.Get("/users", nil)      // TODO: list users
			r.Get("/users/{id}", nil) // TODO: get user by ID
			r.Put("/users/{id}", nil) // TODO: update user

			// Account & MFA routes
			r.Post("/account/lock", nil)           // TODO
			r.Post("/account/unlock", nil)         // TODO
			r.Post("/account/enable-mfa", nil)     // TODO
			r.Post("/account/verify-mfa", nil)     // TODO
			r.Post("/account/verify-email", nil)   // TODO
			r.Post("/account/reset-password", nil) // TODO
		})
	})

	return rt.mux
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}
