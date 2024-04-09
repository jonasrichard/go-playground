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

	from1, to1 := parseFromToOpt(opts.Linear.StartTime)
    from2, to2 := parseFromToOpt(opts.Linear.EndTime)

    // https://en.wikipedia.org/wiki/Linear_interpolation
    //
    //      to - to1           to2 - to1
    //    -------------  =  ---------------
    //     from - from1      from2 - from1
    //
    //  so
    //
    //                                  to2 - to1
    //   to = to1 + (from - from1) * --------------
    //                                from2 - from1

    coeff := float64(to2 - to1) / float64(from2 - from1)

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
