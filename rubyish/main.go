package main

import (
	"errors"
	"fmt"
)

func handler(params ...string) error {
	v := validate()
	v.validateNotEmpty(params[0])
	v.validateLength(params[1], 0, 20)

	if v.hasErrors() {
		return errors.New(v.displayErrors())
	}

	return nil
}

type validation struct {
	result []error
}

func validate() *validation {
	return &validation{
		result: make([]error, 0),
	}
}

func (v *validation) validateNotEmpty(s string) *validation {
	if s == "" {
		v.result = append(v.result, errors.New("string is empty"))
	}

	return v
}

func (v *validation) validateLength(s string, min, max int) *validation {
	paramLength := len(s)

	if paramLength >= min && paramLength <= max {
		return v
	}

	v.result = append(v.result, errors.New("string is too short or long"))

	return v
}

func (v *validation) hasErrors() bool {
	return len(v.result) > 0
}

func (v *validation) displayErrors() string {
	s := ""

	for _, err := range v.result {
		s += err.Error() + " "
	}

	return s
}

func main() {
	if err := handler("", "desc", "price"); err != nil {
		fmt.Println(err)
	}
}
