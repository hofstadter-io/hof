package form

import "fmt"

type FormError interface {
	Field() string
	Message() string
	Value() interface{}
	String() string
}

type Error struct {
	field   string
	message string
	value   interface{}
}

func NewError(message, field string, value interface{}) *Error {
	return &Error{
		field:   field,
		message: message,
		value:   value,
	}
}

func (E *Error) Field() string {
	return E.field
}

func (E *Error) Message() string {
	return E.message
}

func (E *Error) Value() interface{} {
	return E.value
}

func (E *Error) String() string {
	return fmt.Sprintf("FormError: %s - field: %q value: %v", E.message, E.field, E.value)
}

type RequiredError struct {
	*Error
}

func NewRequiredError(field string) *RequiredError {
	return &RequiredError{
		Error: &Error{
			field:   field,
			message: "missing required field",
		},
	}
}

type ValidationError struct {
	*Error
}

func NewValidationError(field, validationMessage string, value interface{}) *ValidationError {
	return &ValidationError{
		Error: &Error{
			field:   field,
			message: fmt.Sprintf("field: %q has bad input %q with value: %v", field, validationMessage, value),
			value:   value,
		},
	}
}
