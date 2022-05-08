package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mrmarble/zmk-viewer/internal/debug"
	"github.com/mrmarble/zmk-viewer/internal/lib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	debugFlag   bool
	VersionFlag string
)

var (
	// Populated by goreleaser during build.
	version = "master"
	commit  = "?"
	date    = ""
)

func (d debugFlag) BeforeApply() error {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return nil
}

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Printf("zmk-viewer has version %s built from %s on %s\n", version, commit, date)
	app.Exit(0)

	return nil
}

type Globals struct {
	Debug debugFlag `short:"D" help:"Enable debug mode"`
}

type CLI struct {
	Globals
	Version  VersionFlag     `name:"version" help:"Print version information and quit"`
	Generate lib.GenerateCmd `cmd:"" help:"Generate layout image."`
	Enbf     debug.EnbfCmd   `cmd:"" help:"Print ENBF from parser."`
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
