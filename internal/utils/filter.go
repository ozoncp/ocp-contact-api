package utils

func defaultList() []string {
	return []string{"1", "2", "3", "4", "5"}
}

func contain(source []string, val string) bool {
	for i := range source {
		if source[i] == val {
			return true
		}
	}
	return false
}

func Filter(source []string) []string {
	result := make([]string, 0, len(source))
	for _, val := range source {
		if !contain(defaultList(), val) {
			result = append(result, val)
		}
	}
	return result
}
