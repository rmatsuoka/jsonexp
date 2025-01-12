package jsonexp

import (
	"fmt"
	"strings"
)

// textExp is a value expression encoded in string.
type textExp struct {
	orig  string
	typ   typ
	match func(any) bool
}

func (e *textExp) Match(x any) bool {
	if !isTyp(x, e.typ) {
		return false
	}
	return e.match(x)
}

func (e *textExp) String() string {
	return e.orig
}

func parseTextExp(text string) (*textExp, error) {
	switch {
	case strings.HasPrefix(text, "=="):
		return &textExp{
			orig:  text,
			typ:   typString,
			match: func(v any) bool { return v.(string) == text[2:] },
		}, nil
	case text == "___":
		return &textExp{
			orig:  text,
			typ:   typAny,
			match: func(v any) bool { return true },
		}, nil
	case strings.HasPrefix(text, "__"):
		return parseTypeExp(text)
	default:
		return &textExp{
			orig:  text,
			typ:   typString,
			match: func(v any) bool { return v.(string) == text },
		}, nil
	}
}

var typeExpMap = map[string]typ{
	"any":     typAny,
	"object":  typObject,
	"array":   typArray,
	"string":  typString,
	"number":  typNumber,
	"boolean": typBool,
}

func parseTypeExp(text string) (*textExp, error) {
	if !strings.HasPrefix(text, "__") || !strings.HasSuffix(text, "__") {
		return nil, fmt.Errorf("jsonexp: parse type expression %s: is not type expression", text)
	}
	typ, ok := typeExpMap[text[2:len(text)-2]]
	if !ok {
		return nil, fmt.Errorf("jsonexp: parse type expression %s: unknown type", text)
	}
	return &textExp{
		orig:  text,
		typ:   typ,
		match: func(a any) bool { return true },
	}, nil
}
