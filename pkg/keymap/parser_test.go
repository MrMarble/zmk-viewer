package keymap

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func getRoot(t *testing.T) string {
	_, b, _, _ := runtime.Caller(0)

	// Root folder of this project
	return filepath.Join(filepath.Dir(b), "..", "..")
}

func getTestdata(t *testing.T) string {
	t.Helper()
	return filepath.Join(getRoot(t), "testdata")
}

func getParams(t *testing.T, name string) (int, int) {
	t.Helper()
	data := strings.Split(name, "_")
	if len(data) != 3 {
		t.Fatalf("Expected file %s to have 3 parts", name)
	}

	out1, err := strconv.Atoi(data[1])
	if err != nil {
		t.Fatalf("Expected %s to be an integer", data[1])
	}

	out2, err := strconv.Atoi(data[2])
	if err != nil {
		t.Fatalf("Expected %s to be an integer", data[2])
	}

	return out1, out2
}

func getKeyamps(t *testing.T) []fs.DirEntry {
	t.Helper()
	keymaps, err := os.ReadDir(filepath.Join(getTestdata(t), "keymaps"))
	if err != nil {
		t.Fatal(err)
	}

	for i, k := range keymaps {
		if k.IsDir() || strings.HasSuffix(k.Name(), ".skip") == true {
			keymaps = append(keymaps[:i], keymaps[i+1:]...)
		}
	}

	return keymaps
}

func TestMain(m *testing.M) {
	v := m.Run()
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
	os.Exit(v)
}

func TestDevicetreeParse(t *testing.T) {
	keymaps := getKeyamps(t)
	for _, k := range keymaps {
		file, err := os.Open(filepath.Join(getTestdata(t), "keymaps", k.Name()))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		keymap, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(k.Name(), func(t *testing.T) {
			t.Parallel()
			tree, err := parse(keymap)
			if err != nil {
				t.Fatal(err)
			}

			snaps.WithConfig(snaps.Filename("devicetree"), snaps.Dir(filepath.Join(getTestdata(t), "__snapshots__"))).MatchSnapshot(t, tree.RootNode().String())
		})
	}
}

func TestParse(t *testing.T) {
	keymaps := getKeyamps(t)
	for _, k := range keymaps {
		n := k.Name()
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			file, err := os.Open(filepath.Join(getTestdata(t), "keymaps", n))
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			keymap, err := Parse(file)
			if err != nil {
				t.Fatal(err)
			}

			if keymap == nil {
				t.Fatal("Expected keymap to be not nil")
			}

			snaps.WithConfig(snaps.Filename("parsed"), snaps.Dir(filepath.Join(getTestdata(t), "__snapshots__"))).MatchSnapshot(t, keymap)
		})
	}
}
