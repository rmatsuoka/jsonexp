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
	case Object:
		return typ == typeAny || typ == typeObject
	case Array:
		return typ == typeAny || typ == typeArray
	case String:
		return typ == typeAny || typ == typeString
	case Number:
		return typ == typeAny || typ == typeNumber
	case Boolean:
		return typ == typeAny || typ == typeBool
	default:
		panic("unexpected type")
	}
}
