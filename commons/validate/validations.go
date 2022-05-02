package validate

import (
	"strings"
	"unicode"
)

type FieldError struct {
	Field      string
	Constraint Constraint
	Message    string
	Args       []interface{}
}

type Constraint string

const (
	Required     Constraint = "Required"
	Range        Constraint = "Range"
	Size         Constraint = "Size"
	Like         Constraint = "Like"
	PasswordRule Constraint = "PasswordRule"
	Within       Constraint = "Within"
)

var messages = map[Constraint]string{
	Required:     "Field is required",
	Range:        "Value is not within range",
	Size:         "Size is not within range",
	Like:         "Field must match regex",
	Within:       "Field must be within one of the allowed values",
	PasswordRule: "Must contain atleast one digit, one lower case alphabet, one upper case alphabet and one special character",
}

type ValidationError struct {
	Errors []FieldError
}

func (ve ValidationError) Error() string {
	return "One or more fields are in error"
}

func New() *ValidationError {
	return &ValidationError{[]FieldError{}}
}

func (ve *ValidationError) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationError) IsRequired(field string, value string) *ValidationError {

	if strings.Trim(field, " ") == "" {
		ve.Errors = append(ve.Errors, FieldError{field, Required, messages[Required], nil})
	}

	return ve
}

func (ve *ValidationError) IsRequiredForInt(field string, value int) *ValidationError {

	if value == 0 {
		ve.Errors = append(ve.Errors, FieldError{field, Required, messages[Required], nil})
	}

	return ve
}

func (ve *ValidationError) IsNumberInRange(field string, value int, lower int, upper int) *ValidationError {

	if value < lower || value > upper {
		ve.Errors = append(ve.Errors, FieldError{field, Range, messages[Range], []interface{}{lower, upper}})
	}

	return ve
}

func (ve *ValidationError) IsSizeInRange(field string, value string, lower int, upper int) *ValidationError {

	if len(value) < lower || len(value) > upper {
		ve.Errors = append(ve.Errors, FieldError{field, Size, messages[Size], []interface{}{lower, upper}})
	}

	return ve
}

func (ve *ValidationError) IsValidPassword(field string, value string) *ValidationError {
	var (
		isMin   bool
		special bool
		number  bool
		upper   bool
		lower   bool
	)

	for _, c := range value {
		// Optimize perf if all become true before reaching the end
		if special && number && upper && lower && isMin {
			break
		}

		// else go on switching
		switch {
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsNumber(c):
			number = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}

	if !(special && upper && lower && number) {
		ve.Errors = append(ve.Errors, FieldError{field, PasswordRule, messages[PasswordRule], nil})
	}

	return ve

}

func (ve *ValidationError) IsWithin(field string, value string, allowedValues []string) *ValidationError {

	found := false
	for _, allowedValue := range allowedValues {
		if value == allowedValue {
			found = true
			break
		}
	}
	if !found {
		ve.Errors = append(ve.Errors, FieldError{field, Within, messages[Within], []interface{}{allowedValues}})
	}

	return ve
}
