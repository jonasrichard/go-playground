package sub

import (
	"fmt"
	"os"
)

func LinearMode(opts GlobalOpts) {
	subtitles, err := ReadSrtFile(opts.Input)
	if err != nil {
        fmt.Printf("Error opening %s: %v\n", opts.Input, err)

        os.Exit(1)
	}

	diff1 := parseFromToOpt(opts.Linear.StartTime)
    // diff2 := parseFromToOpt(t2)

    // TODO diff2-diff1 but checking the times

	for i := range subtitles {
		newFrom := convertTimeToString(convertStringToTime(subtitles[i].From) + diff1)
		newTo := convertTimeToString(convertStringToTime(subtitles[i].To) + diff1)

		subtitles[i].From = newFrom
		subtitles[i].To = newTo
	}

	err = WriteSrtFile(subtitles, opts.Output)

	if err != nil {
		fmt.Println(err)

        os.Exit(1)
	}
}
