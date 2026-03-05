package common

func ParseInt(value any) (int, bool) {
	switch n := value.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

func ParseInt64(value any) (int64, bool) {
	switch n := value.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	default:
		return 0, false
	}
}

func ParseFloat64(value any) (float64, bool) {
	switch n := value.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case float32:
		return float64(n), true
	default:
		return 0, false
	}
}

func ParseStringSlice(value any) ([]string, bool) {
	arr, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		result = append(result, s)
	}
	return result, true
}

func ParseFloat64Slice(value any) ([]float64, bool) {
	arr, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]float64, 0, len(arr))
	for _, item := range arr {
		n, ok := ParseFloat64(item)
		if !ok {
			return nil, false
		}
		result = append(result, n)
	}
	return result, true
}

func ParseIntSlice(value any) ([]int, bool) {
	arr, ok := value.([]any)
	if !ok {
		return nil, false
	}
	result := make([]int, 0, len(arr))
	for _, item := range arr {
		n, ok := ParseInt(item)
		if !ok {
			return nil, false
		}
		result = append(result, n)
	}
	return result, true
}
