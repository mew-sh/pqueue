// Package pqueue provides an intelligent priority queue and sorting library
// that automatically selects the best algorithm based on data characteristics.
package pqueue

import (
	"fmt"
	"reflect"
)

// PQueue represents an intelligent priority queue with adaptive sorting
type PQueue[T any] struct {
	data     []T
	less     func(T, T) bool
	dataType DataType
	size     int
}

// DataType represents the type of data being sorted
type DataType int

const (
	IntegerType DataType = iota
	FloatType
	StringType
	SliceType
	ArrayType
	StructType
	MapType
	PointerType
	InterfaceType
	ChannelType
	FuncType
	GenericType
)

// SortStrategy represents the sorting algorithm to use
type SortStrategy int

const (
	AutoStrategy SortStrategy = iota
	RadixStrategy
	CountingStrategy
	InsertionStrategy
	TimsortStrategy
	IntrosortStrategy
	MergeStrategy
	QuickStrategy
)

// New creates a new PQueue with the given data and comparison function
func New[T any](data []T, less func(T, T) bool) *PQueue[T] {
	pq := &PQueue[T]{
		data: make([]T, len(data)),
		less: less,
		size: len(data),
	}
	copy(pq.data, data)
	pq.dataType = inferDataType(data)
	return pq
}

// NewInts creates a new PQueue for integers
func NewInts(data []int) *PQueue[int] {
	return New(data, func(a, b int) bool { return a < b })
}

// NewFloats creates a new PQueue for floats
func NewFloats(data []float64) *PQueue[float64] {
	return New(data, func(a, b float64) bool { return a < b })
}

// NewStrings creates a new PQueue for strings
func NewStrings(data []string) *PQueue[string] {
	return New(data, func(a, b string) bool { return a < b })
}

// NewBytes creates a new PQueue for byte slices
func NewBytes(data [][]byte) *PQueue[[]byte] {
	return New(data, func(a, b []byte) bool {
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return a[i] < b[i]
			}
		}
		return len(a) < len(b)
	})
}

// NewRunes creates a new PQueue for rune slices
func NewRunes(data [][]rune) *PQueue[[]rune] {
	return New(data, func(a, b []rune) bool {
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return a[i] < b[i]
			}
		}
		return len(a) < len(b)
	})
}

// NewComparable creates a new PQueue for any comparable type
func NewComparable[T comparable](data []T, less func(T, T) bool) *PQueue[T] {
	return New(data, less)
}

// Comparable interface for types that can be compared
type Comparable interface {
	CompareTo(other interface{}) int
}

// NewWithComparable creates a PQueue for types that implement Comparable
func NewWithComparable[T Comparable](data []T) *PQueue[T] {
	return New(data, func(a, b T) bool {
		return a.CompareTo(b) < 0
	})
}

// Size returns the number of elements in the queue
func (pq *PQueue[T]) Size() int {
	return pq.size
}

// IsEmpty returns true if the queue is empty
func (pq *PQueue[T]) IsEmpty() bool {
	return pq.size == 0
}

// Push adds an element to the queue
func (pq *PQueue[T]) Push(item T) {
	if pq.size >= len(pq.data) {
		// Grow the slice
		newSize := len(pq.data) * 2
		if newSize == 0 {
			newSize = 1 // Start with size 1 if empty
		}
		newData := make([]T, newSize)
		copy(newData, pq.data[:pq.size])
		pq.data = newData
	}
	pq.data[pq.size] = item
	pq.size++
}

// Pop removes and returns the smallest element
func (pq *PQueue[T]) Pop() (T, error) {
	var zero T
	if pq.size == 0 {
		return zero, fmt.Errorf("queue is empty")
	}

	// Find minimum element
	minIdx := 0
	for i := 1; i < pq.size; i++ {
		if pq.less(pq.data[i], pq.data[minIdx]) {
			minIdx = i
		}
	}

	result := pq.data[minIdx]
	// Move last element to the position of removed element
	pq.data[minIdx] = pq.data[pq.size-1]
	pq.size--

	return result, nil
}

// Peek returns the smallest element without removing it
func (pq *PQueue[T]) Peek() (T, error) {
	var zero T
	if pq.size == 0 {
		return zero, fmt.Errorf("queue is empty")
	}

	minIdx := 0
	for i := 1; i < pq.size; i++ {
		if pq.less(pq.data[i], pq.data[minIdx]) {
			minIdx = i
		}
	}

	return pq.data[minIdx], nil
}

