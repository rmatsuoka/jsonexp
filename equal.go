package jsonexp

func equalValue(exp expValue, val value) bool {
	switch exp := exp.(type) {
	case expObject:
		obj, ok := val.(object)
		if !ok {
			return false
		}
		return exp.Match(obj)
	case expArray:
		arr, ok := val.(array)
		if !ok {
			return false
		}
		return equalArray(exp, arr)
	case *textExp:
		return exp.Match(val)
	case expNumber:
		return exp == val
	case expBoolean:
		return exp == val
	case nil:
		return val == nil
	default:
		panic("unreachable")
	}
}

func equalArray(exp expArray, arr array) bool {
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
