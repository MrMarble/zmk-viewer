package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/zmk-layout-viewer/internal/lib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type debugFlag bool

func (d debugFlag) BeforeApply() error {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return nil
}

type Globals struct {
	Debug debugFlag `short:"D" help:"Enable debug mode"`
}

type CLI struct {
	Globals
	Generate lib.GenerateCmd `cmd:"" help:"Generate layout image."`
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cli := CLI{}

	ctx := kong.Parse(&cli,
		kong.Name("zmk-viewer"),
		kong.Description("A cli tool for visualizing zmk layouts"),
		kong.UsageOnError())

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
