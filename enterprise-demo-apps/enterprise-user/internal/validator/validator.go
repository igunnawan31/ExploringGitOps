package validator

import (
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password_policy", func(fl validator.FieldLevel) bool {
			password := fl.Field().String()

			if len(password) < 8 {
				return false
			}

			var hasUpper, hasLower, hasNumber, hasSpecial bool

			for _, char := range password {
				switch {
				case unicode.IsUpper(char):
					hasUpper = true
				case unicode.IsLower(char):
					hasLower = true
				case unicode.IsNumber(char):
					hasNumber = true
				case unicode.IsPunct(char) || unicode.IsSymbol(char):
					hasSpecial = true
				}
			}

			// Wajib memenuhi semua kriteria
			return hasUpper && hasLower && hasNumber && hasSpecial
		})
	}
}
