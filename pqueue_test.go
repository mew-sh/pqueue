package pqueue

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

// TestPQueueBasicOperations tests basic priority queue operations
func TestPQueueBasicOperations(t *testing.T) {
	data := []int{6, 5, 4, 9, 2, 7, 1, 8}
	pq := NewInts(data)

	// Test size
	if pq.Size() != 8 {
		t.Errorf("Expected size 8, got %d", pq.Size())
	}

	// Test IsEmpty
	if pq.IsEmpty() {
		t.Error("Expected queue not to be empty")
	}

	// Test Peek
	min, err := pq.Peek()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if min != 1 {
		t.Errorf("Expected min element 1, got %d", min)
	}

	// Test Pop
	popped, err := pq.Pop()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if popped != 1 {
		t.Errorf("Expected popped element 1, got %d", popped)
	}
	if pq.Size() != 7 {
		t.Errorf("Expected size 7 after pop, got %d", pq.Size())
	}

	// Test Push
	pq.Push(0)
	if pq.Size() != 8 {
		t.Errorf("Expected size 8 after push, got %d", pq.Size())
	}

	// Test ToSlice
	slice := pq.ToSlice()
	if len(slice) != 8 {
		t.Errorf("Expected slice length 8, got %d", len(slice))
	}
}

// TestEmptyQueue tests operations on empty queue
func TestEmptyQueue(t *testing.T) {
	pq := NewInts([]int{})

	// Test empty queue properties
	if !pq.IsEmpty() {
		t.Error("Expected empty queue to be empty")
	}
	if pq.Size() != 0 {
		t.Errorf("Expected size 0, got %d", pq.Size())
	}

	// Test Pop on empty queue
	_, err := pq.Pop()
	if err == nil {
		t.Error("Expected error when popping from empty queue")
	}

	// Test Peek on empty queue
	_, err = pq.Peek()
	if err == nil {
		t.Error("Expected error when peeking empty queue")
	}

	// Test Push to empty queue
	pq.Push(42)
	if pq.IsEmpty() {
		t.Error("Expected queue not to be empty after push")
	}
	if pq.Size() != 1 {
		t.Errorf("Expected size 1, got %d", pq.Size())
	}
}

// TestQueueGrowth tests dynamic queue growth
func TestQueueGrowth(t *testing.T) {
	pq := NewInts([]int{1})
	
	// Add many elements to trigger growth
	for i := 2; i <= 100; i++ {
		pq.Push(i)
	}

	if pq.Size() != 100 {
		t.Errorf("Expected size 100, got %d", pq.Size())
	}

	// Verify all elements are present
	pq.Sort()
	sorted := pq.ToSlice()
	for i := 0; i < 100; i++ {
		if sorted[i] != i+1 {
			t.Errorf("Expected element %d at position %d, got %d", i+1, i, sorted[i])
		}
	}
}

