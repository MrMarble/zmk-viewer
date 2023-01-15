package lib

import (
	"image"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/mrmarble/zmk-viewer/pkg/keyboard"
	"github.com/rs/zerolog/log"
)

type GenerateCmd struct {
	KeyboardName string `arg:"" help:"Keyboard name to fetch layout."`

	File        string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	LayoutFile  string `optional:"" short:"l" type:"existingfile" help:"info.json file"`
	Transparent bool   `optional:"" short:"t" help:"Use a transparent background."`
	Output      string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
}

func (g *GenerateCmd) Run() error {
	images := make(map[string]image.Image)

	var keyboardInfo keyboard.Layouts
	var err error
	if g.LayoutFile != "" {
		keyboardInfo, err = keyboard.LoadFile(g.KeyboardName, g.LayoutFile)
	} else {
		keyboardInfo, err = keyboard.Fetch(g.KeyboardName)
	}

	if err != nil {
		return err
	}

	g.KeyboardName = strings.ReplaceAll(g.KeyboardName, "/", "_")

	for layoutName, layout := range keyboardInfo {
		layout := layout
		ctx := createContext(&layout)
		err := drawLayout(ctx, g.Transparent, layout)
		if err != nil {
			return err
		}

		base := ctx.Image()
		images[generateName(g.Output, g.KeyboardName, layoutName, "")] = base

		if keymap, ok := parseKeymap(g.File); ok {
			for _, layer := range keymap.Device.Keymap.Layers {
				ctx := createContext(&layout)
				ctx.DrawImage(base, 0, 0)
				err := drawKeymap(ctx, layout, layer)
				if err != nil {
					return err
				}
				images[generateName(g.Output, g.KeyboardName, layoutName, layer.Name)] = ctx.Image()
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

func generateName(output, name, layout, layer string) string {
	file := name
	if layout != "LAYOUT" {
		file += "_" + strings.ReplaceAll(layout, "LAYOUT_", "")
	}
	if layer != "" {
		file += "_" + layer
	}
	return path.Join(output, file+".png")
}
