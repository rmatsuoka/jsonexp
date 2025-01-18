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
		return equalArray(exp, arr)
	case *textExp:
		return exp.Match(val)
	case numberExp:
		return exp == val
	case booleanExp:
		return exp == val
	case nil:
		return val == nil
	default:
		panic("unreachable")
	}
}

func equalArray(exp arrayExp, arr Array) bool {
	if len(exp) != len(arr) {
		return false
	}
	for i := range exp {
		if !equalValue(exp[i], arr[i]) {
			return false
		}
	}
	return true
}
