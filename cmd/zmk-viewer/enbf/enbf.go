package enbf

import (
	"fmt"

	"github.com/mrmarble/zmk-viewer/pkg/keymap"
)

type Cmd struct{}

func (e *Cmd) Run() error {
	fmt.Print(keymap.Enbf())
	return nil
}
