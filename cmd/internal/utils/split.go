package utils

import "errors"

func Split(source []string, batchLen int) ([][]string, error) {
	if batchLen <= 0 {
		return nil, errors.New("batch length cannot be a zero or negative")
	}
	sourceLen := len(source)
	fullBatchesCount := sourceLen / batchLen
	batchesCount := fullBatchesCount
	if sourceLen % batchLen != 0 {
		batchesCount++
	}
	result := make([][]string, batchesCount)

	for i := 0; i < fullBatchesCount; i++ {
		result[i] = source[i * batchLen:(i + 1) * batchLen]
	}

	if fullBatchesCount != batchesCount {
		result[fullBatchesCount] = source[fullBatchesCount * batchLen:sourceLen]
	}

	return result, nil
}