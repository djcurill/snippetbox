package validators

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) CheckField(ok bool, field, err string) {
	if !ok {
		v.AddFieldError(field, err)
	}
}

func (v *Validator) AddFieldError(field, err string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = err
	}
}

// false if empty string
func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// maximum allowable character
func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) < n
}

// methods do not work with type params, this must be a function
func PermittedValue[T comparable](v T, pvs ...T) bool {
	return slices.Contains(pvs, v)
}
