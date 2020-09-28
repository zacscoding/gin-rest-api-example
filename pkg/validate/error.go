package validate

import (
	"fmt"
	"gin-rest-api-example/pkg/logging"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ValidationErrDetail struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

// NewValidationErrorDetail returns ValidationErrDetail list with given validation errors
func ValidationErrorDetails(obj interface{}, tag string, errs validator.ValidationErrors) []*ValidationErrDetail {
	if len(errs) == 0 {
		return []*ValidationErrDetail{}
	}
	var errors []*ValidationErrDetail
	e := reflect.TypeOf(obj).Elem()
	for _, err := range errs {
		f, _ := e.FieldByName(err.Field())
		tagName, _ := f.Tag.Lookup(tag)
		val := err.Value()
		var message string

		switch err.ActualTag() {
		case "required":
			message = fmt.Sprintf("required %s", tagName)
		case "email":
			message = "required email format"
		case "min":
			message = fmt.Sprintf("%s required at least %s length", tagName, err.Param())
		case "hexadecimal":
			message = "required hexadecimal format"
		case "gte":
			message = fmt.Sprintf("greater than or quauls to %s", err.Param())
		case "numeric":
			message = fmt.Sprintf("%s must be numeric", tagName)
		default:
			logging.DefaultLogger().Warnf("unknown validation tag. tag:%s", err.ActualTag())
			message = fmt.Sprintf("invalid %s", tagName)
		}

		errors = append(errors, &ValidationErrDetail{
			Field:   tagName,
			Value:   val,
			Message: message,
		})
	}
	return errors
}
