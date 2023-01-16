/*
 * infojson - Fetches and parses the info.json file from the QMK API or a local file.
 *
 * See https://docs.qmk.fm/#/reference_info_json for more information.
 */

package infojson

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type Key struct {
	X float64  `json:"x"`
	Y float64  `json:"y"`
	W *float64 `json:"w"`
	H *float64 `json:"h"`
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

func FromName(name string) (Layouts, error) {
	log.Debug().Str("name", name).Send()
	url := "https://keyboards.qmk.fm/v1/keyboards/%v/info.json"

	f, err := fetch(fmt.Sprintf(url, name))
	if err != nil {
		return nil, err
	}

	l := f.Keyboards[name].Layouts
	return l, nil
}

func FromFile(name, path string) (Layouts, error) {
	log.Debug().Str("name", name).Str("path", path).Send()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := Keyboard{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	return f.Layouts, nil
}

func fetch(url string) (*file, error) {
	log.Info().Msg("Fetching keyboard layout.")
	log.Debug().Str("url", url).Send()

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f := file{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
