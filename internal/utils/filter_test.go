package utils

import "testing"

func TestFilter(t *testing.T) {
	source := []string{"0", "1", "2", "3", "6", "10"}
	expected := []string{"0", "6", "10"}

	result := Filter(source)

	if len(result) != len(expected) {
		t.Errorf("Expected %v, received %v", expected, result)
	}

	for i := range result {
		if expected[i] != result[i] {
			t.Errorf("Mismatched at index %v expected %v, received %v", i, expected[i], result[i])
		}
	}
}

func TestFilterEmpty(t *testing.T) {
	result := Filter(nil)

	if len(result) != 0 {
		t.Errorf("Expected empty, received %v", result)
	}
}
