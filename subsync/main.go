package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Number = 0
	Time   = 1
	Text   = 2
)

type SubtitleItem struct {
	Number int
	From   string
	To     string
	Lines  []string
}

func main() {
	input := flag.String("i", "", ".srt input file")
	flag.Parse()

	_, err := ReadSrtFile(*input)
	if err != nil {
		fmt.Println(err)
	}
}

func ReadSrtFile(name string) ([]SubtitleItem, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		step int
		item SubtitleItem
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch step {
		case Number:
			item.Number, err = strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			step++

		case Time:
			p := strings.Index(line, " --> ")
			item.From = line[:p]
			item.To = line[p+5:]
			step++

		case Text:
			if line != "" {
				item.Lines = append(item.Lines, line)
			} else {
				fmt.Println(item)
				item = SubtitleItem{Lines: []string{}}
				step = 0
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return []SubtitleItem{}, nil
}
