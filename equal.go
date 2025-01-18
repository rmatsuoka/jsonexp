package jsonexp

func equalValue(exp valueExp, val Value) bool {
	switch exp := exp.(type) {
	case objectExp:
		obj, ok := val.(Object)
		if !ok {
			return false
		}
		return exp.match(obj)
	case arrayExp:
		arr, ok := val.(Array)
		if !ok {
			return false
		}
		return exp.match(arr)
	case *textExp:
		return exp.match(val)
	case numberExp:
		return exp.match(val)
	case booleanExp:
		return exp.match(val)
	case nil:
		return val == nil
	default:
		panic("unreachable")
	}
}
