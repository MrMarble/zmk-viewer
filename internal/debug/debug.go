package debug

import (
	"fmt"

	"github.com/mrmarble/zmk-viewer/pkg/keymap"
)

type EnbfCmd struct{}

func (e *EnbfCmd) Run() error {
	fmt.Print(keymap.Enbf())
	return nil
}
