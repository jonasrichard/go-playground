package sub

import (
	"fmt"
	"os"
)

func DelayMode(opts GlobalOpts) {
	from, to := parseFromToOpt(opts.Delay.Time)

    diff := to - from

	subtitles, err := ReadSrtFile(opts.Input)
	if err != nil {
        fmt.Printf("Error opening %s: %v\n", opts.Input, err)

        os.Exit(1)
	}

	for i := range subtitles {
		newFrom := convertTimeToString(convertStringToTime(subtitles[i].From) + diff)
		newTo := convertTimeToString(convertStringToTime(subtitles[i].To) + diff)

		subtitles[i].From = newFrom
		subtitles[i].To = newTo
	}

	err = WriteSrtFile(subtitles, opts.Output)

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
