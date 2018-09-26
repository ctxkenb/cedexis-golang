package cedexis

func stringsDiffer(a *string, b *string) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	return *a != *b
}

func boolsDiffer(a *bool, b *bool) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	return *a != *b
}

func intsDiffer(a *int, b *int) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	return *a != *b
}

func float64sDiffer(a *float64, b *float64) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	return *a != *b
}

func stringArraysDiffer(a *[]string, b *[]string) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	if len(*a) != len(*b) {
		return true
	}

	for i := range *a {
		if (*a)[i] != (*b)[i] {
			return true
		}
	}

	return false
}

func intArraysDiffer(a *[]int, b *[]int) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return true
	}

	if len(*a) != len(*b) {
		return true
	}

	for i := range *a {
		if (*a)[i] != (*b)[i] {
			return true
		}
	}

	return false
}
