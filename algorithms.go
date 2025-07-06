package pqueue

import (
	"math"
	"reflect"
)

// insertionSort performs insertion sort on the queue data
func (pq *PQueue[T]) insertionSort() {
	for i := 1; i < pq.size; i++ {
		key := pq.data[i]
		j := i - 1

		for j >= 0 && pq.less(key, pq.data[j]) {
			pq.data[j+1] = pq.data[j]
			j--
		}
		pq.data[j+1] = key
	}
}

// quickSort performs quicksort on the queue data
func (pq *PQueue[T]) quickSort() {
	pq.quickSortRange(0, pq.size-1)
}

func (pq *PQueue[T]) quickSortRange(low, high int) {
	if low < high {
		pi := pq.partition(low, high)
		pq.quickSortRange(low, pi-1)
		pq.quickSortRange(pi+1, high)
	}
}

func (pq *PQueue[T]) partition(low, high int) int {
	pivot := pq.data[high]
	i := low - 1

	for j := low; j < high; j++ {
		if pq.less(pq.data[j], pivot) || (!pq.less(pivot, pq.data[j]) && !pq.less(pq.data[j], pivot)) {
			i++
			pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
		}
	}
	pq.data[i+1], pq.data[high] = pq.data[high], pq.data[i+1]
	return i + 1
}

// mergeSort performs merge sort on the queue data
func (pq *PQueue[T]) mergeSort() {
	if pq.size <= 1 {
		return
	}
	temp := make([]T, pq.size)
	pq.mergeSortRange(0, pq.size-1, temp)
}

func (pq *PQueue[T]) mergeSortRange(left, right int, temp []T) {
	if left < right {
		mid := left + (right-left)/2
		pq.mergeSortRange(left, mid, temp)
		pq.mergeSortRange(mid+1, right, temp)
		pq.merge(left, mid, right, temp)
	}
}

func (pq *PQueue[T]) merge(left, mid, right int, temp []T) {
	// Copy data to temp array
	for i := left; i <= right; i++ {
		temp[i] = pq.data[i]
	}

	i, j, k := left, mid+1, left

	// Merge the two halves
	for i <= mid && j <= right {
		if pq.less(temp[i], temp[j]) {
			pq.data[k] = temp[i]
			i++
		} else if pq.less(temp[j], temp[i]) {
			pq.data[k] = temp[j]
			j++
		} else {
			// Equal elements - take from left array to maintain stability
			pq.data[k] = temp[i]
			i++
		}
		k++
	}

	// Copy remaining elements
	for i <= mid {
		pq.data[k] = temp[i]
		i++
		k++
	}

	for j <= right {
		pq.data[k] = temp[j]
		j++
		k++
	}
}

// introsort performs introspective sort (hybrid of quicksort, heapsort, and insertion sort)
func (pq *PQueue[T]) introsort() {
	maxDepth := int(math.Log2(float64(pq.size))) * 2
	pq.introsortRange(0, pq.size-1, maxDepth)
}

func (pq *PQueue[T]) introsortRange(low, high, depth int) {
	size := high - low + 1

	if size <= 16 {
		pq.insertionSortRange(low, high)
		return
	}

	if depth == 0 {
		pq.heapSortRange(low, high)
		return
	}

	pi := pq.partition(low, high)
	pq.introsortRange(low, pi-1, depth-1)
	pq.introsortRange(pi+1, high, depth-1)
}

func (pq *PQueue[T]) insertionSortRange(low, high int) {
	for i := low + 1; i <= high; i++ {
		key := pq.data[i]
		j := i - 1

		for j >= low && pq.less(key, pq.data[j]) {
			pq.data[j+1] = pq.data[j]
			j--
		}
		pq.data[j+1] = key
	}
}

func (pq *PQueue[T]) heapSortRange(low, high int) {
	size := high - low + 1

	// Build heap
	for i := size/2 - 1; i >= 0; i-- {
		pq.heapify(low, size, low+i)
	}

	// Extract elements from heap
	for i := size - 1; i > 0; i-- {
		pq.data[low], pq.data[low+i] = pq.data[low+i], pq.data[low]
		pq.heapify(low, i, low)
	}
}

func (pq *PQueue[T]) heapify(base, size, root int) {
	largest := root
	left := 2*(root-base) + 1 + base
	right := 2*(root-base) + 2 + base

	if left < base+size && pq.less(pq.data[largest], pq.data[left]) {
		largest = left
	}

	if right < base+size && pq.less(pq.data[largest], pq.data[right]) {
		largest = right
	}

	if largest != root {
		pq.data[root], pq.data[largest] = pq.data[largest], pq.data[root]
		pq.heapify(base, size, largest)
	}
}

// timsort performs a simplified version of Timsort
func (pq *PQueue[T]) timsort() {
	minMerge := 32

	if pq.size <= minMerge {
		pq.insertionSort()
		return
	}

	// Find runs and merge them
	runs := pq.findRuns()
	pq.mergeRuns(runs)
}

