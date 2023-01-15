package img

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/mrmarble/zmk-viewer/pkg/keyboard"
	"github.com/mrmarble/zmk-viewer/pkg/keymap"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	keySize  = 60.0
	spacer   = 5.0
	margin   = keySize / 8
	radius   = 5.0
	fontSize = 10.0
)

type Image struct {
	transparent bool
	raw         bool
	keyboard    keyboard.Keyboard
	keymap      string
}

func New(keyboard keyboard.Keyboard, options ...func(*Image)) *Image {
	i := &Image{
		keyboard: keyboard,
	}
	for _, option := range options {
		option(i)
	}
	return i
}

func WithTransparent() func(*Image) {
	return func(i *Image) {
		i.transparent = true
	}
}

func WithRaw() func(*Image) {
	return func(i *Image) {
		i.raw = true
	}
}

func WithKeymap(keymap string) func(*Image) {
	return func(i *Image) {
		i.keymap = keymap
	}
}

func (i *Image) GenerateLayouts() (map[string]image.Image, error) {
	images := make(map[string]image.Image)
	for layoutName, layout := range i.keyboard.Layouts {
		layout := layout
		ctx := createContext(&layout)
		err := drawLayout(ctx, i.transparent, layout)
		if err != nil {
			return nil, err
		}

		base := ctx.Image()
		images[generateName(i.keyboard.Name, layoutName, "")] = base

		if keymap, ok := parseKeymap(i.keymap); ok {
			for _, layer := range keymap.Device.Keymap.Layers {
				ctx := createContext(&layout)
				ctx.DrawImage(base, 0, 0)
				err := drawKeymap(ctx, layout, layer, i.raw, -1)
				if err != nil {
					return nil, err
				}
				images[generateName(i.keyboard.Name, layoutName, layer.Name)] = ctx.Image()
			}
		}
	}

	return images, nil
}

func (i *Image) GenerateSingle() (image.Image, error) {
	layers, err := i.GenerateLayouts()
	if err != nil {
		return nil, err
	}
	first := true
	var output *image.RGBA
	var rect image.Rectangle
	height := 0
	for _, layer := range layers {
		if first {
			first = false
			rect = image.Rect(0, 0, layer.Bounds().Dx(), layer.Bounds().Dy()*(len(layers)-1))
			output = image.NewRGBA(rect)
			continue
		}
		draw.Draw(output, image.Rect(0, height, layer.Bounds().Dx(), layer.Bounds().Dy()+height), layer, image.Point{0, 0}, draw.Src)
		height += layer.Bounds().Dy()
	}
	return output, nil
}

func (i *Image) GenerateUnified() (image.Image, error) {
	for _, layout := range i.keyboard.Layouts {
		layout := layout
		ctx := createContext(&layout)
		err := drawLayout(ctx, i.transparent, layout)
		if err != nil {
			return nil, err
		}

		base := ctx.Image()
		if keymap, ok := parseKeymap(i.keymap); ok {
			for layerIndex, layer := range keymap.Device.Keymap.Layers {
				err := drawKeymap(ctx, layout, layer, i.raw, layerIndex)
				if err != nil {
					return nil, err
				}
			}
		}
		return base, nil
	}
	return nil, fmt.Errorf("no layout found")
}

func generateName(name, layout, layer string) string {
	file := name
	if layout != "LAYOUT" {
		file += "_" + strings.ReplaceAll(layout, "LAYOUT_", "")
	}
	if layer != "" {
		file += "_" + layer
	}
	return file + ".png"
}

// parseKeymap returns struct from a .keymap file.
func parseKeymap(file string) (*keymap.File, bool) {
	if file == "" {
		return nil, false
	}
	log.Info().Msg("Parsing keymap file.")
	r, err := os.Open(file)
	if err != nil {
		log.Err(err).Send()
		return nil, false
	}

	ast, err := keymap.Parse(r)
	defer r.Close()

	if err != nil {
		log.Err(err).Send()
		return nil, false
	}
	return ast, true
}

