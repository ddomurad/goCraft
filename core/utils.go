package core

func IfThenElse(condition bool, a, b Any) Any {
	if condition {
		return a
	}
	return b
}

func MinOfUint(vars ...uint) (min uint) {
	min = vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return
}
