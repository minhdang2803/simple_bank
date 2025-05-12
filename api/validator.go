package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/minhdang2803/simple_bank/utils"
)

var validatorCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currenct, isOk := fieldLevel.Field().Interface().(string); isOk {
		isSupport := utils.IsSupportedCurrency(currenct)
		fmt.Print(isSupport)
		return isSupport
	}
	return false
}
