package keyboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mrmarble/zmk-viewer/internal/template"
	"github.com/rs/zerolog/log"
)

type Key struct {
	X     float64  `json:"x"`
	Y     float64  `json:"y"`
	W     float64  `json:"w"`
	H     *float64 `json:"h"`
	Label string   `json:"label"`
}

type Layout struct {
	Layout []Key `json:"layout"`
}

type (
	Layouts  map[string]Layout
	Keyboard struct {
		Name    string  `json:"keyboard_name"`
		Layouts Layouts `json:"layouts"`
	}
)

type file struct {
	Keyboards map[string]Keyboard `json:"keyboards"`
}

func fetch(url string) (*file, error) {
	log.Info().Msg("Fetching keyboard layout.")
	log.Debug().Str("url", url).Send()

	client := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "zmk-layout-viewer")

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	f := file{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func Fetch(name string) (Layouts, error) {
	log.Debug().Str("name", name).Send()
	url := "https://keyboards.qmk.fm/v1/keyboards/%v/info.json"

	f, err := fetch(fmt.Sprintf(url, name))
	if err != nil {
		return nil, err
	}

	l := f.Keyboards[name].Layouts
	return l, nil
}

func FromTemplate(tpl template.Template) Layout {
	layout := Layout{}
	lastColumn := 0

	generate := func(offset, last int) {
		for x, column := range tpl.Keyboard.Columns {
			for y := 0; y < column.Keys; y++ {
				var h float64 = 1
				layout.Layout = append(layout.Layout, Key{X: float64(x + offset), Y: float64(y) + column.Step, W: 1, H: &h})
			}
			lastColumn++
		}
		for _, key := range tpl.Keyboard.Keys {
			w := key.W
			h := key.H

			if w == 0 { // TODO: Implement custom Unmarshaller to handle this
				w = 1
			}
			if h == 0 {
				h = 1
			}
			layout.Layout = append(layout.Layout, Key{X: float64(key.Column + offset - last), Y: float64(key.Row) + key.Step, H: &h, W: w})
			if key.Column > lastColumn {
				lastColumn = key.Column
			}
		}
	}
	generate(0, 0)
	if tpl.Keyboard.Mirror {
		tpl.Keyboard.Columns = reverse(tpl.Keyboard.Columns)
		generate(lastColumn+3, lastColumn+1)
	}
	return layout
}

func reverse(columns []template.Column) []template.Column {
	rev := []template.Column{}
	for i := range columns {
		// reverse the order
		rev = append(rev, columns[len(columns)-1-i])
	}
	return rev
}
