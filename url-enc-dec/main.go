package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
)

func main() {
	input := bufio.NewReader(os.Stdin)

	for {
		line, err := input.ReadString('\n')
		if err != nil {
			panic(err)
		}

		decoded, err := url.PathUnescape(line)
		if err != nil {
			panic(err)
		}

		fmt.Println(decoded)
	}
}
