package main

import (
	"fmt"
	"os"

	"subsync/sub"
)

func main() {
    globalOpts := sub.Parse()

    switch globalOpts.Mode {
    case "delay":
        sub.DelayMode(globalOpts)
    case "linear":
        sub.LinearMode(globalOpts)
    default:
        fmt.Printf("Unknown mode '%s'\n", globalOpts.Mode)

        os.Exit(1)
    }
}

//func main3() {
//	input := flag.String("i", "", ".srt input file")
//	flag.Parse()
//
//	subtitles, err := ReadSrtFile(*input)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	avg := 0
//	n := 0
//
//	for i := range subtitles {
//		item := subtitles[i]
//
//		delta := item.ToMS - item.FromMS
//		avg += delta
//		n++
//
//		if delta < 1500 {
//			fmt.Printf("%v %dms\n", item, delta)
//
//			before := item.FromMS
//			if i > 0 {
//				before = subtitles[i-1].ToMS
//			}
//
//			after := item.ToMS
//			if i < len(subtitles)-1 {
//				after = subtitles[i+1].FromMS
//			}
//
//			fmt.Printf("%d %d\n", item.FromMS-before, after-item.ToMS)
//
//			// We need to add this ms before and after the timeframe
//			plus := (1500 - delta) / 2
//
//			if plus > item.FromMS-before {
//				subtitles[i].FromMS = before + 1
//			} else {
//				subtitles[i].FromMS -= plus
//			}
//
//			if plus > after-item.ToMS {
//				subtitles[i].ToMS = after - 1
//			} else {
//				subtitles[i].ToMS += plus
//			}
//
//			subtitles[i].From = convertTimeToString(subtitles[i].FromMS)
//			subtitles[i].To = convertTimeToString(subtitles[i].ToMS)
//
//			fmt.Printf("%v\n", subtitles[i])
//		}
//	}
//
//	fmt.Printf("%vms\n", avg/n)
//
//	sub.WriteSrtFile(subtitles, "out.srt")
//}

// TODO check if any time will be negative
// TODO support simple delay mode like -d 1500ms -d 1,5s

