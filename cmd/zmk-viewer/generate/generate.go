package generate

import (
	"image"
	"image/png"
	"os"
	"path"
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

	Output string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
}

func (g *Cmd) Run() error {
	return generate(strings.ReplaceAll(g.KeyboardName, "/", "_"), g.LayoutFile, g.Output, g.File, g.Transparent)
}

func generate(keyboardName, layoutFile, output, keymapFile string, isTransparent bool) error {
	images := make(map[string]image.Image)

	var keyboardInfo keyboard.Layouts
	var err error
	if layoutFile != "" {
		keyboardInfo, err = keyboard.LoadFile(keyboardName, layoutFile)
	} else {
		keyboardInfo, err = keyboard.Fetch(keyboardName)
	}

	if err != nil {
		return err
	}

	for layoutName, layout := range keyboardInfo {
		layout := layout
		ctx := img.CreateContext(&layout)
		err := img.DrawLayout(ctx, isTransparent, layout)
		if err != nil {
			return err
		}

		base := ctx.Image()
		images[generateName(output, keyboardName, layoutName, "")] = base

		if keymap, ok := img.ParseKeymap(keymapFile); ok {
			for _, layer := range keymap.Device.Keymap.Layers {
				ctx := img.CreateContext(&layout)
				ctx.DrawImage(base, 0, 0)
				err := img.DrawKeymap(ctx, layout, layer)
				if err != nil {
					return err
				}
				images[generateName(output, keyboardName, layoutName, layer.Name)] = ctx.Image()
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
