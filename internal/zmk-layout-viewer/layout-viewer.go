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

func Generate(board string, file string, trasnparent bool, output string) error {

	l, err := keyboard.Fetch(board)
	if err != nil {
		return err
	}

	if file != "" {
		log.Info().Msg("Parsing keymap file.")
		log.Debug().Str("board", board).Str("file", file).Send()
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

func drawLayer(layer *keymap.Layer, layout *keyboard.Layout, output string, transparent bool) error {
	mx := maxX(layout.Layout) + 1
	my := maxY(layout.Layout) + 1

	imageW := (mx * 40) + (mx+1)*5
	imageH := int(math.Ceil((my*40.)+(my+1.)*5.) + 20)

	log.Debug().Int("Image Width", imageW).Int("Image Height", imageH).Send()

	dc := gg.NewContext(imageW, imageH)
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
		dc.DrawString(layer.Name, 5., 15.)
	}

	for i, key := range layout.Layout {
		w := 40.0
		x := key.X*w + 5.0*key.X + 5
		y := key.Y*w + 5.0*key.Y + 20

		dc.DrawRoundedRectangle(x, y, key.W*w, w, 5.)
		dc.SetRGB(0., 0., 0.)
		dc.StrokePreserve()
		dc.SetRGB(1.0, 1.0, 1.0)
		dc.Fill()

		if layer != nil {
			drawBehavior(dc, layer.Bindings[i], x+w/2., y+w/2.)
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

func maxX(l []keyboard.Key) int {
	curr := 0
	for _, v := range l {
		if v.X > float64(curr) {
			curr = int(v.X)
		}
	}
	log.Debug().Int("Max X", curr).Send()
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
