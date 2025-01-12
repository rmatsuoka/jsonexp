package jsonexp

type typ int

const (
	typeNull typ = iota
	typeObject
	typeArray
	typeString
	typeNumber
	typeBool
	typeAny
)

func isType(value any, typ typ) bool {
	switch value.(type) {
	case nil:
		return typ == typeAny || typ == typeNull
	case map[string]any:
		return typ == typeAny || typ == typeObject
	case []any:
		return typ == typeAny || typ == typeArray
	case string:
		return typ == typeAny || typ == typeString
	case float64:
		return typ == typeAny || typ == typeNumber
	case bool:
		return typ == typeAny || typ == typeBool
	default:
		panic("unexpected type")
	}
}
