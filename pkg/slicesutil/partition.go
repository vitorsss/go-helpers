package slicesutil

import "math"

func Partition[T any](values []T, size int) [][]T {
	valuesSize := len(values)
	partitions := int(math.Ceil(float64(valuesSize) / float64(size)))

	result := make([][]T, 0, partitions)
	for i := range partitions {
		startIdx := i * size
		endIdx := min((i+1)*size, valuesSize)
		result = append(result, values[startIdx:endIdx])
	}

	return result
}
