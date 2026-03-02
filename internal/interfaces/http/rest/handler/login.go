// internal/interfaces/http/rest/handler/login.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/requestctx"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/request"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/response"
)

type LoginHandler struct {
	commandBus *messaging.CommandBus
	validator  appContracts.Validator
}

func NewLoginHandler(
	commandBus *messaging.CommandBus,
	validator appContracts.Validator,
) *LoginHandler {
	return &LoginHandler{
		commandBus: commandBus,
		validator:  validator,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := new(request.LoginRequest)

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
		response.Error(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"One or more fields are invalid",
			err.Error(),
		)
		return
	}

	meta, _ := requestctx.MetaFrom(r.Context())

	cmd := command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,

		IPAddress:         meta.IPAddress,
		UserAgent:         meta.UserAgent,
		DeviceID:          meta.DeviceID,
		DeviceFingerprint: meta.DeviceFingerprint,
		DeviceName:        meta.DeviceName,
		RequestID:         meta.RequestID,
	}

	result, err := messaging.Execute[
		command.LoginCommand,
		*dto.LoginResultDTO,
	](h.commandBus, r.Context(), cmd)
	if err != nil {
		// Later: map specific errors to codes (ACCOUNT_LOCKED, INVALID_CREDENTIALS, etc.)
		response.Error(
			w,
			http.StatusUnauthorized,
			"LOGIN_FAILED",
			"Invalid credentials or account locked",
			err.Error(),
		)
		return
	}

	resp := response.LoginResponse{
		MFARequired: result.MFARequired,
		LastLogin:   result.LastLogin,
		Status:      result.Status,
		ChallengeID: result.ChallengeID,
	}

	if !result.MFARequired {
		resp.Tokens = response.AuthTokensResponse{
			AccessToken:  result.Tokens.AccessToken.Value,
			RefreshToken: result.Tokens.RefreshToken.Value,
			ExpiresAt:    result.Tokens.AccessToken.ExpiresAt,
		}
	}

	response.Success(
		w,
		http.StatusOK,
		"LOGIN_SUCCESS",
		"Login successful",
		&resp,
	)
}