func (pq *PQueue[T]) findRuns() []int {
	runs := []int{0}
	i := 0

	for i < pq.size-1 {
		start := i

		// Find ascending or descending run
		if pq.less(pq.data[i], pq.data[i+1]) {
			// Ascending run
			for i < pq.size-1 && pq.less(pq.data[i], pq.data[i+1]) {
				i++
			}
		} else {
			// Descending run - reverse it
			for i < pq.size-1 && (pq.less(pq.data[i+1], pq.data[i]) || (!pq.less(pq.data[i], pq.data[i+1]) && !pq.less(pq.data[i+1], pq.data[i]))) {
				i++
			}
			pq.reverse(start, i)
		}

		i++
		runs = append(runs, i)
	}

	if runs[len(runs)-1] != pq.size {
		runs = append(runs, pq.size)
	}

	return runs
}

func (pq *PQueue[T]) mergeRuns(runs []int) {
	temp := make([]T, pq.size)

	for len(runs) > 2 {
		newRuns := []int{runs[0]}

		for i := 1; i < len(runs)-1; i += 2 {
			left := runs[i-1]
			mid := runs[i] - 1
			right := runs[i+1] - 1

			pq.merge(left, mid, right, temp)
			newRuns = append(newRuns, runs[i+1])
		}

		// Handle odd number of runs
		if len(runs)%2 == 0 {
			newRuns = append(newRuns, runs[len(runs)-1])
		}

		runs = newRuns
	}
}

func (pq *PQueue[T]) reverse(start, end int) {
	for start < end {
		pq.data[start], pq.data[end] = pq.data[end], pq.data[start]
		start++
		end--
	}
}

// radixSort performs radix sort for integer types
func (pq *PQueue[T]) radixSort() {
	// This is a simplified implementation that works with reflect
	// In practice, you'd want type-specific implementations for better performance
	if pq.dataType != IntegerType {
		pq.quickSort()
		return
	}

	// Get the maximum value to determine number of digits
	maxVal := pq.getMaxInt()
	if maxVal <= 0 {
		return
	}

	// Do counting sort for every digit
	for exp := 1; maxVal/exp > 0; exp *= 10 {
		pq.countingSortByDigit(exp)
	}
}

func (pq *PQueue[T]) getMaxInt() int {
	if pq.size == 0 {
		return 0
	}

	max := 0
	for i := 0; i < pq.size; i++ {
		val := reflect.ValueOf(pq.data[i])
		if val.Kind() == reflect.Int || val.Kind() == reflect.Int64 || val.Kind() == reflect.Int32 {
			intVal := int(val.Int())
			if intVal > max {
				max = intVal
			}
		}
	}
	return max
}

func (pq *PQueue[T]) countingSortByDigit(exp int) {
	output := make([]T, pq.size)
	count := make([]int, 10)

	// Count occurrences of each digit
	for i := 0; i < pq.size; i++ {
		val := reflect.ValueOf(pq.data[i])
		digit := (int(val.Int()) / exp) % 10
		count[digit]++
	}

	// Change count[i] to actual position
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// Build output array
	for i := pq.size - 1; i >= 0; i-- {
		val := reflect.ValueOf(pq.data[i])
		digit := (int(val.Int()) / exp) % 10
		output[count[digit]-1] = pq.data[i]
		count[digit]--
	}

	// Copy output array to data
	copy(pq.data[:pq.size], output)
}

// countingSort performs counting sort for small integer ranges
func (pq *PQueue[T]) countingSort() {
	if pq.dataType != IntegerType {
		pq.quickSort()
		return
	}

	minVal, maxVal := pq.getMinMaxInt()
	if maxVal-minVal > 10000 { // Don't use counting sort for large ranges
		pq.quickSort()
		return
	}

	count := make([]int, maxVal-minVal+1)

	// Count each element
	for i := 0; i < pq.size; i++ {
		val := reflect.ValueOf(pq.data[i])
		index := int(val.Int()) - minVal
		count[index]++
	}

	// Reconstruct the array
	pos := 0
	for i := 0; i < len(count); i++ {
		for count[i] > 0 {
			val := reflect.ValueOf(minVal + i)
			pq.data[pos] = val.Interface().(T)
			pos++
			count[i]--
		}
	}
}

func (pq *PQueue[T]) getMinMaxInt() (int, int) {
	if pq.size == 0 {
		return 0, 0
	}

	val := reflect.ValueOf(pq.data[0])
	min, max := int(val.Int()), int(val.Int())

	for i := 1; i < pq.size; i++ {
		val := reflect.ValueOf(pq.data[i])
		intVal := int(val.Int())
		if intVal < min {
			min = intVal
		}
		if intVal > max {
			max = intVal
		}
	}

	return min, max
}

// isNearlySorted checks if the data is nearly sorted
func (pq *PQueue[T]) isNearlySorted() bool {
	if pq.size <= 1 {
		return true
	}

	inversions := 0
	threshold := pq.size / 10 // Allow up to 10% inversions
	if threshold < 1 {
		threshold = 1 // Allow at least 1 inversion for small arrays
	}

	for i := 0; i < pq.size-1; i++ {
		if pq.less(pq.data[i+1], pq.data[i]) {
			inversions++
			if inversions > threshold {
				return false
			}
		}
	}

	return true
}

// hasSmallRange checks if integer data has a small range
func (pq *PQueue[T]) hasSmallRange() bool {
	if pq.dataType != IntegerType || pq.size == 0 {
		return false
	}

	min, max := pq.getMinMaxInt()
	return (max - min) <= 1000 // Consider small if range is <= 1000
}
