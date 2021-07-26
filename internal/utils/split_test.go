package utils

import (
	. "github.com/ozoncp/ocp-contact-api/internal/models"
	"testing"
)

func makeTestContact(id uint64, text string) Contact {
	return Contact{Id: id , UserId: id, Type: 42, Text: text}
}

func makeTestData() []Contact {
	return []Contact{
	makeTestContact(1, "one"),
	makeTestContact(2,"two"),
	makeTestContact(3, "three"),
	makeTestContact(4, "four"),
	makeTestContact(5, "five"),
	makeTestContact(6,"six")}
}

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

func compareSlices(left []Contact, right []Contact) bool {
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
	source := makeTestData()
	batchLen := 2
	expectedCount := 3
	val, err := Split(source, batchLen)
	if err != nil {
		t.Errorf("Expected value, received error %v", err)
	}

	if len(val) != expectedCount {
		t.Fatalf("Expected %v slices, received %v", batchLen, val)
	}

	for i := 0; i < expectedCount; i++ {
		expectedSlice := source[i * batchLen:(i+1) * batchLen]
		if !compareSlices(expectedSlice, val[i]) {
			t.Fatalf("Expected %v, received %v", expectedSlice, val[i])
		}
	}
}

func TestSplitPartialBatches(t *testing.T) {
	source := makeTestData()
	sourceLen := len(source)
	batchLen := 4
	expectedCount := 2
	val, err := Split(source, batchLen)
	if err != nil {
		t.Errorf("Expected value, received error %v", err)
	}

	if len(val) != expectedCount {
		t.Fatalf("Expected %v slices, received %v", batchLen, val)
	}

	for i := 0; i < expectedCount; i++ {
		first := i * batchLen
		last := first + batchLen
		if last > sourceLen {
			last = sourceLen
		}
		expectedSlice := source[first:last]
		if !compareSlices(expectedSlice, val[i]) {
			t.Fatalf("Expected %v, received %v", expectedSlice, val[i])
		}
	}
}

func TestSplitUnchanges(t *testing.T) {
	source := makeTestData()
	res, _ := Split(source, 3)
	baseVal := res[0][0]
	source[0] = makeTestContact(42, "forty two")
	newVal := res[0][0]
	if baseVal != newVal {
		t.Errorf("res[0][0] is <%v>, <%v> expected", newVal, baseVal)
	}
}
