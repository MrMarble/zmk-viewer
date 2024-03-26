package keymap

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alecthomas/participle/v2"
)

func getRoot(t *testing.T) string {
	t.Helper()
	_, b, _, _ := runtime.Caller(0)

	// Root folder of this project
	return filepath.Join(filepath.Dir(b), "..", "..")
}

func TestParse(t *testing.T) {
	testdata := filepath.Join(getRoot(t), "testdata", "keymaps")

	keymaps, err := os.ReadDir(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for _, keymap := range keymaps {
		t.Run(keymap.Name(), func(t *testing.T) {
			file, err := os.Open(filepath.Join(testdata, keymap.Name()))
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			_, err = Parse(file)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGrammar(t *testing.T) {
	gr := "#define HRML(k1,k2,k3,k4) &ht LSHFT k1  &ht LALT k2  &ht LCTRL k3  &ht LGUI k4"

	parser := participle.MustBuild[Define](participle.UseLookahead(2),
		participle.Unquote("String"))

	_, err := parser.ParseString("", gr)
	if err != nil {
		t.Fatal(err)
	}
}
