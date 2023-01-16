package img

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	keyboard "github.com/mrmarble/zmk-viewer/pkg/infojson"
	"github.com/mrmarble/zmk-viewer/pkg/keymap"
	"github.com/rs/zerolog/log"
)

//go:embed FiraCode-Bold.ttf
var firaCode []byte

const (
	keySize  = 70.0
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
			rect = image.Rect(0, 0, layer.Bounds().Dx(), layer.Bounds().Dy()*len(layers))
			output = image.NewRGBA(rect)
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
		if !i.transparent {
			ctx.SetHexColor("#eeeeee")
			ctx.Clear()
		}
		base := ctx.Image()
		if keymap, ok := parseKeymap(i.keymap); ok {
			keys := make([]*keycap, len(layout.Layout))
			for layerIndex, layer := range keymap.Device.Keymap.Layers {
				for keyIndex, key := range layer.Bindings {
					if layerIndex == 0 {
						log.Debug().Msgf("Adding key %d", keyIndex)
						x, y := getKeyCoords(layout.Layout[keyIndex])
						keys[keyIndex] = newKeycap(x, y, keySize, keySize).fromKey(key, !i.raw)
					} else {
						log.Debug().Msgf("Updating key %d", keyIndex)
						keys[keyIndex].setLayer(layerIndex, key, !i.raw)
					}
				}
			}
			for _, key := range keys {
				key.draw(ctx)
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

	ctx := gg.NewContext(imageW, imageH)
	f, err := truetype.Parse(firaCode)
	if err != nil {
		log.Err(err).Send()
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: 12.0,
		// Hinting: font.HintingFull,
	})
	ctx.SetFontFace(face)

	return ctx
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
		newKeycap(x, y, w, h).draw(ctx)
	}
	return nil
}

// drawKeymap of the keyboard. Legend on top of the keys.
func drawKeymap(ctx *gg.Context, layout keyboard.Layout, layer *keymap.Layer, raw bool, layerNum int) error {
	ctx.SetRGB(0., 0., 0.)
	if layerNum == -1 {
		ctx.DrawString(layer.Name, spacer, fontSize+spacer)
	}

	for i, key := range layout.Layout {
		x, y := getKeyCoords(key)
		newKeycap(x, y, keySize, keySize).fromKey(layer.Bindings[i], !raw).draw(ctx)
	}
	return nil
}

func getKeyCoords(key keyboard.Key) (float64, float64) {
	x := key.X*keySize + spacer*key.X + spacer
	y := key.Y*keySize + spacer*key.Y + (fontSize + spacer*2)

	return x, y
}

type keycap struct {
	x      float64
	y      float64
	w      float64
	h      float64
	base   string
	layer1 string
	layer2 string
	layer3 string
	mod    string
}

func newKeycap(x, y, w, h float64) *keycap {
	return &keycap{x: x, y: y, w: w, h: h}
}

func (k *keycap) fromKey(key *keymap.Behavior, parseKeyCode bool) *keycap {
	if key.Params == nil || len(key.Params) == 0 {
		return k
	}

	k.base = formatKeyCode(key.Params[0], parseKeyCode)
	if len(key.Params) > 1 {
		k.base = formatKeyCode(key.Params[1], parseKeyCode)
		k.mod = formatKeyCode(key.Params[0], parseKeyCode)
	}
	if key.Action == "mo" {
		k.base = "L" + k.base
	}
	return k
}

func formatKeyCode(key *keymap.List, parseKeyCode bool) string {
	str := ""
	if key.KeyCode == nil {
		str += fmt.Sprintf("%v", *key.Number)
	} else if parseKeyCode {
		str += keymap.GetSymbol(*key.KeyCode)
	} else {
		str += *key.KeyCode
	}
	if strings.HasPrefix(str, "LC") {
		str = "⌃" + str[3:len(str)-1]
	}

	return str
}

func (k *keycap) setLayer(layer int, key *keymap.Behavior, parseKeyCode bool) {
	if key.Params == nil || len(key.Params) == 0 {
		return
	}
	switch layer {
	case 1:
		k.layer1 = formatKeyCode(key.Params[0], parseKeyCode)
		if key.Action == "mo" {
			k.layer1 = "L" + k.layer1
		}
	case 2:
		k.layer2 = formatKeyCode(key.Params[0], parseKeyCode)
		if key.Action == "mo" {
			k.layer2 = "L" + k.layer2
		}
	case 3:
		k.layer3 = formatKeyCode(key.Params[0], parseKeyCode)
		if key.Action == "mo" {
			k.layer3 = "L" + k.layer3
		}
	}
}

func (k *keycap) draw(ctx *gg.Context) {
	log.Debug().Str("Base", k.base).Str("MOD", k.mod).Str("Layer1", k.layer1).Str("Layer2", k.layer2).Str("Layer3", k.layer3).Msg("Drawing key")
	ctx.DrawRoundedRectangle(k.x, k.y, k.w, k.h, radius)

	// Border
	ctx.SetColor(color.Black)
	ctx.StrokePreserve()

	// Shadow
	ctx.SetColor(color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff})
	ctx.Fill()

	// Highlight
	ctx.DrawRoundedRectangle(k.x+margin, k.y+4, k.w-(margin*2), k.h-(margin*2), radius/2.0)
	ctx.SetColor(color.RGBA{R: 0xa7, G: 0xa7, B: 0xa7, A: 0xff})
	ctx.Fill()

	if k.base != "" {
		ctx.SetColor(color.Black)
		_, sh := ctx.MeasureString(k.base)
		ctx.DrawString(k.base, k.x+margin+3, k.y+margin+sh)
	}

	if k.mod != "" {
		sw, _ := ctx.MeasureString(k.mod)
		ctx.SetColor(color.RGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff})
		ctx.DrawString(k.mod, k.x+(k.w/2)-sw/2, k.y+k.h-3)
	}

	if k.layer1 != "" {
		ctx.SetColor(color.RGBA{R: 0x05, G: 0x90, B: 0x33, A: 0xff})
		sw, sh := ctx.MeasureString(k.layer1)
		ctx.DrawString(k.layer1, k.x+k.w-sw-margin-3, k.y+margin+sh)
	}

	if k.layer2 != "" {
		ctx.SetColor(color.RGBA{R: 0x05, G: 0x06, B: 0xb1, A: 0xff})
		sw, sh := ctx.MeasureString(k.layer2)
		ctx.DrawString(k.layer2, k.x+k.w-sw-margin-3, k.y+k.h-sh-margin)
	}

	if k.layer3 != "" {
		ctx.SetColor(color.RGBA{R: 0xe3, G: 0x00, B: 0x52, A: 0xff})
		_, sh := ctx.MeasureString(k.layer3)
		ctx.DrawString(k.layer3, k.x+margin+3, k.y+k.h-sh-margin)
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
