package generate

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrmarble/zmk-viewer/internal/img"
	"github.com/mrmarble/zmk-viewer/pkg/infojson"
	"github.com/rs/zerolog/log"
)

type Cmd struct {
	KeyboardName string `arg:"" help:"Keyboard name to fetch layout."`

	File       string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	LayoutFile string `optional:"" short:"l" type:"existingfile" help:"info.json file"`

	Transparent bool   `optional:"" short:"t" help:"Use a transparent background."`
	Raw         bool   `optional:"" short:"r" help:"Draw the ZMK codes instead of the key labels."`
	Single      bool   `optional:"" short:"s" help:"Generate a single image."`
	Unified     bool   `optional:"" short:"u" help:"Generate a single image with all the layers."`
	Output      string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
}

func (g *Cmd) Run() error {
	return generate(g.KeyboardName, g.LayoutFile, g.Output, g.File, g.Transparent, g.Raw, g.Single, g.Unified)
}

func generate(keyboardName, layoutFile, output, keymapFile string, isTransparent, isRaw, single, unified bool) error {
	var layouts infojson.Layouts
	var err error
	if layoutFile != "" {
		layouts, err = infojson.FromFile(keyboardName, layoutFile)
	} else {
		layouts, err = infojson.FromName(keyboardName)
	}

	if err != nil {
		return err
	}

	kbd := infojson.Keyboard{
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

	kbdImage := img.New(kbd, options...)

	var images []img.Layer
	switch {
	case single:
		outputImage, err := kbdImage.GenerateSingle()
		if err != nil {
			return err
		}
		images = []img.Layer{
			{Name: keyboardName + ".png", Image: outputImage},
		}
	case unified:
		outputImage, err := kbdImage.GenerateUnified()
		if err != nil {
			return err
		}
		images = []img.Layer{
			{Name: keyboardName + ".png", Image: outputImage},
		}
	default:
		images, err = kbdImage.GenerateLayouts(true)

	}

	if err != nil {
		return err
	}

	for _, image := range images {
		sanitized := strings.ReplaceAll(image.Name, "/", "_")
		f, err := os.Create(filepath.Join(output, sanitized))
		if err != nil {
			return err
		}
		defer f.Close()
		if err = png.Encode(f, image); err != nil {
			return nil
		}
		log.Info().Str("Path", sanitized).Msg("Image saved")
	}

	return nil
}
