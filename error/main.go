package main

import (
	"encoding/json"
	"fmt"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ae AppError) Error() string {
	return fmt.Sprintf("Error in app (%v - %v)", ae.Code, ae.Message)
}

func Load() error {
	return AppError{
		Code:    102,
		Message: "Cannot load",
	}
}

func main() {
	var err = AppError{
		Code:    101,
		Message: "Event not found",
	}

	fmt.Printf("Error is %v\n", err)

	fmt.Printf("Error is %v\n", Load())

	b, _ := json.Marshal(Load())
	fmt.Println(string(b))

	aerr := Load()
	shit, _ := aerr.(AppError)
	b, _ = json.Marshal(shit)
	fmt.Println(string(b))
}
