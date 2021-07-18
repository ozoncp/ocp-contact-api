package utils

func Flip(source map[int]string) map[string]int {
	result := make(map[string]int, len(source))

	for key, value := range source {
		result[value] = key
	}

	return result
}
