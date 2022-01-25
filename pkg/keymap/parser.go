package keymap

import (
	"io"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

var parser = participle.MustBuild(&File{},
	participle.UseLookahead(2),
	participle.Unquote("String"),
)

type File struct {
	Pos lexer.Position

	Sections []*Section `parser:"( @@ ';'* )*"`
}

type Section struct {
	Pos lexer.Position

	Includes []*Include `parser:"@@+"`
	Configs  []*Config  `parser:"| @@"`
	Device   *Device    `parser:"| '/' '{' @@ '}'';'"`
}

type Include struct {
	Pos lexer.Position

	Value string `parser:"'#'Ident'<'@((Ident ('-' Ident)? '/'?)* ('.' Ident))'>'"`
}

type Config struct {
	Pos lexer.Position

	Behavior string `parser:"'&'@Ident '{' '}'"`
}

type Device struct {
	Pos lexer.Position

	Combos *Combos `parser:"'combos' '{' @@"`
	Keymap *Keymap `parser:"'keymap' '{' @@"`
}

type Combos struct {
	Pos lexer.Position

	Compatible string   `parser:"'compatible' '=' @String';'"`
	Combos     []*Combo `parser:"@@* '}'';'"`
}

type Combo struct {
	Pos lexer.Position

	Name         string      `parser:"@Ident '{'"`
	Timeout      int32       `parser:"'timeout''-''ms' '=' '<'@Int'>'';'"`
	KeyPositions []*List     `parser:"'key''-''positions' '=' '<'@@+'>'';'"`
	Bindings     []*Behavior `parser:"'bindings' '=' '<'@@+'>'';' '}'';'"`
}

type Keymap struct {
	Pos lexer.Position

	Compatible string   `parser:"'compatible' '=' @String';'"`
	Layers     []*Layer `parser:"@@+ '}'';'"`
}

type Layer struct {
	Pos lexer.Position

	Name     string      `parser:"@Ident '{'"`
	Bindings []*Behavior `parser:"'bindings' '=' '<'@@+'>'';' '}'';'"`
}

type List struct {
	Number  *int32  `parser:"@Int"`
	KeyCode *string `parser:"| @(Ident ('('Ident')')?)"`
}

type Behavior struct {
	Pos lexer.Position

	Action string  `parser:"'&'@Ident"`
	Params []*List `parser:"@@*"`
}

func Parse(r io.Reader) (*File, error) {
	ast := &File{}
	err := parser.Parse(r, ast)
	return ast, err
}
