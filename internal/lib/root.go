package lib

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/mrmarble/zmk-layout-viewer/pkg/keyboard"
	"github.com/rs/zerolog/log"
)

type GenerateCmd struct {
	KeyboardName string `arg:"" help:"Keyboard name to fetch layout."`

	File        string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	Transparent bool   `optional:"" short:"t" help:"Use a transparent background."`
	Output      string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
}

func (g *GenerateCmd) Run() error {
	images := make(map[string]image.Image)

	keyboardInfo, err := keyboard.Fetch(g.KeyboardName)
	if err != nil {
		return err
	}

	for _, layout := range keyboardInfo {
		ctx := createContext(&layout)
		err := drawLayout(ctx, g.Transparent, layout)
		if err != nil {
			return err
		}

		base := ctx.Image()
		images[fmt.Sprintf("%s/%s.png", g.Output, g.KeyboardName)] = base

		if keymap, ok := parseKeymap(g.File); ok {
			for _, layer := range keymap.Device.Keymap.Layers {
				ctx := createContext(&layout)
				ctx.DrawImage(base, 0, 0)
				err := drawKeymap(ctx, layout, layer)
				if err != nil {
					return err
				}
				images[fmt.Sprintf("%s/%s_%s.png", g.Output, g.KeyboardName, layer.Name)] = ctx.Image()
			}
		}
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
