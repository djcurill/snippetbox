package validators

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// parsing this regex once and storing at startup is more performanct
// This regex is the one currently recommended by W3C and Web hyptertext application technology group for
// validating email addresses
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) CheckField(ok bool, field, err string) {
	if !ok {
		v.AddFieldError(field, err)
	}
}

func (v *Validator) AddNonFieldError(err string) {
	v.NonFieldErrors = append(v.NonFieldErrors, err)
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

func (v *Validator) MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) > n
}

func (v *Validator) MaxBytes(value string, n int) bool {
	return len(value) < n
}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// methods do not work with type params, this must be a function
func PermittedValue[T comparable](v T, pvs ...T) bool {
	return slices.Contains(pvs, v)
}
