package bootstrap

import (
	"net/http"

	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest"
	restHandler "github.com/victorotene80/authentication_api/internal/interfaces/http/rest/handler"
	"github.com/victorotene80/authentication_api/internal/infrastructure/validation"
)

func initializeHTTP(commandBus *messaging.CommandBus) http.Handler {
	validate := validation.NewPlaygroundValidator()

	createUserHandler := restHandler.NewCreateUserHandler(commandBus, validate)
	loginHandler := restHandler.NewLoginHandler(commandBus, validate)
	logoutHandler := restHandler.NewLogoutHandler(commandBus, validate)

	router := rest.NewRouter(
		createUserHandler,
		loginHandler,
		logoutHandler,
	)

	return router.Setup()
}
