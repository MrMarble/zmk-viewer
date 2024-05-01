package devicetree

// #cgo CFLAGS: -std=c11 -fPIC
// #include "./tree-sitter-devicetree/src/parser.c"
import "C"

import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_devicetree())
	return sitter.NewLanguage(ptr)
}
