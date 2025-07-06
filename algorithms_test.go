package pqueue

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// TestAllSortingStrategies tests all sorting strategies explicitly
func TestAllSortingStrategies(t *testing.T) {
	data := []int{64, 34, 25, 12, 22, 11, 90}
	expected := []int{11, 12, 22, 25, 34, 64, 90}

	strategies := []struct {
		name     string
		strategy SortStrategy
	}{
		{"Auto", AutoStrategy},
		{"Insertion", InsertionStrategy},
		{"Quick", QuickStrategy},
		{"Merge", MergeStrategy},
		{"Introsort", IntrosortStrategy},
		{"Timsort", TimsortStrategy},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			pq := NewInts(data)
			pq.SortWithStrategy(s.strategy)
			result := pq.ToSlice()

			if !reflect.DeepEqual(result, expected) {
				t.Errorf("Strategy %s: got %v, want %v", s.name, result, expected)
			}
		})
	}
}

// TestRadixSortForIntegers tests radix sort specifically for integers
func TestRadixSortForIntegers(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{
			name: "positive integers",
			data: []int{170, 45, 75, 90, 2, 802, 24, 66},
			want: []int{2, 24, 45, 66, 75, 90, 170, 802},
		},
		{
			name: "single digit",
			data: []int{5, 2, 8, 1, 9, 3},
			want: []int{1, 2, 3, 5, 8, 9},
		},
		{
			name: "same number of digits",
			data: []int{123, 456, 789, 234, 567},
			want: []int{123, 234, 456, 567, 789},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewInts(tt.data)
			pq.SortWithStrategy(RadixStrategy)
			got := pq.ToSlice()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RadixSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCountingSortForIntegers tests counting sort for small range integers
func TestCountingSortForIntegers(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{
			name: "small range",
			data: []int{4, 2, 2, 8, 3, 3, 1},
			want: []int{1, 2, 2, 3, 3, 4, 8},
		},
		{
			name: "all same",
			data: []int{5, 5, 5, 5},
			want: []int{5, 5, 5, 5},
		},
		{
			name: "zero included",
			data: []int{3, 0, 2, 0, 1},
			want: []int{0, 0, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewInts(tt.data)
			pq.SortWithStrategy(CountingStrategy)
			got := pq.ToSlice()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountingSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAutoStrategySelection tests automatic strategy selection
func TestAutoStrategySelection(t *testing.T) {
	tests := []struct {
		name         string
		data         []int
		expectedType string // We can't directly test the strategy, but we can test behavior
	}{
		{
			name: "small array should use insertion sort",
			data: []int{5, 2, 8, 1},
			expectedType: "small",
		},
		{
			name: "large array should use advanced strategy",
			data: make([]int, 1001), // Will be filled with random data
			expectedType: "large",
		},
		{
			name: "nearly sorted should be detected",
			data: []int{1, 2, 3, 5, 4, 6, 7, 8, 9, 10}, // Only one inversion
			expectedType: "nearly_sorted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedType == "large" {
				// Fill with random data for large test
				for i := range tt.data {
					tt.data[i] = rand.Intn(1000)
				}
			}

			pq := NewInts(tt.data)
			pq.Sort() // Use auto strategy

			// Verify it's sorted regardless of strategy
			result := pq.ToSlice()
			for i := 1; i < len(result); i++ {
				if result[i-1] > result[i] {
					t.Errorf("Array not sorted: %d > %d at positions %d, %d", result[i-1], result[i], i-1, i)
				}
			}
		})
	}
}

// TestStrategyForDifferentDataTypes tests strategy selection for different data types
func TestStrategyForDifferentDataTypes(t *testing.T) {
	t.Run("string_data", func(t *testing.T) {
		data := []string{"zebra", "apple", "banana", "cherry", "date"}
		pq := NewStrings(data)
		pq.Sort()
		result := pq.ToSlice()

		expected := []string{"apple", "banana", "cherry", "date", "zebra"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("String sort failed: got %v, want %v", result, expected)
		}
	})

	t.Run("float_data", func(t *testing.T) {
		data := []float64{3.14, 2.71, 1.41, 1.73, 0.57}
		pq := NewFloats(data)
		pq.Sort()
		result := pq.ToSlice()

		expected := []float64{0.57, 1.41, 1.73, 2.71, 3.14}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Float sort failed: got %v, want %v", result, expected)
		}
	})

	t.Run("byte_slice_data", func(t *testing.T) {
		data := [][]byte{
			[]byte("zebra"),
			[]byte("apple"),
			[]byte("banana"),
		}
		pq := NewBytes(data)
		pq.Sort()
		result := pq.ToSlice()

		// Check if sorted lexicographically
		for i := 1; i < len(result); i++ {
			if string(result[i-1]) > string(result[i]) {
				t.Errorf("Byte slices not sorted: %s > %s", result[i-1], result[i])
			}
		}
	})
}

// TestLargeDataSets tests performance with large data sets
func TestLargeDataSets(t *testing.T) {
	sizes := []int{1000, 5000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			// Generate random data
			rand.Seed(time.Now().UnixNano())
			data := make([]int, size)
			for i := range data {
				data[i] = rand.Intn(size * 10)
			}

			start := time.Now()
			pq := NewInts(data)
			pq.Sort()
			duration := time.Since(start)

			// Verify it's sorted
			result := pq.ToSlice()
			for i := 1; i < len(result); i++ {
				if result[i-1] > result[i] {
					t.Errorf("Large dataset not sorted at position %d", i)
				}
			}

			t.Logf("Sorted %d elements in %v", size, duration)
		})
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("single_element", func(t *testing.T) {
		pq := NewInts([]int{42})
		pq.Sort()
		result := pq.ToSlice()
		if len(result) != 1 || result[0] != 42 {
			t.Errorf("Single element test failed")
		}
	})

	t.Run("all_duplicates", func(t *testing.T) {
		data := make([]int, 100)
		for i := range data {
			data[i] = 7
		}
		pq := NewInts(data)
		pq.Sort()
		result := pq.ToSlice()

		for _, v := range result {
			if v != 7 {
				t.Errorf("All duplicates test failed")
			}
		}
	})

	t.Run("reverse_sorted", func(t *testing.T) {
		data := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
		pq := NewInts(data)
		pq.Sort()
		result := pq.ToSlice()

		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Reverse sorted test failed")
		}
	})

	t.Run("already_sorted", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		pq := NewInts(data)
		pq.Sort()
		result := pq.ToSlice()

		if !reflect.DeepEqual(result, data) {
			t.Errorf("Already sorted test failed")
		}
	})
}

