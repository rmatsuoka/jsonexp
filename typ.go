package jsonexp

type typ int

const (
	typNull typ = iota
	typObject
	typArray
	typString
	typNumber
	typBool
	typAny
)

func isTyp(value any, typ typ) bool {
	switch value.(type) {
	case nil:
		return typ == typAny || typ == typNull
	case map[string]any:
		return typ == typAny || typ == typObject
	case []any:
		return typ == typAny || typ == typArray
	case string:
		return typ == typAny || typ == typString
	case float64:
		return typ == typAny || typ == typNumber
	case bool:
		return typ == typAny || typ == typBool
	default:
		panic("unexpected type")
	}
}
