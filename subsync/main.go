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
	time1 := flag.String("t1", "", "simple delay in from->to format")
	time2 := flag.String("t2", "", "interpolated delay in from->to format")
	flag.Parse()

	_, err := ReadSrtFile(*input)
	if err != nil {
		fmt.Println(err)
	}

	diff1 := parseFromToOpt(*time1)

    if *time2 == "" {
        fmt.Println(diff1)
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

func convertTime(s string) int {
	var hour, minute, second, milli int
	fmt.Sscanf(s, "%2d:%2d:%2d,%3d", &hour, &minute, &second, &milli)

	return (((hour*60)+minute)*60+second)*1000 + milli
}

func parseFromToOpt(ftOpt string) int {
	p := strings.Index(ftOpt, "->")

	if p == -1 {
		return 0
	}

	from := ftOpt[:p]
	to := ftOpt[p+2:]

	fmt.Printf("|%s|%s|\n", from, to)
    fmt.Printf("%d %d\n", convertTime(from), convertTime(to))

    return 0
}