// TestIntegerQueue tests integer queue operations
func TestIntegerQueue(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{
			name: "already sorted",
			data: []int{1, 2, 3, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "reverse sorted",
			data: []int{5, 4, 3, 2, 1},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "random order",
			data: []int{3, 1, 4, 1, 5, 9, 2, 6},
			want: []int{1, 1, 2, 3, 4, 5, 6, 9},
		},
		{
			name: "single element",
			data: []int{42},
			want: []int{42},
		},
		{
			name: "duplicates",
			data: []int{3, 3, 3, 1, 1, 2, 2},
			want: []int{1, 1, 2, 2, 3, 3, 3},
		},
		{
			name: "negative numbers",
			data: []int{-3, -1, -4, 1, 5, 0, -2},
			want: []int{-4, -3, -2, -1, 0, 1, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewInts(tt.data)
			pq.Sort()
			got := pq.ToSlice()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFloatQueue tests float queue operations
func TestFloatQueue(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		want []float64
	}{
		{
			name: "basic floats",
			data: []float64{3.14, 2.71, 1.41, 1.73},
			want: []float64{1.41, 1.73, 2.71, 3.14},
		},
		{
			name: "negative floats",
			data: []float64{-1.5, 2.3, -0.7, 0.0, 1.2},
			want: []float64{-1.5, -0.7, 0.0, 1.2, 2.3},
		},
		{
			name: "very small differences",
			data: []float64{1.0001, 1.0002, 1.0000},
			want: []float64{1.0000, 1.0001, 1.0002},
		},
		{
			name: "special values",
			data: []float64{math.Inf(1), math.Inf(-1), 0.0, 1.0},
			want: []float64{math.Inf(-1), 0.0, 1.0, math.Inf(1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewFloats(tt.data)
			pq.Sort()
			got := pq.ToSlice()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStringQueue tests string queue operations
func TestStringQueue(t *testing.T) {
	tests := []struct {
		name string
		data []string
		want []string
	}{
		{
			name: "basic strings",
			data: []string{"zebra", "apple", "banana", "cherry"},
			want: []string{"apple", "banana", "cherry", "zebra"},
		},
		{
			name: "case sensitive",
			data: []string{"Apple", "apple", "Banana", "banana"},
			want: []string{"Apple", "Banana", "apple", "banana"},
		},
		{
			name: "empty strings",
			data: []string{"", "a", "", "b"},
			want: []string{"", "", "a", "b"},
		},
		{
			name: "unicode strings",
			data: []string{"世界", "hello", "αβγ", "مرحبا"},
			want: []string{"hello", "αβγ", "مرحبا", "世界"},
		},
		{
			name: "single character",
			data: []string{"z", "a", "m", "b"},
			want: []string{"a", "b", "m", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewStrings(tt.data)
			pq.Sort()
			got := pq.ToSlice()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestByteSliceQueue tests byte slice queue operations
func TestByteSliceQueue(t *testing.T) {
	tests := []struct {
		name string
		data [][]byte
		want [][]byte
	}{
		{
			name: "basic byte slices",
			data: [][]byte{
				[]byte("zebra"),
				[]byte("apple"),
				[]byte("banana"),
			},
			want: [][]byte{
				[]byte("apple"),
				[]byte("banana"),
				[]byte("zebra"),
			},
		},
		{
			name: "different lengths",
			data: [][]byte{
				[]byte("abc"),
				[]byte("ab"),
				[]byte("abcd"),
				[]byte("a"),
			},
			want: [][]byte{
				[]byte("a"),
				[]byte("ab"),
				[]byte("abc"),
				[]byte("abcd"),
			},
		},
		{
			name: "empty slices",
			data: [][]byte{
				[]byte("hello"),
				[]byte(""),
				[]byte("world"),
				[]byte(""),
			},
			want: [][]byte{
				[]byte(""),
				[]byte(""),
				[]byte("hello"),
				[]byte("world"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewBytes(tt.data)
			pq.Sort()
			got := pq.ToSlice()

			if !deepEqualByteSlices(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuneSliceQueue tests rune slice queue operations
func TestRuneSliceQueue(t *testing.T) {
	tests := []struct {
		name string
		data [][]rune
		want [][]rune
	}{
		{
			name: "unicode runes",
			data: [][]rune{
				[]rune("世界"),
				[]rune("hello"),
				[]rune("αβγ"),
			},
			want: [][]rune{
				[]rune("hello"),
				[]rune("αβγ"),
				[]rune("世界"),
			},
		},
		{
			name: "mixed runes",
			data: [][]rune{
				[]rune("😊🎉"),
				[]rune("abc"),
				[]rune("123"),
			},
			want: [][]rune{
				[]rune("123"),
				[]rune("abc"),
				[]rune("😊🎉"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewRunes(tt.data)
			pq.Sort()
			got := pq.ToSlice()

			if !deepEqualRuneSlices(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCustomTypes tests custom type queue operations
func TestCustomTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"David", 25},
	}

	// Sort by age, then by name
	pq := New(people, func(a, b Person) bool {
		if a.Age != b.Age {
			return a.Age < b.Age
		}
		return a.Name < b.Name
	})

	pq.Sort()
	sorted := pq.ToSlice()

	expected := []Person{
		{"Bob", 25},
		{"David", 25},
		{"Alice", 30},
		{"Charlie", 35},
	}

	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Sort() = %v, want %v", sorted, expected)
	}
}

// TestDataTypeInference tests data type inference
func TestDataTypeInference(t *testing.T) {
	tests := []struct {
		name     string
		pq       interface{}
		expected DataType
	}{
		{
			name:     "integers",
			pq:       NewInts([]int{1, 2, 3}),
			expected: IntegerType,
		},
		{
			name:     "floats",
			pq:       NewFloats([]float64{1.1, 2.2, 3.3}),
			expected: FloatType,
		},
		{
			name:     "strings",
			pq:       NewStrings([]string{"a", "b", "c"}),
			expected: StringType,
		},
		{
			name:     "byte slices",
			pq:       NewBytes([][]byte{[]byte("a"), []byte("b")}),
			expected: SliceType,
		},
		{
			name:     "rune slices",
			pq:       NewRunes([][]rune{[]rune("a"), []rune("b")}),
			expected: SliceType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dataType DataType
			switch v := tt.pq.(type) {
			case *PQueue[int]:
				dataType = v.GetDataType()
			case *PQueue[float64]:
				dataType = v.GetDataType()
			case *PQueue[string]:
				dataType = v.GetDataType()
			case *PQueue[[]byte]:
				dataType = v.GetDataType()
			case *PQueue[[]rune]:
				dataType = v.GetDataType()
			}

			if dataType != tt.expected {
				t.Errorf("GetDataType() = %v, want %v", dataType, tt.expected)
			}
		})
	}
}

// TestPointerTypes tests pointer type handling
func TestPointerTypes(t *testing.T) {
	values := []int{5, 2, 8, 1, 9}
	ptrs := make([]*int, len(values))
	for i := range values {
		ptrs[i] = &values[i]
	}

	// Add a nil pointer
	ptrs = append(ptrs, nil)

	pq := New(ptrs, func(a, b *int) bool {
		if a == nil {
			return b != nil // nil values come first
		}
		if b == nil {
			return false
		}
		return *a < *b
	})

	pq.Sort()
	sorted := pq.ToSlice()

	// First element should be nil
	if sorted[0] != nil {
		t.Error("Expected nil pointer to be first")
	}

	// Check that non-nil pointers are sorted
	for i := 1; i < len(sorted)-1; i++ {
		if sorted[i] == nil || sorted[i+1] == nil {
			continue
		}
		if *sorted[i] > *sorted[i+1] {
			t.Errorf("Pointers not sorted correctly: %d > %d", *sorted[i], *sorted[i+1])
		}
	}
}

// ComparableInt is a test type that implements Comparable
type ComparableInt int

func (c ComparableInt) CompareTo(other interface{}) int {
	if o, ok := other.(ComparableInt); ok {
		if c < o {
			return -1
		} else if c > o {
			return 1
		}
		return 0
	}
	return 0
}

// TestComparableInterface tests the Comparable interface
func TestComparableInterface(t *testing.T) {
	data := []ComparableInt{3, 1, 4, 1, 5}
	pq := NewWithComparable(data)
	pq.Sort()
	sorted := pq.ToSlice()

	expected := []ComparableInt{1, 1, 3, 4, 5}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Sort() = %v, want %v", sorted, expected)
	}
}

// Helper functions
func deepEqualByteSlices(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func deepEqualRuneSlices(a, b [][]rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// TestRandomData tests sorting with random data
func TestRandomData(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			// Generate random data
			data := make([]int, size)
			for i := range data {
				data[i] = rand.Intn(1000)
			}
			
			// Sort with PQueue
			pq := NewInts(data)
			pq.Sort()
			pqSorted := pq.ToSlice()
			
			// Sort with standard library
			stdSorted := make([]int, len(data))
			copy(stdSorted, data)
			sort.Ints(stdSorted)
			
			// Compare results
			if !reflect.DeepEqual(pqSorted, stdSorted) {
				t.Errorf("PQueue sort doesn't match standard sort for size %d", size)
			}
		})
	}
}
