package jsonexp

func equalValue(exp valueExp, val Value) bool {
	switch exp := exp.(type) {
	case objectExp:
		obj, ok := val.(Object)
		if !ok {
			return false
		}
		return exp.matchValue(obj)
	case arrayExp:
		arr, ok := val.(Array)
		if !ok {
			return false
		}
		return exp.matchValue(arr)
	case *textExp:
		return exp.matchValue(val)
	case numberExp:
		return exp.matchValue(val)
	case booleanExp:
		return exp.matchValue(val)
	case nil:
		return val == nil
	default:
		panic("unreachable")
	}
}