// Sort sorts the queue using the optimal algorithm based on data characteristics
func (pq *PQueue[T]) Sort() {
	pq.SortWithStrategy(AutoStrategy)
}

// SortWithStrategy sorts using a specific strategy
func (pq *PQueue[T]) SortWithStrategy(strategy SortStrategy) {
	if pq.size <= 1 {
		return
	}

	actualStrategy := strategy
	if strategy == AutoStrategy {
		actualStrategy = pq.chooseOptimalStrategy()
	}

	switch actualStrategy {
	case InsertionStrategy:
		pq.insertionSort()
	case TimsortStrategy:
		pq.timsort()
	case IntrosortStrategy:
		pq.introsort()
	case MergeStrategy:
		pq.mergeSort()
	case QuickStrategy:
		pq.quickSort()
	case RadixStrategy:
		if pq.dataType == IntegerType {
			pq.radixSort()
		} else {
			pq.quickSort() // fallback
		}
	case CountingStrategy:
		if pq.dataType == IntegerType {
			pq.countingSort()
		} else {
			pq.quickSort() // fallback
		}
	default:
		pq.quickSort()
	}
}

// ToSlice returns a copy of the internal data
func (pq *PQueue[T]) ToSlice() []T {
	result := make([]T, pq.size)
	copy(result, pq.data[:pq.size])
	return result
}

// GetDataType returns the inferred data type for debugging purposes
func (pq *PQueue[T]) GetDataType() DataType {
	return pq.dataType
}

// GetDataTypeName returns a human-readable name for the data type
func (pq *PQueue[T]) GetDataTypeName() string {
	switch pq.dataType {
	case IntegerType:
		return "Integer"
	case FloatType:
		return "Float"
	case StringType:
		return "String"
	case SliceType:
		return "Slice"
	case ArrayType:
		return "Array"
	case StructType:
		return "Struct"
	case MapType:
		return "Map"
	case PointerType:
		return "Pointer"
	case InterfaceType:
		return "Interface"
	case ChannelType:
		return "Channel"
	case FuncType:
		return "Function"
	default:
		return "Generic"
	}
}

// inferDataType attempts to determine the data type using reflection
func inferDataType[T any](data []T) DataType {
	if len(data) == 0 {
		return GenericType
	}

	// Get the type of the first element
	t := reflect.TypeOf(data[0])

	// Handle pointers by getting the underlying type
	if t.Kind() == reflect.Ptr {
		return PointerType
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return IntegerType
	case reflect.Float32, reflect.Float64:
		return FloatType
	case reflect.String:
		return StringType
	case reflect.Slice:
		return SliceType
	case reflect.Array:
		return ArrayType
	case reflect.Struct:
		return StructType
	case reflect.Map:
		return MapType
	case reflect.Interface:
		return InterfaceType
	case reflect.Chan:
		return ChannelType
	case reflect.Func:
		return FuncType
	default:
		return GenericType
	}
}

// chooseOptimalStrategy selects the best sorting algorithm based on data characteristics
func (pq *PQueue[T]) chooseOptimalStrategy() SortStrategy {
	n := pq.size

	// For very small arrays, use insertion sort
	if n <= 16 {
		return InsertionStrategy
	}

	// Check if data is nearly sorted
	if pq.isNearlySorted() {
		return InsertionStrategy
	}

	// For integer data with small range, use counting or radix sort
	if pq.dataType == IntegerType && n > 100 {
		if pq.hasSmallRange() {
			return CountingStrategy
		}
		return RadixStrategy
	}

	// For strings, use specialized string sorting
	if pq.dataType == StringType {
		if n > 1000 {
			return IntrosortStrategy // Good for large string datasets
		}
		return TimsortStrategy // Good for strings with patterns
	}

	// For slices and arrays, use stable sorting
	if pq.dataType == SliceType || pq.dataType == ArrayType {
		return MergeStrategy // Stable and predictable
	}

	// For structs and complex types, use comparison-based sorts
	if pq.dataType == StructType || pq.dataType == InterfaceType {
		if n > 1000 {
			return IntrosortStrategy
		}
		return TimsortStrategy
	}

	// For pointers, maps, channels, functions - use generic approach
	if pq.dataType == PointerType || pq.dataType == MapType ||
		pq.dataType == ChannelType || pq.dataType == FuncType {
		return QuickStrategy // Simple and effective for these types
	}

	// For large datasets, use introsort (hybrid approach)
	if n > 1000 {
		return IntrosortStrategy
	}

	// Default to timsort for general purpose
	return TimsortStrategy
}
