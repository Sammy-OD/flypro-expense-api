package validators

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var categories = map[string]bool{
	"travel": true, "meals": true, "office": true, "supplies": true,
}

var iso4217 = map[string]bool{
	"USD": true, "EUR": true, "GBP": true, "NGN": true, "CAD": true, "AUD": true, "JPY": true,
}

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
			val := strings.ToUpper(fl.Field().String())
			return iso4217[val]
		})
		_ = v.RegisterValidation("category", func(fl validator.FieldLevel) bool {
			val := strings.ToLower(fl.Field().String())
			return categories[val]
		})
	}
}
