package handler

import (
	"encoding/json"
	"net/http"

	"github.com/victorotene80/authentication_api/internal/application/command"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/request"
	"github.com/victorotene80/authentication_api/internal/interfaces/http/rest/response"
)

type LoginHandler struct {
	commandBus *messaging.CommandBus
	validator  contracts.Validator
}

func NewLoginHandler(commandBus *messaging.CommandBus, validator contracts.Validator) *LoginHandler {
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.LoginCommand{
		Email:     req.Email,
		Password:  req.Password,
		DeviceID:  req.DeviceID,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
	}

	result, err := messaging.Execute[command.LoginCommand, dto.LoginResultDTO](h.commandBus, r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
