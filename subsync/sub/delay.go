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

    delaySubtitles(subtitles, diff)

	err = WriteSrtFile(subtitles, opts.Output)

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func delaySubtitles(subtitles []SubtitleItem, diff int) {
	for i := range subtitles {
        (&subtitles[i]).Delay(diff)
	}
}
