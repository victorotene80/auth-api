package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/handler"
	appmw "github.com/victorotene80/authentication_api/internal/interfaces/middleware"
)

type Router struct {
	mux *chi.Mux

	Logger         *zap.Logger
	AuthMiddleware *appmw.AuthMiddleware

	CreateUserHandler *handler.CreateUserHandler
	LoginHandler      *handler.LoginHandler
	// LogoutHandler     *handler.LogoutHandler
}

func NewRouter(
	createUserHandler *handler.CreateUserHandler,
	authMiddleware *appmw.AuthMiddleware,
	logger *zap.Logger,
	loginHandler *handler.LoginHandler,
	// logoutHandler *handler.LogoutHandler,
) *Router {
	return &Router{
		mux:               chi.NewRouter(),
		Logger:            logger,
		CreateUserHandler: createUserHandler,
		AuthMiddleware:    authMiddleware,
		LoginHandler:      loginHandler,
		// LogoutHandler:     logoutHandler,
	}
}

func (rt *Router) Setup() http.Handler {
	rt.mux.Use(chiMw.RequestID) // injects request IDs
	rt.mux.Use(chiMw.RealIP)    // extracts real IP
	rt.mux.Use(appmw.PanicRecovery(rt.Logger))
	rt.mux.Use(appmw.RequestMetadata)
	rt.mux.Use(chiMw.Logger) // optional: Chi's default logger (can remove if noisy)

	rt.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	rt.mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	rt.mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
			rt.CreateUserHandler.CreateUser(w, r)
		})

		r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
			rt.LoginHandler.Login(w, r)
		})

		r.Group(func(r chi.Router) {
			r.Use(rt.AuthMiddleware.Handle)

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