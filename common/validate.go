package common

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s any) error
	ValidateField(f any, tag string) error
}

type validatorImpl struct {
	validator *validator.Validate
}

func NewValidator() Validator {
	return &validatorImpl{
		validator: validator.New(),
	}
}

func (v *validatorImpl) ValidateStruct(s any) error {
	if s == nil {
		return errors.New("validation failed on nil struct")
	}

	err := v.validator.Struct(s)

	if err != nil {
		return v.buildValidationError(err)
	}

	return nil
}

func (v *validatorImpl) ValidateField(f any, tag string) error {
	err := v.validator.Var(f, tag)

	if err != nil {
		return v.buildValidationError(err)
	}

	return nil
}

func (v *validatorImpl) buildValidationError(err error) error {
	msg := ""
	for i, err := range err.(validator.ValidationErrors) {
		if i != 0 {
			msg += " and "
		}
		msg = msg + fmt.Sprintf("validation failed on field %v", err.Field())
	}
	return errors.New(msg)
}
