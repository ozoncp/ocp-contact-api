package utils

import (
	"errors"
	. "github.com/ozoncp/ocp-contact-api/internal/models"
)

func Split(source []Contact, batchLen int) ([][]Contact, error) {
	if batchLen <= 0 {
		return nil, errors.New("batch length cannot be a zero or negative")
	}
	sourceLen := len(source)
	fullBatchesCount := sourceLen / batchLen
	batchesCount := fullBatchesCount
	if sourceLen % batchLen != 0 {
		batchesCount++
	}
	result := make([][]Contact, batchesCount)

	for i := 0; i < fullBatchesCount; i++ {
		result[i] = append([]Contact{}, source[i * batchLen:(i + 1) * batchLen]...)
	}

	if fullBatchesCount != batchesCount {
		result[fullBatchesCount] = append([]Contact{}, source[fullBatchesCount * batchLen:sourceLen]...)
	}

	return result, nil
}