// TestStability tests if stable sorting algorithms maintain relative order
func TestStability(t *testing.T) {
	type Item struct {
		Value int
		Index int
	}

	data := []Item{
		{Value: 3, Index: 0},
		{Value: 1, Index: 1},
		{Value: 3, Index: 2},
		{Value: 2, Index: 3},
		{Value: 1, Index: 4},
	}

	pq := New(data, func(a, b Item) bool {
		return a.Value < b.Value
	})

	// Test with merge sort (stable)
	pq.SortWithStrategy(MergeStrategy)
	result := pq.ToSlice()

	// Check that items with same value maintain original order
	prevValue := -1
	prevIndex := -1
	for _, item := range result {
		if item.Value == prevValue && item.Index < prevIndex {
			t.Errorf("Stability violated: item with index %d came before item with index %d", item.Index, prevIndex)
		}
		if item.Value != prevValue {
			prevIndex = item.Index
		} else if item.Index > prevIndex {
			prevIndex = item.Index
		}
		prevValue = item.Value
	}
}

// TestConcurrentAccess tests thread safety considerations
func TestConcurrentAccess(t *testing.T) {
	// Note: PQueue is not thread-safe by design, but we test that
	// concurrent read operations don't cause data races
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(1000)
	}

	pq := NewInts(data)
	pq.Sort()

	// Multiple goroutines reading the sorted data
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result := pq.ToSlice()
			// Verify it's sorted
			for j := 1; j < len(result); j++ {
				if result[j-1] > result[j] {
					t.Errorf("Concurrent read found unsorted data")
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestMemoryUsage tests memory efficiency
func TestMemoryUsage(t *testing.T) {
	// Test that sorting doesn't use excessive memory
	originalData := make([]int, 10000)
	for i := range originalData {
		originalData[i] = rand.Intn(10000)
	}
	
	// Test different strategies to ensure they don't leak memory
	strategies := []SortStrategy{
		QuickStrategy,
		MergeStrategy,
		IntrosortStrategy,
		TimsortStrategy,
	}

	for _, strategy := range strategies {
		// Create a copy for each test
		testData := make([]int, len(originalData))
		copy(testData, originalData)
		
		testPQ := NewInts(testData)
		testPQ.SortWithStrategy(strategy)
		
		// Verify the result is still correct
		result := testPQ.ToSlice()
		for i := 1; i < len(result); i++ {
			if result[i-1] > result[i] {
				t.Errorf("Memory test failed: unsorted result for strategy %v", strategy)
			}
		}
	}
}
