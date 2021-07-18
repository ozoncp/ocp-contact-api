package utils

import "testing"

func TestSplitNegativeBatchLen(t *testing.T) {
	val, err := Split(nil, -1)
	if err == nil {
		t.Errorf("Expected error, received %v", val)
	}
}

func TestSplitZeroBatchLen(t *testing.T) {
	val, err := Split(nil, 0)
	if err == nil {
		t.Errorf("Expected error, received %v", val)
	}
}

func TestSplitNilSource(t *testing.T) {
	val, err := Split(nil, 3)
	if err != nil {
		t.Errorf("Expected value, received error %v", err)
	}

	if len(val) != 0 {
		t.Errorf("Expected empty value, received %v", val)
	}
}

func compareSlices(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := 0; i < len(left); i++ {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func TestSplitFullBatches(t *testing.T) {
	source := []string{"one", "two", "three", "four", "five", "six"}
	batchLen := 2
	expectedCount := 3
	val, err := Split(source, batchLen)
	if err != nil {
		t.Errorf("Expected value, received error %v", err)
	}

	if len(val) != expectedCount {
		t.Errorf("Expected %v slices, received %v", batchLen, val)
		return
	}

	for i := 0; i < expectedCount; i++ {
		expectedSlice := source[i * batchLen:(i+1) * batchLen]
		if !compareSlices(expectedSlice, val[i]) {
			t.Errorf("Expected %v, received %v", expectedSlice, val[i])
			return
		}
	}
}

func TestSplitPartialBatches(t *testing.T) {
	source := []string{"one", "two", "three", "four", "five", "six"}
	sourceLen := len(source)
	batchLen := 4
	expectedCount := 2
	val, err := Split(source, batchLen)
	if err != nil {
		t.Errorf("Expected value, received error %v", err)
	}

	if len(val) != expectedCount {
		t.Errorf("Expected %v slices, received %v", batchLen, val)
		return
	}

	for i := 0; i < expectedCount; i++ {
		first := i * batchLen
		last := first + batchLen
		if last > sourceLen {
			last = sourceLen
		}
		expectedSlice := source[first:last]
		if !compareSlices(expectedSlice, val[i]) {
			t.Errorf("Expected %v, received %v", expectedSlice, val[i])
			return
		}
	}
}
