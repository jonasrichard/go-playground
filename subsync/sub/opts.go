package sub

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type GlobalOpts struct {
	Mode   string     `short:"m" long:"mode" description:"Processing mode" choice:"delay" choice:"linear" required:"true"`
	Input  string     `short:"i" long:"input" description:"Input .srt file" required:"true" value-name:"FILE"`
	Output string     `short:"o" long:"output" description:"Output .srt file" required:"true" value-name:"FILE"`
	Delay  DelayOpts  `group:"delay" required:"false"`
	Linear LinearOpts `group:"linear" required:"false"`
}

type DelayOpts struct {
	Time string `short:"t" long:"time" description:"Source dest mapping hh:mm:ss,msc->hh:mm:ss,msc form"`
}

type LinearOpts struct {
	StartTime string `short:"s" long:"start" description:"Source dest mapping of starting point in hh:mm:ss,msc->hh:mm:ss,msc form"`
	EndTime   string `short:"e" long:"end" description:"Source dest mapping of end point in hh:mm:ss,msc->hh:mm:ss,msc form"`
}

func Parse() GlobalOpts {
	var globalOpts GlobalOpts

    parser := flags.NewParser(&globalOpts, flags.Default)

    _, err := parser.Parse()

	if err != nil {
        parser.WriteHelp(os.Stdout)

		os.Exit(1)
	}

	return globalOpts
}
