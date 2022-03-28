package lib

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"

	"github.com/mrmarble/zmk-viewer/internal/template"
	"github.com/mrmarble/zmk-viewer/pkg/keyboard"
	"github.com/rs/zerolog/log"
)

type GenerateCmd struct {
	KeyboardName string `arg:"" help:"Keyboard name to fetch layout."`

	File        string `optional:"" short:"f" type:"existingfile" help:"ZMK .keymap file"`
	Transparent bool   `optional:"" short:"t" help:"Use a transparent background."`
	Output      string `optional:"" short:"o" type:"existingdir" default:"." help:"Output directory."`
	Template    string `optional:"" type:"existingfile" help:"Template to generate Layout"`
}

func (g *GenerateCmd) Run() error {
	images := make(map[string]image.Image)

	var keyboardInfo keyboard.Layouts
	var err error

	if g.Template != "" {
		keyboardInfo, err = fromTemplate(g.Template)
	} else {
		keyboardInfo, err = fromRemote(g.KeyboardName)
		if err != nil {
			return err
		}
	}

	for _, layout := range keyboardInfo {
		ctx := createContext(&layout)
		err := drawLayout(ctx, g.Transparent, layout)
		if err != nil {
			return err
		}

		base := ctx.Image()
		images[path.Join(g.Output, fmt.Sprintf("%s.png", g.KeyboardName))] = base

		if keymap, ok := parseKeymap(g.File); ok {
			for _, layer := range keymap.Device.Keymap.Layers {
				ctx := createContext(&layout)
				ctx.DrawImage(base, 0, 0)
				err := drawKeymap(ctx, layout, layer)
				if err != nil {
					return err
				}
				images[path.Join(g.Output, fmt.Sprintf("%s_%s.png", g.KeyboardName, layer.Name))] = ctx.Image()
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

func fromTemplate(name string) (keyboard.Layouts, error) {
	tpl, err := template.FromFile(name)
	if err != nil {
		return nil, err
	}
	l := keyboard.Layouts{}
	l[tpl.Keyboard.Name] = keyboard.FromTemplate(*tpl)
	return l, nil
}

func fromRemote(name string) (keyboard.Layouts, error) {
	keyboardInfo, err := keyboard.Fetch(name)
	if err != nil {
		return nil, err
	}
	return keyboardInfo, nil
}
