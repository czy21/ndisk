package web

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

func DateValidator(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	fmt.Println(date, ok)
	return true
}
