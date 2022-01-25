package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"
	zmklayoutviewer "github.com/mrmarble/zmk-layout-viewer/internal/zmk-layout-viewer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type debugFlag bool

func (d debugFlag) BeforeApply() error {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return nil
}

var cli struct {
	Keyboard string `arg:"" help:"Keyboard name to fetch layout."`

	File        string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	Transparent bool   `optional:"" short:"t" help:"Use a transparent background."`
	Output      string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`

	Debug debugFlag `help:"Enable debug logging."`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx := kong.Parse(&cli)

	err := zmklayoutviewer.Generate(cli.Keyboard, cli.File, cli.Transparent, cli.Output)
	ctx.FatalIfErrorf(err)
}
