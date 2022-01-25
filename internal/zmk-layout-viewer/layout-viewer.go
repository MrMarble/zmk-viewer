package zmklayoutviewer

import (
	"fmt"
	"math"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/mrmarble/zmk-layout-viewer/pkg/keyboard"
	"github.com/mrmarble/zmk-layout-viewer/pkg/keymap"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	keySize  = 40.0
	spacer   = 5.0
	fontSize = 10.0
)

func Generate(board string, file string, trasnparent bool, output string) error {
	l, err := keyboard.Fetch(board)
	if err != nil {
		return err
	}

	if file != "" {
		err = withKeymap(file, output, trasnparent, l)
		if err != nil {
			return err
		}
	} else {
		for _, layout := range l {
			err = drawLayer(nil, &layout, output, trasnparent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func withKeymap(file string, output string, trasnparent bool, l map[string]keyboard.Layout) error {
	log.Info().Msg("Parsing keymap file.")
	r, err := os.Open(file)
	if err != nil {
		return err
	}

	ast, err := keymap.Parse(r)
	defer r.Close()

	if err != nil {
		return err
	}

	for _, layer := range ast.Sections[len(ast.Sections)-1].Device.Keymap.Layers {
		for _, layout := range l {
			err = drawLayer(layer, &layout, output, trasnparent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createContext(layout *keyboard.Layout) *gg.Context {
	mx := maxX(layout.Layout) + 1
	my := maxY(layout.Layout) + 1

	imageW := int((mx*keySize)+(mx*spacer)) + spacer
	imageH := int(math.Ceil((my*keySize)+(my+1.)*spacer) + (fontSize + spacer*2))

	log.Debug().Int("Image Width", imageW).Int("Image Height", imageH).Send()

	return gg.NewContext(imageW, imageH)
}

func drawLayer(layer *keymap.Layer, layout *keyboard.Layout, output string, transparent bool) error {
	dc := createContext(layout)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})
	dc.SetFontFace(face)

	if err != nil {
		return err
	}

	if !transparent {
		dc.SetHexColor("#bfb6af")
		dc.Clear()
	}

	if layer != nil {
		dc.SetRGB(0., 0., 0.)
		dc.DrawString(layer.Name, spacer, fontSize+spacer)
	}

	for i, key := range layout.Layout {
		x := key.X*keySize + spacer*key.X + spacer
		y := key.Y*keySize + spacer*key.Y + (fontSize + spacer*2)

		if key.H != nil {
			dc.DrawRoundedRectangle(x, y, key.W*keySize, *key.H*keySize, spacer)
		} else {
			dc.DrawRoundedRectangle(x, y, key.W*keySize, keySize, spacer)
		}

		dc.SetRGB(0., 0., 0.)
		dc.StrokePreserve()
		dc.SetRGB(1.0, 1.0, 1.0)
		dc.Fill()

		if layer != nil {
			drawBehavior(dc, layer.Bindings[i], x+keySize/2., y+keySize/2.)
		}
	}
	if layer != nil {
		log.Info().Str("output", fmt.Sprintf("%v/%v.png", output, layer.Name)).Msg("Layout generaged!")
		err = dc.SavePNG(fmt.Sprintf("%v/%v.png", output, layer.Name))
		if err != nil {
			return err
		}
	} else {
		log.Info().Str("output", fmt.Sprintf("%v/layout.png", output)).Msg("Layout generated!")

		err = dc.SavePNG(fmt.Sprintf("%v/layout.png", output))
		if err != nil {
			return err
		}
	}
	return nil
}

func drawBehavior(ctx *gg.Context, key *keymap.Behavior, x float64, y float64) {
	log.Debug().Str("Action", key.Action).Interface("Params", key.Params).Send()
	ctx.SetRGB(0., 0., 0.)
	for i, v := range key.Params {
		str := ""
		if v.KeyCode == nil {
			str = str + fmt.Sprintf("%v", *v.Number)
		} else {
			str = str + *v.KeyCode
		}

		dw, dh := ctx.MeasureString(str)
		ctx.DrawString(str, x-dw/2, y-dh/2.+float64(i)*10.)

	}
}

func maxX(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.X > curr {
			curr = v.X
		}
	}
	log.Debug().Float64("Max X", curr).Send()
	return curr
}

func maxY(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.Y > curr {
			curr = v.Y
		}
	}
	log.Debug().Float64("Max Y", curr).Send()
	return curr
}
