package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Parse(jsonText string) interface{} {
	var i interface{}

	if err := json.Unmarshal([]byte(jsonText), &i); err != nil {
		panic(err)
	}

	return i
}

// [0].markets.[2].id
func SplitPath(path string) error {
	segments := strings.Split(path, ".")

	for _, segment := range segments {
		if segment == "" {
			return errors.New("empty path segment")
		}

		if segment[0] == '[' {
			// array reference
			num := segment[1 : len(segment)-1]
			index, err := strconv.Atoi(num)
			if err != nil {
				return err
			}

			fmt.Printf("index %d\n", index)
		} else {
			// object field reference
			fmt.Printf("field %s\n", segment)
		}
	}

	return nil
}

func main() {
	j := `[
    {"event_id": 2, "success": true}
    ]`

	fmt.Println(Parse(j))

	SplitPath("[1].event_id")
}
