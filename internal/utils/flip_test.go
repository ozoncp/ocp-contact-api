package utils

import "testing"

func TestEmptyFlip(t *testing.T) {
	val := Flip(nil)
	if len(val) != 0 {
		t.Errorf("Expected empty, received %v", val)
	}
}

func TestFlip(t *testing.T) {
	length := 3
	source := make(map[int]string, length)
	source[1] = "one"
	source[2] = "two"
	source[3] = "three"

	val := Flip(source)
	if len(val) != length {
		t.Errorf("Expected length %v , received %v", length, val)
		return
	}

	for key, value := range val {
		if source[value] != key {
			t.Errorf("Expected key %v, received %v", source[value], key)
			return
		}
	}
}
