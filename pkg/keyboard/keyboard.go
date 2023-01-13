package keyboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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

type localfile struct {
	Layouts map[string]Layout
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

func load(path string) (*localfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := localfile{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func LoadKeyboard(name, path string) (Layouts, error) {
	log.Debug().Str("name", name).Send()
	if path != "" {
		f, err := load(path)
		if err != nil {
			return nil, err
		}
		l := f.Layouts
		return l, nil
	}

	url := "https://keyboards.qmk.fm/v1/keyboards/%v/info.json"

	f, err := fetch(fmt.Sprintf(url, name))
	if err != nil {
		return nil, err
	}

	l := f.Keyboards[name].Layouts
	return l, nil
}
