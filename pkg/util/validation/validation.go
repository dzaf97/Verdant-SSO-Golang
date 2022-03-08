package validation

import (
	"fmt"
	"log"

	"gopkg.in/go-playground/validator.v9"
)

func ValidateForm(a interface{}) bool {

	v := validator.New()
	err := v.Struct(a)

	if err != nil {

		log.Println(err.(validator.ValidationErrors))
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return false
		}

		// from here you can create your own error messages in whatever language you wish
		return false
	}

	return true
}
