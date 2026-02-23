package handler

import (
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
)

type LogoutHandler struct {
	commandBus *messaging.CommandBus
	validator  contracts.Validator
}

func NewLogoutHandler(commandBus *messaging.CommandBus, validator contracts.Validator) *LogoutHandler {
	return &LogoutHandler{
		commandBus: commandBus,
		validator:  validator,
	}
}
