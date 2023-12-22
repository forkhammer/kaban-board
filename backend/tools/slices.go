package tools

func Unique[T interface{}, K int | string | uint](slice []T, keyFunc func(T) K) []T {
	result := make([]T, 0)

	for _, el := range slice {
		i := IndexOf(result, func(s T) bool {
			return keyFunc(s) == keyFunc(el)
		})

		if i == -1 {
			result = append(result, el)
		}
	}

	return result
}

func IndexOf[T interface{}](slice []T, predicate func(T) bool) int {
	for i, el := range slice {
		if predicate(el) {
			return i
		}
	}

	return -1
}

func Filter[T interface{}](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)

	for _, el := range slice {
		if predicate(el) {
			result = append(result, el)
		}
	}

	return result
}

func Map[T interface{}, K interface{}](slice []T, mapFunc func(T) K) []K {
	result := make([]K, 0)

	for i := range slice {
		el := slice[i]
		result = append(result, mapFunc(el))
	}

	return result
}

func Find[T interface{}](slice []T, predicate func(T) bool) *T {
	for _, el := range slice {
		if predicate(el) {
			return &el
		}
	}

	return nil
}
