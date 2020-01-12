package template

import (
	"errors"
	"fmt"
)

type ExecutionError struct {
	error
}
func NewExecutionError(err error) ExecutionError {
	return ExecutionError{errors.New(fmt.Sprintf("template execution failed: %s", err))}
}

type FileReadError struct {
	error
}
func NewFileReadError(err error) FileReadError {
	return FileReadError{errors.New(fmt.Sprintf("failed to read template file: %s", err))}
}

type CreateTemplateError struct {
	error
}
func NewCreateTemplateError(err error) CreateTemplateError {
	return CreateTemplateError{errors.New(fmt.Sprintf("failed to create template struct: %s", err))}
}

type GoFormatError struct {
	error
}
func NewGoFormatError(err error) GoFormatError {
	return GoFormatError{errors.New(fmt.Sprintf("failed to gofmt file: %s", err))}
}

type FileOpenError struct {
	error
}
func NewFileOpenError(err error) FileOpenError {
	return FileOpenError{errors.New(fmt.Sprintf("unable to open template file: %s", err))}
}

type ParseError struct {
	error
}
func NewParseError(err error) ParseError {
	return ParseError{errors.New(fmt.Sprintf("unable to parse template: %s", err))}
}
