package keymap

import (
	"context"
	"io"
	"strings"

	"github.com/mrmarble/zmk-viewer/pkg/devicetree"
	sitter "github.com/smacker/go-tree-sitter"
)

type Keymap struct {
	Layers []*Layer
}

type Layer struct {
	Name     string
	Bindings []Binding
}

type Binding struct {
	Action    string
	Modifiers []string
}

func parse(source []byte) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(devicetree.GetLanguage())

	return parser.ParseCtx(context.Background(), nil, source)
}

func getKeymap(tree *sitter.Tree, source []byte) (*sitter.Node, error) {
	q, _ := sitter.NewQuery([]byte(`((node (identifier) @keymap) @node (#eq? @keymap "keymap"))`), devicetree.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	var keymap *sitter.Node

	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, source)
		if m == nil || len(m.Captures) == 0 {
			continue
		}
		if m.Captures[0].Node != nil {
			keymap = m.Captures[0].Node
			break
		}
	}

	return keymap, nil
}

func getLayers(keymap *sitter.Node, source []byte) ([]*sitter.Node, error) {
	q, _ := sitter.NewQuery([]byte(`((node (identifier) @ident) @node (#not-eq? @ident "keymap"))`), devicetree.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, keymap)

	var layers []*sitter.Node

	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, source)
		if m == nil || len(m.Captures) == 0 {
			continue
		}
		if m.Captures[0].Node != nil {
			layers = append(layers, m.Captures[0].Node)
		}
	}
	return layers, nil
}

func parseLayer(layer *sitter.Node, source []byte) (*Layer, error) {
	l := &Layer{}
	l.Name = layer.Child(0).Content(source)
	for i := 0; i < int(layer.ChildCount()); i++ {
		if layer.Child(i).Type() == "property" && layer.Child(i).Content(source)[:8] == "bindings" {
			var bindings []Binding
			b := layer.Child(i).Child(2).Child(1)
			for {
				action := b.Content(source)
				if action == "&bootloader" {
					b = b.NextSibling()
					continue
				}
				var modifiers []string
				for {
					if b.NextSibling() == nil || b.NextSibling().Type() == ">" {
						break
					}
					modifier := b.NextSibling().Content(source)
					if strings.HasPrefix(modifier, "&") {
						break
					}
					modifiers = append(modifiers, modifier)
					b = b.NextSibling()
				}

				bindings = append(bindings, Binding{Action: action, Modifiers: modifiers})

				b = b.NextSibling()
				if b == nil || b.Type() == ">" {
					break
				}
			}
			l.Bindings = bindings
		}
	}
	return l, nil
}

func Parse(r io.Reader) (*Keymap, error) {
	source, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tree, err := parse(source)
	if err != nil {
		return nil, err
	}

	keymap, err := getKeymap(tree, source)
	if err != nil {
		return nil, err
	}

	layers, err := getLayers(keymap, source)
	if err != nil {
		return nil, err
	}

	parsedLayers := make([]*Layer, 0, len(layers))
	for _, layer := range layers {
		parsedLayer, err := parseLayer(layer, source)
		if err != nil {
			return nil, err
		}
		parsedLayers = append(parsedLayers, parsedLayer)
	}

	return &Keymap{Layers: parsedLayers}, nil
}
