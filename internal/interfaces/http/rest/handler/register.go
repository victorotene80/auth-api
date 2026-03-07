package handler

import (
	"encoding/json"
	"net/http"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/requestctx"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/httperr"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/request"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/response"
)

type CreateUserHandler struct {
	commandBus *messaging.CommandBus
	validator  appContracts.Validator
}

func NewCreateUserHandler(commandBus *messaging.CommandBus, validator appContracts.Validator) *CreateUserHandler {
	return &CreateUserHandler{
		commandBus: commandBus,
		validator:  validator,
	}
}

func (h *CreateUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := new(request.CreateUserRequest)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"Invalid JSON payload",
			nil,
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		// You can later parse this into []response.ValidationError
		response.Error(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"One or more fields are invalid",
			err.Error(), // or a structured list
		)
		return
	}

	meta, _ := requestctx.MetaFrom(r.Context())

	cmd := command.CreateUserCommand{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		MiddleName:  req.MiddleName,
		Phone:       req.Phone,
		Locale:      req.Locale,
		Timezone:    req.Timezone,
		AcceptTerms: req.AcceptTerms,

		IPAddress:         meta.IPAddress,
		UserAgent:         meta.UserAgent,
		DeviceID:          meta.DeviceID,
		RequestID:         meta.RequestID,
		DeviceFingerprint: meta.DeviceFingerprint,
		DeviceName:        meta.DeviceName,
	}

	result, err := messaging.Execute[
		command.CreateUserCommand,
		*dto.RegisterUserResponseDTO, // pointer result to match your handler
	](h.commandBus, r.Context(), cmd)
	if err != nil {

		if err != nil {
			response.Error(
				w,
				httperr.StatusFrom(err),
				"USER_CREATION_FAILED",
				"Could not create user",
				err.Error(),
			)

			return
		}
	}

	resp := response.CreateUserResponse{
		UserID:                result.UserID,
		Email:                 result.Email,
		FirstName:             result.FirstName,
		LastName:              result.LastName,
		AccessToken:           result.AccessToken,
		AccessTokenExpiresAt:  result.AccessTokenExpiresAt,
		RefreshToken:          result.RefreshToken,
		RefreshTokenExpiresAt: result.RefreshTokenExpiresAt,
	}

	response.Success(
		w,
		http.StatusCreated,
		"USER_CREATED",
		"User registered successfully",
		&resp,
	)
}
