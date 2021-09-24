package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type LogRecord struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	File    string    `json:"file"`
	Line    string    `json:"line"`
	Message string    `json:"message"`

	Method       string `json:"method"`
	Host         string `json:"host"`
	URI          string `json:"uri"`
	Status       int    `json:"status"`
	Error        string `json:"error"`
	LatencyHuman string `json:"latency"`
}

func main() {
	input := bufio.NewReader(os.Stdin)

	for {
		lineBytes, err := input.ReadBytes('\n')
		if err != nil {
			break
		}

		r := LogRecord{}
		err = json.Unmarshal(lineBytes, &r)
		if err != nil {
			fmt.Println(err)
		}

		if r.Level != "" {
			fmt.Printf("%s [%s] %s:%s %s\n", r.Time, r.Level, r.File, r.Line, r.Message)
		} else {
			color.HiBlue("%s %s %d %s %s %s %s\n", r.Time, r.Method, r.Status, r.Host, r.URI, r.Error, r.LatencyHuman)
		}
	}
}
