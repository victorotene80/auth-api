package dto

import "github.com/victorotene80/authentication_api/internal/domain/contracts"

type RefreshResultDTO struct{
	AccessToken contracts.Token
	RefreshToken string
}