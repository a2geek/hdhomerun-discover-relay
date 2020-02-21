package main

import (
	"hdhomerun-discover-relay/cmd"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information; use multiple times for more detail"`

	Discover cmd.DiscoverCommand `command:"discover" alias:"d" description:"Test HDHomeRun discovery mechanism"`
	Relay    cmd.RelayCommand    `command:"relay" alias:"r" description:"Relay HDHomeRun discovery packets"`
}

func main() {
	opts := &Options{}
	parser := flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)
	parser.Parse()
}
