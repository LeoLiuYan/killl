package main

import (
	"context"
	"github.com/ogier/pflag"
	"killl/app"
	"killl/lib/log"
	"os"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	cfg := app.NewConfig()
	fs := pflag.NewFlagSet("killl", pflag.ExitOnError)
	cfg.AddFlagSet(fs)
	fs.Parse(os.Args[1:])

	log.InitLogger(cfg.Verbose)
	go app.Run(ctx, cfg)
	InitSignal()
}
