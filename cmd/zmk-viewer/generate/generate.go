package generate

import (
	"image/png"
	"os"
	"strings"

	"github.com/mrmarble/zmk-viewer/internal/img"
	"github.com/mrmarble/zmk-viewer/pkg/keyboard"
	"github.com/rs/zerolog/log"
)

type Cmd struct {
	KeyboardName string `arg:"" help:"Keyboard name to fetch layout."`

	File       string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	LayoutFile string `optional:"" short:"l" type:"existingfile" help:"info.json file"`

	Transparent bool `optional:"" short:"t" help:"Use a transparent background."`
	Raw         bool `optional:"" short:"r" help:"Draw the ZMK codes instead of the key labels."`

	Output string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
}

func (g *Cmd) Run() error {
	return generate(strings.ReplaceAll(g.KeyboardName, "/", "_"), g.LayoutFile, g.Output, g.File, g.Transparent, g.Raw)
}

func generate(keyboardName, layoutFile, output, keymapFile string, isTransparent, isRaw bool) error {
	var layouts keyboard.Layouts
	var err error
	if layoutFile != "" {
		layouts, err = keyboard.LoadFile(keyboardName, layoutFile)
	} else {
		layouts, err = keyboard.Fetch(keyboardName)
	}

	if err != nil {
		return err
	}

	kbd := keyboard.Keyboard{
		Name:    keyboardName,
		Layouts: layouts,
	}

	options := []func(*img.Image){}
	if isTransparent {
		options = append(options, img.WithTransparent())
	}

	if keymapFile != "" {
		options = append(options, img.WithKeymap(keymapFile))
	}

	if isRaw {
		options = append(options, img.WithRaw())
	}

	img := img.New(kbd, options...) // TODO: add options

	images, err := img.GenerateLayouts()
	if err != nil {
		return err
	}

	for path, image := range images {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if err = png.Encode(f, image); err != nil {
			return nil
		}
		log.Info().Str("Path", path).Msg("Image saved")
	}

	return nil
}
