package sub

import (
	"fmt"
	"os"
)

type LinearParam struct {
	startFrom int
	startTo   int
	endFrom   int
	endTo     int
	coeff     float64
}

func LinearMode(opts GlobalOpts) {
	subtitles, err := ReadSrtFile(opts.Input)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", opts.Input, err)

		os.Exit(1)
	}

	param := LinearParam{}

	param.startFrom, param.startTo = parseFromToOpt(opts.Linear.StartTime)
	param.endFrom, param.endTo = parseFromToOpt(opts.Linear.EndTime)

	linearInterpolation(subtitles, param)

	err = WriteSrtFile(subtitles, opts.Output)

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func linearInterpolation(subtitles []SubtitleItem, p LinearParam) {
	// https://en.wikipedia.org/wiki/Linear_interpolation
	//
	// Mapping x to y
	//
	//      y - y1        y2 - y1
	//    --------- =  ------------
	//      x - x1        x2 - x1
	//
	//  so
	//
	//                        y2 - y1
	//   y = y1 + (x - x1) * ---------
	//                        x2 - x1
	//
	// x1 = startFrom
	// x2 = endFrom
	// y1 = startTo
	// y2 = endTo

	p.coeff = float64(p.endTo-p.startTo) / float64(p.endFrom-p.startFrom)

	for i := range subtitles {
		subtitles[i].SetFromMS(p.TransformMS(subtitles[i].FromMS))
		subtitles[i].SetToMS(p.TransformMS(subtitles[i].ToMS))
	}
}

func (p LinearParam) TransformMS(ms int) int {
    return p.startTo + int(float64(ms - p.startFrom) * p.coeff)
}
