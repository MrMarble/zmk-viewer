package keymap

import "strings"

func GetSymbol(key string) string {
	if strings.HasPrefix(key, "NUM_") {
		return key[4:]
	}

	switch key {
	case "LSFT", "RSFT":
		return "Shift" // "⇧"
	case "LCTL", "RCTL":
		return "Ctrl" // "⌃"
	case "LALT", "RALT":
		return "Alt" // "⌥"
	case "SPC":
		return "Spc" // "␣"
	case "RET", "RETURN", "ENTER":
		return "Enter" //"⏎"
	case "CMMA":
		return ","
	case "DOT":
		return "."
	case "QUOT":
		return "'"
	case "TAB":
		return "⇥"
	case "ESC":
		return "⎋"
	case "DEL":
		return "⌦"
	case "RBKT":
		return "]"
	case "LBKT":
		return "["
	case "RBRC":
		return "}"
	case "LBRC":
		return "{"
	case "EQL", "KP_EQUAL":
		return "="
	case "MINUS":
		return "-"
	case "BACKSLASH", "BSLH", "NON_US_BACKSLASH", "NON_US_BSLH", "NUBS":
		return "\\"
	case "HASH", "NON_US_HASH", "NUHS", "POUND":
		return "#"
	case "TILDE", "TILDE2":
		return "~"
	case "PIPE", "PIPE2":
		return "|"
	case "DQT", "DOUBLE_QUOTES":
		return "\""
	case "GRAVE":
		return "`"
	case "LGUI", "RGUI":
		return "⌘"
	case "UP":
		return "↑"
	case "DOWN":
		return "↓"
	case "LEFT":
		return "←"
	case "RIGHT":
		return "→"
	case "BKSP":
		return "⌫"
	case "UNDER":
		return "_"
	case "CARET", "CRRT":
		return "^"
	case "AMPS":
		return "&"
	case "KMLT":
		return "*"
	case "LPRN":
		return "("
	case "RPRN":
		return ")"
	case "COLN":
		return ":"
	case "SCLN":
		return ";"
	case "DLLR":
		return "$"
	case "PRCT":
		return "%"
	case "ATSN":
		return "@"
	case "BANG":
		return "!"
	case "QMARK":
		return "?"
	case "FSLH":
		return "/"
	case "KPLS":
		return "+"
	case "PG_UP":
		return "PgUp"
	case "PG_DN":
		return "PgDn"
	case "K_VOL_UP":
		return "VolUp" // "🔊"
	case "K_VOL_DN":
		return "VolDn" // "🔉"
	case "K_MUTE":
		return "Mute" // "🔇"
	case "BT_NXT":
		return "BT+" // "⏭"
	case "BT_PRV":
		return "BT-" // "⏮"
	case "BT_CLR":
		return "BTClr" // "⏯"
	}
	return key
}
