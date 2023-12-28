package api

import (
	"github/tdadadavid/fingreat/utils"

	"github.com/go-playground/validator/v10"
)

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsValidCurrency(currency)
	}
	return false;
}