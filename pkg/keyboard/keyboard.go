package keyboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Key struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	W     float64 `json:"w"`
	Label string  `json:"label"`
}

type Layout struct {
	Layout []Key `json:"layout"`
}

type Keyboard struct {
	Name    string            `json:"keyboard_name"`
	Layouts map[string]Layout `json:"layouts"`
}

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

func Fetch(name string) (map[string]Layout, error) {
	log.Debug().Str("name", name).Send()
	url := "https://keyboards.qmk.fm/v1/keyboards/%v/info.json"

	f, err := fetch(fmt.Sprintf(url, name))
	if err != nil {
		return nil, err
	}

	l := f.Keyboards[name].Layouts
	return l, nil
}
