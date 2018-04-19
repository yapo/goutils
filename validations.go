package goutils

import (
	"reflect"
	"regexp"
	"time"
)

// ValidatorFunc define a common type on all proyects
type ValidatorFunc func(interface{}) (bool, string)

func isString(param interface{}) bool {
	return reflect.TypeOf(param).Kind() == reflect.String
}

// IntegerValidator checks that param is numeric
func IntegerValidator(param interface{}) (bool, string) {
	//checking if a variable is an int checking if it is a float64 is the moda
	if reflect.TypeOf(param).Kind() != reflect.Float64 && reflect.TypeOf(param).Kind() != reflect.Int {
		return false, "is not an integer"
	}
	return true, ""
}

// StringValidator checks that param is a string
func StringValidator(param interface{}) (bool, string) {
	if !isString(param) {
		return false, "is not a string"
	}
	return true, ""
}

// EmptyStringValidator checks that param is not an empty string
func EmptyStringValidator(param interface{}) (bool, string) {
	if !isString(param) || param.(string) == "" {
		return false, "can not be an empty string"
	}
	return true, ""
}

// PlatesValidator checks that param is formated correctly as a car plate
func PlatesValidator(param interface{}) (bool, string) {
	var regexValidator = regexp.MustCompile("^[A-Za-z]{2}([A-Za-z]{1,2}0?|[0-9]{1,2})[0-9]{2}$")
	if !isString(param) || !regexValidator.MatchString(param.(string)) {
		return false, "is not in the required format"
	}
	return true, ""
}

// EmailValidator checks that param is formated correctly as a email
func EmailValidator(param interface{}) (bool, string) {
	var regexValidator = regexp.MustCompile("^[^@]+@[a-zA-Z0-9]+(\\.[a-zA-Z0-9]+)+$")
	if !isString(param) || !regexValidator.MatchString(param.(string)) {
		return false, "is not in the required format"
	}
	return true, ""
}

// ValidateDate Check if date YYYY-mm-dd is valid
func ValidateDate(param interface{}) (bool, string) {
	if !isString(param) {
		return false, "is not a string"
	}
	date := param.(string)
	_, e := time.Parse(time.RFC3339, date+"T00:00:00Z")
	if e != nil {
		return false, "is not a date"
	}
	return true, ""
}
