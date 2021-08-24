package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	b := []byte(`
    {"field_a": 1, "field_2": [{"a": 2}, {"b": 3}]}
    `)

	var i interface{}

	if err := json.Unmarshal(b, &i); err != nil {
		panic(err)
	}

	fmt.Printf("Object %v\n", i)
}
