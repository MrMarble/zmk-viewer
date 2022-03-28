package template

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Keyboard Keyboard
}

type Keyboard struct {
	Name    string
	Columns []Column
	Keys    []Key
	Mirror  bool
}

type Column struct {
	Keys int
	Step float64
}

type Key struct {
	Column int
	Row    int
	Step   float64
	H      float64
	W      float64
}

func FromFile(name string) (*Template, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

func Parse(data []byte) (*Template, error) {
	tpl := Template{}

	err := yaml.Unmarshal(data, &tpl)
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}
