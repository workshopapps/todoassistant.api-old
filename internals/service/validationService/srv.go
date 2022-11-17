package validationService

import "github.com/go-playground/validator/v10"

type ValidationSrv interface {
	Validate(any) error
}

type validationStruct struct{}

func (v *validationStruct) Validate(a any) error {
	validate := validator.New()
	return validate.Struct(a)
}
