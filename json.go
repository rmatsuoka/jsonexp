package jsonexp

type (
	// Value is Object | Array | String | Number | Boolean | nil
	Value   = any
	Object  = map[string]Value
	Array   = []Value
	String  = string
	Number  = float64
	Boolean = bool
)
