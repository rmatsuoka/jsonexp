package jsonexp

import (
	"encoding/json"
	"fmt"
	"strings"
)

// textExp is a value expression encoded in string.
type textExp struct {
	orig       string
	typ        typ
	matchOnTyp func(any) bool
}

func (*textExp) valueExp() {}

func (e *textExp) match(x any) bool {
	if !isType(x, e.typ) {
		return false
	}
	return e.matchOnTyp(x)
}

func (e *textExp) String() string {
	return e.orig
}

func (e *textExp) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.orig)
}

func parseTextExp(text string) (*textExp, error) {
	switch {
	case strings.HasPrefix(text, "=="):
		return &textExp{
			orig:       text,
			typ:        typeString,
			matchOnTyp: func(v any) bool { return v.(string) == text[2:] },
		}, nil
	case strings.HasPrefix(text, "((") && strings.HasSuffix(text, "))"):
		return parseTextTypeExp(text)
	default:
		return &textExp{
			orig:       text,
			typ:        typeString,
			matchOnTyp: func(v any) bool { return v.(string) == text },
		}, nil
	}
}

var typeExpMap = map[string]typ{
	"any":     typeAny,
	"object":  typeObject,
	"array":   typeArray,
	"string":  typeString,
	"number":  typeNumber,
	"boolean": typeBool,
}

func parseTextTypeExp(text string) (*textExp, error) {
	if !strings.HasPrefix(text, "((") || !strings.HasSuffix(text, "))") {
		return nil, fmt.Errorf("jsonexp: parse type expression %s: is not type expression", text)
	}
	typ, ok := typeExpMap[text[2:len(text)-2]]
	if !ok {
		return nil, fmt.Errorf("jsonexp: parse type expression %s: unknown type", text)
	}
	return &textExp{
		orig:       text,
		typ:        typ,
		matchOnTyp: func(a any) bool { return true },
	}, nil
}