// createContext from the calculated keyboard size.
func createContext(layout *keyboard.Layout) *gg.Context {
	mx := maxX(layout.Layout) + 1
	my := maxY(layout.Layout) + 1

	imageW := int((mx*keySize)+(mx*spacer)) + spacer
	imageH := int(math.Ceil((my*keySize)+(my+1.)*spacer) + (fontSize + spacer*2))

	log.Debug().Int("Image Width", imageW).Int("Image Height", imageH).Float64("Max X", mx).Float64("Max Y", my).Send()

	return gg.NewContext(imageW, imageH)
}

// drawLaout of the keyboard. Blank keys.
func drawLayout(ctx *gg.Context, transparent bool, layout keyboard.Layout) error {
	if !transparent {
		ctx.SetHexColor("#eeeeee")
		ctx.Clear()
	}

	for _, key := range layout.Layout {
		x, y := getKeyCoords(key)
		w := keySize
		h := keySize
		if key.H != nil {
			h = *key.H * keySize
		}
		if key.W != nil {
			w = *key.W * keySize
		}
		ctx.DrawRoundedRectangle(x, y, w, h, radius)
		ctx.SetRGB(0., 0., 0.)
		ctx.StrokePreserve()
		ctx.SetHexColor("#888888")
		ctx.Fill()
		ctx.DrawRoundedRectangle(x+(margin*2)/2, y+2, w-(margin*2), h-(margin*2), radius)
		ctx.SetHexColor("#a7a7a7")
		ctx.Fill()
	}
	return nil
}

// drawKeymap of the keyboard. Legend on top of the keys.
func drawKeymap(ctx *gg.Context, layout keyboard.Layout, layer *keymap.Layer, raw bool, layerNum int) error {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})
	ctx.SetFontFace(face)

	ctx.SetRGB(0., 0., 0.)
	if layerNum == -1 {
		ctx.DrawString(layer.Name, spacer, fontSize+spacer)
	}

	for i, key := range layout.Layout {
		x, y := getKeyCoords(key)
		drawBehavior(ctx, layer.Bindings[i], x+margin+3, y+margin*2.5, raw, layerNum)
	}
	return nil
}

func getKeyCoords(key keyboard.Key) (float64, float64) {
	x := key.X*keySize + spacer*key.X + spacer
	y := key.Y*keySize + spacer*key.Y + (fontSize + spacer*2)

	return x, y
}

func drawBehavior(ctx *gg.Context, key *keymap.Behavior, x float64, y float64, raw bool, layerNum int) {
	log.Debug().Str("Action", key.Action).Interface("Params", key.Params).Send()
	if key.Params == nil || len(key.Params) == 0 {
		return
	}
	leyend := key.Params[0]
	if len(key.Params) > 1 {
		leyend = key.Params[1]
	}
	str := ""
	if leyend.KeyCode == nil {
		str += fmt.Sprintf("%v", *leyend.Number)
	} else if raw {
		str += *leyend.KeyCode
	} else {
		str += keymap.GetSymbol(*leyend.KeyCode)
	}

	w, h := ctx.MeasureString(str)
	dx, dy := 0., 0.
	if layerNum == 1 {
		ctx.SetHexColor("#5ff84a")
		dx, dy = 38-w, 0
	} else if layerNum == 2 {
		ctx.SetHexColor("#f84a4a")
		dx, dy = 38-w, 28
	} else if layerNum == 3 {
		ctx.SetHexColor("#482af8")
		dx, dy = 0, 28
	}
	ctx.DrawString(str, x+dx, y-h/2.+dy)
	if len(key.Params) > 1 {
		ctx.SetHexColor("#eeeeee")
		str := ""
		if key.Params[0].KeyCode == nil {
			str += fmt.Sprintf("%v", *key.Params[0].Number)
		} else if raw {
			str += *key.Params[0].KeyCode
		} else {
			str += keymap.GetSymbol(*key.Params[0].KeyCode)
		}
		ctx.DrawString(str, x, y-h/2.+44)
		ctx.SetRGB(0., 0., 0.)
	}
}

func maxX(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.X > curr {
			curr = v.X
		}
	}
	return curr
}

func maxY(l []keyboard.Key) float64 {
	curr := 0.
	for _, v := range l {
		if v.Y > curr {
			curr = v.Y
		}
	}
	return curr
}
