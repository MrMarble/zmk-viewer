package devicetree_test

import (
	"testing"

	"github.com/mrmarble/zmk-viewer/pkg/devicetree"
)

func TestCanLoadGrammar(t *testing.T) {
	language := devicetree.GetLanguage()
	if language == nil {
		t.Errorf("Error loading Devicetree grammar")
	}
}
