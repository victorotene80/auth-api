package bootstrap

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	appmw "github.com/victorotene80/authentication_api/internal/interfaces/middleware"
	"github.com/victorotene80/authentication_api/internal/infrastructure/validation"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest"
	restHandler "github.com/victorotene80/authentication_api/internal/interfaces/http/rest/handler"
)

func initializeHTTP(
	commandBus *messaging.CommandBus,
	authSvc contracts.AuthService,
	logger *zap.Logger,
) http.Handler {
	validate := validation.NewPlaygroundValidator()

	loginHandler := restHandler.NewLoginHandler(commandBus, validate)
	createUserHandler := restHandler.NewCreateUserHandler(commandBus, validate)
	authMiddleware := appmw.NewAuthMiddleware(authSvc)

	router := rest.NewRouter(
		createUserHandler,
		authMiddleware,
		logger,
		loginHandler,
	)

	return router.Setup()
}