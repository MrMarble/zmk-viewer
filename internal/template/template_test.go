package template_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mrmarble/zmk-viewer/internal/template"
)

func Diff(t *testing.T, x interface{}, y interface{}) {
	t.Helper()

	diff := cmp.Diff(x, y)
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestParse(t *testing.T) {
	tpl, err := template.FromFile("testdata/template.yaml.golden")
	if err != nil {
		t.Fatal(err)
	}
	Diff(t, tpl.Keyboard.Mirror, false)
	Diff(t, len(tpl.Keyboard.Columns), 2)
	Diff(t, tpl.Keyboard.Columns[0].Step, 0)
}
