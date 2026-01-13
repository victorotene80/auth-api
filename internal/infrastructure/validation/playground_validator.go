package validation

import (
	"github.com/go-playground/validator/v10"
)

type PlaygroundValidator struct {
	v *validator.Validate
}

func NewPlaygroundValidator() *PlaygroundValidator {
	return &PlaygroundValidator{
		v: validator.New(),
	}
}

func (p *PlaygroundValidator) Struct(i any) error {
	return p.v.Struct(i)
}
