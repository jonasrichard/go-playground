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
	FromMS int
	ToMS   int
	Lines  []string
}

func main() {
	input := flag.String("i", "", ".srt input file")
	flag.Parse()

	subtitles, err := ReadSrtFile(*input)
	if err != nil {
		fmt.Println(err)
	}

	avg := 0
	n := 0

	for i := range subtitles {
		item := subtitles[i]

		delta := item.ToMS - item.FromMS
		avg += delta
		n++

		if delta < 1500 {
			fmt.Printf("%v %dms\n", item, delta)

			before := item.FromMS
			if i > 0 {
				before = subtitles[i-1].ToMS
			}

			after := item.ToMS
			if i < len(subtitles)-1 {
				after = subtitles[i+1].FromMS
			}

			fmt.Printf("%d %d\n", item.FromMS-before, after-item.ToMS)

			// We need to add this ms before and after the timeframe
			plus := (1500 - delta) / 2

			if plus > item.FromMS-before {
				subtitles[i].FromMS = before + 1
			} else {
				subtitles[i].FromMS -= plus
			}

			if plus > after-item.ToMS {
				subtitles[i].ToMS = after - 1
			} else {
				subtitles[i].ToMS += plus
			}

			subtitles[i].From = convertTimeToString(subtitles[i].FromMS)
			subtitles[i].To = convertTimeToString(subtitles[i].ToMS)

			fmt.Printf("%v\n", subtitles[i])
		}
	}

	fmt.Printf("%vms\n", avg/n)

	WriteSrtFile(subtitles, "out.srt")
}

// TODO check if any time will be negative
// TODO support simple delay mode like -d 1500ms -d 1,5s

func main2() {
	input := flag.String("i", "", ".srt input file")
	time1 := flag.String("t1", "", "simple delay in from->to format")
	time2 := flag.String("t2", "", "interpolated delay in from->to format")
	flag.Parse()

	subtitles, err := ReadSrtFile(*input)
	if err != nil {
		fmt.Println(err)
	}

	diff1 := parseFromToOpt(*time1)
	fmt.Printf("Applying diff %dms\n", diff1)

	if *time2 == "" {
		fmt.Println(diff1)
	}

	for i := range subtitles {
		fmt.Println(subtitles[i])
		newFrom := convertTimeToString(convertStringToTime(subtitles[i].From) + diff1)
		newTo := convertTimeToString(convertStringToTime(subtitles[i].To) + diff1)

		subtitles[i].From = newFrom
		subtitles[i].To = newTo
		fmt.Println(subtitles[i])
	}

	p := strings.Index(*input, ".srt")
	output := (*input)[:p] + fmt.Sprintf("%+dms", diff1) + ".srt"
	err = WriteSrtFile(subtitles, output)

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
		step      int
		item      SubtitleItem
		subtitles []SubtitleItem = make([]SubtitleItem, 0)
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
				item.FromMS = convertStringToTime(item.From)
				item.ToMS = convertStringToTime(item.To)

				subtitles = append(subtitles, item)
				// fmt.Println(item)
				item = SubtitleItem{Lines: []string{}}
				step = 0
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subtitles, nil
}

func WriteSrtFile(subtitles []SubtitleItem, name string) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range subtitles {
		fmt.Fprintf(file, "%d\n", item.Number)
		fmt.Fprintf(file, "%s --> %s\n", item.From, item.To)

		for _, line := range item.Lines {
			fmt.Fprintf(file, "%s\n", line)
		}

		fmt.Fprint(file, "\n")
	}

	return nil
}

func convertStringToTime(s string) int {
	var hour, minute, second, milli int

	fmt.Sscanf(s, "%2d:%2d:%2d,%3d", &hour, &minute, &second, &milli)

	return (((hour*60)+minute)*60+second)*1000 + milli
}

func convertTimeToString(t int) string {
	milli := t % 1000
	t /= 1000
	second := t % 60
	t /= 60
	minute := t % 60
	t /= 60

	return fmt.Sprintf("%02d:%02d:%02d,%03d", t, minute, second, milli)
}

func parseFromToOpt(ftOpt string) int {
	p := strings.Index(ftOpt, "->")

	if p == -1 {
		return 0
	}

	from := ftOpt[:p]
	to := ftOpt[p+2:]

	return convertStringToTime(to) - convertStringToTime(from)
}
