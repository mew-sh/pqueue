package pqueue

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// BenchmarkPQueueVsStandardSort compares PQueue with Go's standard sort
func BenchmarkPQueueVsStandardSort(b *testing.B) {
	sizes := []int{100, 1000, 5000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("PQueue_Size_%d", size), func(b *testing.B) {
			data := generateRandomInts(size)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testData := make([]int, len(data))
				copy(testData, data)
				pq := NewInts(testData)
				b.StartTimer()
				
				pq.Sort()
			}
		})

		b.Run(fmt.Sprintf("StandardSort_Size_%d", size), func(b *testing.B) {
			data := generateRandomInts(size)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testData := make([]int, len(data))
				copy(testData, data)
				b.StartTimer()
				
				sort.Ints(testData)
			}
		})
	}
}

// BenchmarkSortingStrategies benchmarks different sorting strategies
func BenchmarkSortingStrategies(b *testing.B) {
	strategies := []struct {
		name     string
		strategy SortStrategy
	}{
		{"Auto", AutoStrategy},
		{"Quick", QuickStrategy},
		{"Merge", MergeStrategy},
		{"Introsort", IntrosortStrategy},
		{"Timsort", TimsortStrategy},
		{"Insertion", InsertionStrategy},
	}

	sizes := []int{100, 1000, 5000}

	for _, size := range sizes {
		data := generateRandomInts(size)
		
		for _, s := range strategies {
			b.Run(fmt.Sprintf("%s_Size_%d", s.name, size), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					testData := make([]int, len(data))
					copy(testData, data)
					pq := NewInts(testData)
					b.StartTimer()
					
					pq.SortWithStrategy(s.strategy)
				}
			})
		}
	}
}

// BenchmarkSpecializedSorts benchmarks radix and counting sort for integers
func BenchmarkSpecializedSorts(b *testing.B) {
	b.Run("RadixSort", func(b *testing.B) {
		data := generateRandomInts(5000)
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.SortWithStrategy(RadixStrategy)
		}
	})

	b.Run("CountingSort", func(b *testing.B) {
		// Generate data with small range for counting sort
		data := make([]int, 5000)
		for i := range data {
			data[i] = rand.Intn(100) // Small range
		}
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.SortWithStrategy(CountingStrategy)
		}
	})
}

// BenchmarkDataTypes benchmarks different data types
func BenchmarkDataTypes(b *testing.B) {
	size := 1000

	b.Run("Integers", func(b *testing.B) {
		data := generateRandomInts(size)
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})

	b.Run("Floats", func(b *testing.B) {
		data := generateRandomFloats(size)
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]float64, len(data))
			copy(testData, data)
			pq := NewFloats(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})

	b.Run("Strings", func(b *testing.B) {
		data := generateRandomStrings(size)
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			pq := NewStrings(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})

	b.Run("ByteSlices", func(b *testing.B) {
		data := generateRandomByteSlices(size)
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([][]byte, len(data))
			for j, bs := range data {
				testData[j] = make([]byte, len(bs))
				copy(testData[j], bs)
			}
			pq := NewBytes(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})
}

// BenchmarkPriorityQueueOperations benchmarks priority queue operations
func BenchmarkPriorityQueueOperations(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		pq := NewInts([]int{})
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pq.Push(rand.Intn(1000))
		}
	})

	b.Run("Pop", func(b *testing.B) {
		data := generateRandomInts(b.N + 1000) // Ensure we have enough elements
		pq := NewInts(data)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if pq.Size() > 0 {
				pq.Pop()
			}
		}
	})

	b.Run("Peek", func(b *testing.B) {
		data := generateRandomInts(1000)
		pq := NewInts(data)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pq.Peek()
		}
	})
}

// BenchmarkWorstCaseScenarios benchmarks worst-case scenarios
func BenchmarkWorstCaseScenarios(b *testing.B) {
	size := 1000

	b.Run("ReverseSorted", func(b *testing.B) {
		data := make([]int, size)
		for i := range data {
			data[i] = size - i
		}
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})

	b.Run("AllDuplicates", func(b *testing.B) {
		data := make([]int, size)
		for i := range data {
			data[i] = 42
		}
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})

	b.Run("NearlySorted", func(b *testing.B) {
		data := make([]int, size)
		for i := range data {
			data[i] = i
		}
		// Create a few inversions
		for i := 0; i < size/20; i++ {
			j := rand.Intn(size-1)
			data[j], data[j+1] = data[j+1], data[j]
		}
		
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]int, len(data))
			copy(testData, data)
			pq := NewInts(testData)
			b.StartTimer()
			
			pq.Sort()
		}
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("SmallArrays", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := generateRandomInts(16)
			pq := NewInts(data)
			pq.Sort()
		}
	})

	b.Run("MediumArrays", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := generateRandomInts(1000)
			pq := NewInts(data)
			pq.Sort()
		}
	})

	b.Run("LargeArrays", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := generateRandomInts(10000)
			pq := NewInts(data)
			pq.Sort()
		}
	})
}

// Helper functions for generating test data
func generateRandomInts(size int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(size * 10)
	}
	return data
}

func generateRandomFloats(size int) []float64 {
	rand.Seed(time.Now().UnixNano())
	data := make([]float64, size)
	for i := range data {
		data[i] = rand.Float64() * 1000
	}
	return data
}

func generateRandomStrings(size int) []string {
	rand.Seed(time.Now().UnixNano())
	data := make([]string, size)
	for i := range data {
		length := rand.Intn(10) + 1
		bytes := make([]byte, length)
		for j := range bytes {
			bytes[j] = byte(rand.Intn(26) + 'a')
		}
		data[i] = string(bytes)
	}
	return data
}

func generateRandomByteSlices(size int) [][]byte {
	rand.Seed(time.Now().UnixNano())
	data := make([][]byte, size)
	for i := range data {
		length := rand.Intn(10) + 1
		bytes := make([]byte, length)
		for j := range bytes {
			bytes[j] = byte(rand.Intn(256))
		}
		data[i] = bytes
	}
	return data
}

// BenchmarkCustomTypes benchmarks custom struct types
func BenchmarkCustomTypes(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	people := make([]Person, 1000)
	for i := range people {
		people[i] = Person{
			Name: fmt.Sprintf("Person%d", rand.Intn(1000)),
			Age:  rand.Intn(100),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testData := make([]Person, len(people))
		copy(testData, people)
		pq := New(testData, func(a, b Person) bool {
			if a.Age != b.Age {
				return a.Age < b.Age
			}
			return a.Name < b.Name
		})
		b.StartTimer()
		
		pq.Sort()
	}
}
