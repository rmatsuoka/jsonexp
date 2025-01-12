package jsonexp

type (
	// value is object | array | string | number | boolean | nil
	value  = any
	object = map[string]value
	array  = []value
	// string is just string
	number  = float64
	boolean = bool
)

var null any = nil
