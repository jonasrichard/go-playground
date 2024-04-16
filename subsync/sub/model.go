package sub

import (
	"bufio"
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

// Delay a subtitle item in place
func (si *SubtitleItem) Delay(diff int) {
    si.FromMS += diff
    si.ToMS += diff

    si.From = convertTimeToString(si.FromMS)
    si.To = convertTimeToString(si.ToMS)
}

func (si *SubtitleItem) SetFromMS(f int) {
    si.FromMS = f
    si.From = convertTimeToString(si.FromMS)
}

func (si *SubtitleItem) SetToMS(f int) {
    si.ToMS = f
    si.To = convertTimeToString(si.ToMS)
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
		subtitles = make([]SubtitleItem, 0)
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

// parseFromToOpt converts hh:mm:ss,msc->hh:mm:ss,msc into two ints representing
// the millis of the two times
func parseFromToOpt(ftOpt string) (int, int) {
	p := strings.Index(ftOpt, "->")

	if p == -1 {
        fmt.Printf("Time format error: %s\n", ftOpt)

        os.Exit(1)
	}

	from := ftOpt[:p]
	to := ftOpt[p+2:]

	return convertStringToTime(from), convertStringToTime(to)
}
