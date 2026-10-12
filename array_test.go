package js

import (
	"cmp"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestArray_Readme(t *testing.T) {
	yourSlice := []string{"hello", "world"}
	joined := Array[string](yourSlice).Join("-")
	if joined != "hello-world" {
		t.Error("Expected 'hello-world', got ", joined)
	}
	t.Run("item filtering", func(t *testing.T) {
		type Item struct {
			Name  string
			Price float64
			Keep  bool
		}

		shoppingCart := []Item{
			{Name: "Laptop", Price: 1000, Keep: true},
			{Name: "Mouse", Price: 50, Keep: false},
			{Name: "Keyboard", Price: 100, Keep: true},
			{Name: "HDMI Cable", Price: 20, Keep: true},
		}

		res := ArrayOf(shoppingCart...).
			Filter(func(v Item) bool { return v.Keep }).                            // Step 1: Keep only checked items
			Map(func(v Item) float64 { return v.Price * 0.9 }).                     // Step 2: Apply 10% discount
			Reduce(func(sum float64, price float64) float64 { return sum + price }) // Step 3: Accumulate the total

		expected := 1008.0
		if res != expected {
			t.Errorf("expected: %f, got %f", expected, res)
		}
	})

	t.Run("processing text numbers", func(t *testing.T) {
		convertToNumber := func(s string) *float64 {
			value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if err != nil {
				return nil
			}
			return &value
		}
		res := ArrayOf("10", "abc", "25.5", "  42  ").
			Map(convertToNumber).
			Filter(func(n *float64) bool { return n != nil }).
			Reduce(func(acc float64, v *float64) float64 { return acc + *v }, 0)

		expected := 77.5
		if res != expected {
			t.Errorf("expected: %f, got %f", expected, res)
		}
	})
}

func TestArray_ForRange(t *testing.T) {
	for i, v := range ArrayOf(1, 2) {
		if i == 0 && v != 1 {
			t.Errorf("expected: %d, got %d", 1, v)
		}
		if i == 1 && v != 2 {
			t.Errorf("expected: %d, got %d", 2, v)
		}
	}
}

func TestArray_Length(t *testing.T) {
	var a = Array[int]{1, 2, 3}
	if a.Length() != 3 {
		t.Errorf("len(a) should be 3, got %d", len(a))
	}
}

func TestArray_Push(t *testing.T) {
	a := Array[int]{1, 2}

	// Push single element
	newLen := a.Push(3)
	if newLen != 3 {
		t.Errorf("expected new length 3, got %d", newLen)
	}
	if !reflect.DeepEqual(a, Array[int]{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", a)
	}

	// Push multiple elements
	newLen = a.Push(4, 5)
	if newLen != 5 {
		t.Errorf("expected new length 5, got %d", newLen)
	}
	if !reflect.DeepEqual(a, Array[int]{1, 2, 3, 4, 5}) {
		t.Errorf("expected [1 2 3 4 5], got %v", a)
	}

	// Push 0 elements
	newLen = a.Push()
	if newLen != 5 {
		t.Errorf("expected new length 5, got %d", newLen)
	}
	if !reflect.DeepEqual(a, Array[int]{1, 2, 3, 4, 5}) {
		t.Errorf("expected [1 2 3 4 5], got %v", a)
	}
}

func TestArray_Pop(t *testing.T) {
	a := Array[string]{"a", "b", "c"}

	// Pop "c"
	val, ok := a.Pop()
	if !ok || val != "c" {
		t.Errorf("expected ('c', true), got ('%s', %v)", val, ok)
	}
	if a.Length() != 2 {
		t.Errorf("expected length 2, got %d", a.Length())
	}

	// Pop "b"
	val, ok = a.Pop()
	if !ok || val != "b" {
		t.Errorf("expected ('b', true), got ('%s', %v)", val, ok)
	}

	// Pop "a"
	val, ok = a.Pop()
	if !ok || val != "a" {
		t.Errorf("expected ('a', true), got ('%s', %v)", val, ok)
	}
	if a.Length() != 0 {
		t.Errorf("expected length 0, got %d", a.Length())
	}

	// Pop from empty
	val, ok = a.Pop()
	if ok || val != "" {
		t.Errorf("expected ('', false), got ('%s', %v)", val, ok)
	}
	if a.Length() != 0 {
		t.Errorf("expected length 0, got %d", a.Length())
	}
}

func TestArray_Shift(t *testing.T) {
	a := Array[int]{10, 20, 30}

	// Shift 10
	val, ok := a.Shift()
	if !ok || val != 10 {
		t.Errorf("expected (10, true), got (%d, %v)", val, ok)
	}
	if !reflect.DeepEqual(a, Array[int]{20, 30}) {
		t.Errorf("expected [20 30], got %v", a)
	}

	// Shift 20
	val, ok = a.Shift()
	if !ok || val != 20 {
		t.Errorf("expected (20, true), got (%d, %v)", val, ok)
	}

	// Shift 30
	val, ok = a.Shift()
	if !ok || val != 30 {
		t.Errorf("expected (30, true), got (%d, %v)", val, ok)
	}
	if a.Length() != 0 {
		t.Errorf("expected length 0, got %d", a.Length())
	}

	// Shift from empty
	val, ok = a.Shift()
	if ok || val != 0 {
		t.Errorf("expected (0, false), got (%d, %v)", val, ok)
	}
}

func TestArray_Unshift(t *testing.T) {
	a := Array[int]{3, 4}

	// Unshift single element
	newLen := a.Unshift(2)
	if newLen != 3 {
		t.Errorf("expected new length 3, got %d", newLen)
	}
	if !reflect.DeepEqual(a, Array[int]{2, 3, 4}) {
		t.Errorf("expected [2 3 4], got %v", a)
	}

	// Unshift multiple elements
	newLen = a.Unshift(0, 1)
	if newLen != 5 {
		t.Errorf("expected new length 5, got %d", newLen)
	}
	if !reflect.DeepEqual(a, Array[int]{0, 1, 2, 3, 4}) {
		t.Errorf("expected [0 1 2 3 4], got %v", a)
	}

	// Unshift 0 elements
	newLen = a.Unshift()
	if newLen != 5 {
		t.Errorf("expected new length 5, got %d", newLen)
	}

	// Unshift onto empty array
	empty := Array[string]{}
	newLen = empty.Unshift("first", "second")
	if newLen != 2 {
		t.Errorf("expected new length 2, got %d", newLen)
	}
	if !reflect.DeepEqual(empty, Array[string]{"first", "second"}) {
		t.Errorf("expected ['first', 'second'], got %v", empty)
	}
}

func TestArray_MemoryLeakZeroing(t *testing.T) {
	type Item struct{ Val int }
	item1 := &Item{Val: 1}
	item2 := &Item{Val: 2}

	a := Array[*Item]{item1, item2}
	underlying := a[:2]

	val, ok := a.Pop()
	if !ok || val != item2 {
		t.Fatalf("expected item2, got %v", val)
	}
	if underlying[1] != nil {
		t.Errorf("expected underlying[1] to be nil after Pop, got %v", underlying[1])
	}

	val, ok = a.Shift()
	if !ok || val != item1 {
		t.Fatalf("expected item1, got %v", val)
	}
	if underlying[0] != nil {
		t.Errorf("expected underlying[0] to be nil after Shift, got %v", underlying[0])
	}
}

func TestArray_Splice(t *testing.T) {
	// 1. Remove 0 elements, insert elements at index 2
	t.Run("insert without removal", func(t *testing.T) {
		a := Array[string]{"angel", "clown", "mandarin", "sturgeon"}
		removed := a.Splice(2, 0, "drum")
		if len(removed) != 0 {
			t.Errorf("expected empty removed array, got %v", removed)
		}
		expected := Array[string]{"angel", "clown", "drum", "mandarin", "sturgeon"}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 2. Remove 1 element from index 3
	t.Run("remove 1 element", func(t *testing.T) {
		a := Array[string]{"angel", "clown", "drum", "mandarin", "sturgeon"}
		removed := a.Splice(3, 1)
		expectedRemoved := Array[string]{"mandarin"}
		if !reflect.DeepEqual(removed, expectedRemoved) {
			t.Errorf("expected removed %v, got %v", expectedRemoved, removed)
		}
		expected := Array[string]{"angel", "clown", "drum", "sturgeon"}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 3. Remove 1 element and insert new element (replace)
	t.Run("replace 1 element", func(t *testing.T) {
		a := Array[string]{"angel", "clown", "trumpet", "sturgeon"}
		removed := a.Splice(0, 2, "parrot", "anemone", "blue")
		expectedRemoved := Array[string]{"angel", "clown"}
		if !reflect.DeepEqual(removed, expectedRemoved) {
			t.Errorf("expected removed %v, got %v", expectedRemoved, removed)
		}
		expected := Array[string]{"parrot", "anemone", "blue", "trumpet", "sturgeon"}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 4. Negative start index
	t.Run("negative start index", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		removed := a.Splice(-2, 1)
		expectedRemoved := Array[int]{4}
		if !reflect.DeepEqual(removed, expectedRemoved) {
			t.Errorf("expected removed %v, got %v", expectedRemoved, removed)
		}
		expected := Array[int]{1, 2, 3, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 5. Negative start index exceeding length
	t.Run("negative start exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		removed := a.Splice(-10, 2)
		expectedRemoved := Array[int]{1, 2}
		if !reflect.DeepEqual(removed, expectedRemoved) {
			t.Errorf("expected removed %v, got %v", expectedRemoved, removed)
		}
		expected := Array[int]{3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 6. Start index exceeding length
	t.Run("start exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2}
		removed := a.Splice(10, 5, 3, 4)
		if len(removed) != 0 {
			t.Errorf("expected empty removed, got %v", removed)
		}
		expected := Array[int]{1, 2, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 7. Negative deleteCount (clamped to 0)
	t.Run("negative deleteCount", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		removed := a.Splice(1, -5, 99)
		if len(removed) != 0 {
			t.Errorf("expected empty removed, got %v", removed)
		}
		expected := Array[int]{1, 99, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 8. DeleteCount exceeding remaining elements
	t.Run("deleteCount exceeding remaining", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		removed := a.Splice(2, 99)
		expectedRemoved := Array[int]{3, 4, 5}
		if !reflect.DeepEqual(removed, expectedRemoved) {
			t.Errorf("expected removed %v, got %v", expectedRemoved, removed)
		}
		expected := Array[int]{1, 2}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 9. Splice on empty array
	t.Run("empty array", func(t *testing.T) {
		a := Array[int]{}
		removed := a.Splice(0, 0, 10, 20)
		if len(removed) != 0 {
			t.Errorf("expected empty removed, got %v", removed)
		}
		expected := Array[int]{10, 20}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 10. Memory leak prevention test on old slice
	t.Run("memory leak zeroing", func(t *testing.T) {
		type Item struct{ Val int }
		i1 := &Item{Val: 1}
		i2 := &Item{Val: 2}
		i3 := &Item{Val: 3}

		a := Array[*Item]{i1, i2, i3}
		oldUnderlying := a[:3]

		removed := a.Splice(1, 1)
		if len(removed) != 1 || removed[0] != i2 {
			t.Fatalf("expected i2 removed")
		}
		if oldUnderlying[1] != nil {
			t.Errorf("expected oldUnderlying[1] to be zeroed out (nil), got %v", oldUnderlying[1])
		}
	})
}

func TestArray_Sort(t *testing.T) {
	// 1. Default string sorting
	t.Run("default string sorting", func(t *testing.T) {
		months := Array[string]{"March", "Jan", "Feb", "Dec"}
		ret := months.Sort()
		if ret != &months {
			t.Errorf("expected sort to return pointer to same array")
		}
		expected := Array[string]{"Dec", "Feb", "Jan", "March"}
		if !reflect.DeepEqual(months, expected) {
			t.Errorf("expected %v, got %v", expected, months)
		}
	})

	// 2. Default numeric sorting (converted to strings, classic JS behavior)
	t.Run("default numeric sorting (lexicographical)", func(t *testing.T) {
		numbers := Array[int]{1, 30, 4, 21, 100000}
		numbers.Sort()
		// In JS: "1" < "100000" < "21" < "30" < "4"
		expected := Array[int]{1, 100000, 21, 30, 4}
		if !reflect.DeepEqual(numbers, expected) {
			t.Errorf("expected %v, got %v", expected, numbers)
		}
	})

	// 3. Numeric ascending sorting with compareFn
	t.Run("numeric ascending with compareFn", func(t *testing.T) {
		numbers := Array[int]{1, 30, 4, 21, 100000}
		numbers.Sort(func(a, b int) int {
			return a - b
		})
		expected := Array[int]{1, 4, 21, 30, 100000}
		if !reflect.DeepEqual(numbers, expected) {
			t.Errorf("expected %v, got %v", expected, numbers)
		}
	})

	// 4. Numeric descending sorting with compareFn
	t.Run("numeric descending with compareFn", func(t *testing.T) {
		numbers := Array[int]{1, 30, 4, 21, 100000}
		numbers.Sort(func(a, b int) int {
			return b - a
		})
		expected := Array[int]{100000, 30, 21, 4, 1}
		if !reflect.DeepEqual(numbers, expected) {
			t.Errorf("expected %v, got %v", expected, numbers)
		}
	})

	// 5. In-place mutation and reference return
	t.Run("in-place mutation and chaining", func(t *testing.T) {
		a := Array[int]{3, 1, 2}
		a.Sort(func(x, y int) int { return x - y }).Push(4)
		expected := Array[int]{1, 2, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 6. Stability (ECMAScript 2019+ requirement)
	t.Run("stability", func(t *testing.T) {
		type Item struct {
			Key int
			Val string
		}
		items := Array[Item]{
			{Key: 2, Val: "first-2"},
			{Key: 1, Val: "first-1"},
			{Key: 2, Val: "second-2"},
			{Key: 1, Val: "second-1"},
		}
		items.Sort(func(a, b Item) int {
			return a.Key - b.Key
		})
		expected := Array[Item]{
			{Key: 1, Val: "first-1"},
			{Key: 1, Val: "second-1"},
			{Key: 2, Val: "first-2"},
			{Key: 2, Val: "second-2"},
		}
		if !reflect.DeepEqual(items, expected) {
			t.Errorf("expected %v, got %v", expected, items)
		}
	})

	// 7. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		ret := empty.Sort()
		if ret != &empty {
			t.Errorf("expected sort to return pointer to same array")
		}
		if len(empty) != 0 {
			t.Errorf("expected empty array, got %v", empty)
		}
	})

	// 8. Single element array
	t.Run("single element", func(t *testing.T) {
		single := Array[string]{"only"}
		single.Sort()
		expected := Array[string]{"only"}
		if !reflect.DeepEqual(single, expected) {
			t.Errorf("expected %v, got %v", expected, single)
		}
	})

	// 9. Nil array pointer
	t.Run("nil array pointer", func(t *testing.T) {
		var nilArr *Array[int]
		ret := nilArr.Sort()
		if ret != nil {
			t.Errorf("expected nil return for nil receiver, got %v", ret)
		}
	})

	// 10. Explicit nil compareFn fallback to default
	t.Run("nil compareFn fallback", func(t *testing.T) {
		a := Array[int]{40, 5}
		a.Sort(nil)
		// "40" < "5" lexicographically
		expected := Array[int]{40, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 11. Struct sorting with custom comparator
	t.Run("struct sorting", func(t *testing.T) {
		type Student struct {
			Name  string
			Grade int
		}
		students := Array[Student]{
			{"Charlie", 88},
			{"Alice", 95},
			{"Bob", 72},
		}
		students.Sort(func(a, b Student) int {
			if a.Grade < b.Grade {
				return -1
			} else if a.Grade > b.Grade {
				return 1
			}
			return 0
		})
		expected := Array[Student]{
			{"Bob", 72},
			{"Charlie", 88},
			{"Alice", 95},
		}
		if !reflect.DeepEqual(students, expected) {
			t.Errorf("expected %v, got %v", expected, students)
		}
	})
}

func TestArray_Reverse(t *testing.T) {
	// 1. Odd length array
	t.Run("odd length array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		ret := a.Reverse()
		if ret != &a {
			t.Errorf("expected reverse to return pointer to same array")
		}
		expected := Array[int]{3, 2, 1}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 2. Even length array
	t.Run("even length array", func(t *testing.T) {
		a := Array[string]{"one", "two", "three", "four"}
		a.Reverse()
		expected := Array[string]{"four", "three", "two", "one"}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 3. In-place mutation and method chaining
	t.Run("in-place mutation and chaining", func(t *testing.T) {
		a := Array[int]{3, 1, 2}
		a.Sort(func(x, y int) int { return x - y }).Reverse().Push(0)
		expected := Array[int]{3, 2, 1, 0}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 4. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		ret := empty.Reverse()
		if ret != &empty {
			t.Errorf("expected reverse to return pointer to same array")
		}
		if len(empty) != 0 {
			t.Errorf("expected empty array, got %v", empty)
		}
	})

	// 5. Single element array
	t.Run("single element", func(t *testing.T) {
		single := Array[int]{42}
		single.Reverse()
		expected := Array[int]{42}
		if !reflect.DeepEqual(single, expected) {
			t.Errorf("expected %v, got %v", expected, single)
		}
	})

	// 6. Nil array pointer
	t.Run("nil array pointer", func(t *testing.T) {
		var nilArr *Array[int]
		ret := nilArr.Reverse()
		if ret != nil {
			t.Errorf("expected nil return for nil receiver, got %v", ret)
		}
	})
}

func TestArray_Fill(t *testing.T) {
	// 1. Fill entire array
	t.Run("fill entire array", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		ret := a.Fill(0)
		if ret != &a {
			t.Errorf("expected fill to return pointer to same array")
		}
		expected := Array[int]{0, 0, 0, 0}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 2. Fill from start index
	t.Run("fill from start index", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.Fill(5, 1)
		expected := Array[int]{1, 5, 5, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 3. Fill from start to end index
	t.Run("fill from start to end", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.Fill(0, 2, 4)
		expected := Array[int]{1, 2, 0, 0}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 4. Negative start index
	t.Run("negative start index", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.Fill(6, -2)
		expected := Array[int]{1, 2, 6, 6}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 5. Negative start and end indices
	t.Run("negative start and end", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.Fill(7, -3, -1)
		expected := Array[int]{1, 7, 7, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 6. Start >= end (no-op)
	t.Run("start >= end", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.Fill(8, 3, 2)
		expected := Array[int]{1, 2, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 7. Negative start exceeding length
	t.Run("negative start exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.Fill(9, -10, 2)
		expected := Array[int]{9, 9, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 8. Start exceeding length
	t.Run("start exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.Fill(9, 5)
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 9. End exceeding length
	t.Run("end exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.Fill(9, 1, 10)
		expected := Array[int]{1, 9, 9}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 10. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		ret := empty.Fill(1)
		if ret != &empty {
			t.Errorf("expected fill to return pointer to same array")
		}
		if len(empty) != 0 {
			t.Errorf("expected empty array, got %v", empty)
		}
	})

	// 11. Nil array pointer
	t.Run("nil array pointer", func(t *testing.T) {
		var nilArr *Array[int]
		ret := nilArr.Fill(1)
		if ret != nil {
			t.Errorf("expected nil return for nil receiver, got %v", ret)
		}
	})

	// 12. In-place mutation and chaining
	t.Run("in-place mutation and chaining", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.Fill(0, 1).Push(4)
		expected := Array[int]{1, 0, 0, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})
}

func TestArray_CopyWithin(t *testing.T) {
	// 1. Target only (negative target)
	t.Run("target only with negative index", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		ret := a.CopyWithin(-2)
		if ret != &a {
			t.Errorf("expected copyWithin to return pointer to same array")
		}
		expected := Array[int]{1, 2, 3, 1, 2}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 2. Target and start
	t.Run("target and start", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		a.CopyWithin(0, 3)
		expected := Array[int]{4, 5, 3, 4, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 3. Target, start, and end
	t.Run("target, start, and end", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		a.CopyWithin(0, 3, 4)
		expected := Array[int]{4, 2, 3, 4, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 4. All negative indices
	t.Run("negative indices", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		a.CopyWithin(-2, -3, -1)
		expected := Array[int]{1, 2, 3, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 5. Overlapping copy shifting right
	t.Run("overlapping shift right", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		a.CopyWithin(1, 0, 3)
		expected := Array[int]{1, 1, 2, 3, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 6. Overlapping copy shifting left
	t.Run("overlapping shift left", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		a.CopyWithin(0, 1, 4)
		expected := Array[int]{2, 3, 4, 4, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 7. Target >= length (no-op)
	t.Run("target >= length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.CopyWithin(3)
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 8. Start >= end (no-op)
	t.Run("start >= end", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.CopyWithin(0, 2, 1)
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 9. Target < -length (clamped to 0)
	t.Run("target < -length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.CopyWithin(-10, 1)
		expected := Array[int]{2, 3, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 10. Start < -length (clamped to 0)
	t.Run("start < -length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.CopyWithin(1, -10)
		expected := Array[int]{1, 1, 2}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 11. End > length (clamped to length)
	t.Run("end > length", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.CopyWithin(0, 2, 10)
		expected := Array[int]{3, 4, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 12. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		ret := empty.CopyWithin(0, 1)
		if ret != &empty {
			t.Errorf("expected copyWithin to return pointer to same array")
		}
		if len(empty) != 0 {
			t.Errorf("expected empty array, got %v", empty)
		}
	})

	// 13. Nil array pointer
	t.Run("nil array pointer", func(t *testing.T) {
		var nilArr *Array[int]
		ret := nilArr.CopyWithin(0)
		if ret != nil {
			t.Errorf("expected nil return for nil receiver, got %v", ret)
		}
	})

	// 14. In-place mutation and chaining
	t.Run("in-place mutation and chaining", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		a.CopyWithin(0, 2).Push(5)
		expected := Array[int]{3, 4, 3, 4, 5}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})
}

func TestArray_Filter(t *testing.T) {
	// 1. Filter strings by length (MDN example)
	t.Run("filter words longer than 6", func(t *testing.T) {
		words := Array[string]{"spray", "elite", "exuberant", "destruction", "present"}
		result := words.Filter(func(word string) bool {
			return len(word) > 6
		})
		expected := Array[string]{"exuberant", "destruction", "present"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 2. Filter numbers by threshold (MDN example)
	t.Run("filter numbers >= 10", func(t *testing.T) {
		a := Array[int]{12, 5, 8, 130, 44}
		result := a.Filter(func(value int) bool {
			return value >= 10
		})
		expected := Array[int]{12, 130, 44}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 3. Filter primes (MDN example)
	t.Run("filter primes", func(t *testing.T) {
		a := Array[int]{-3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		result := a.Filter(func(num int) bool {
			for i := 2; num > i; i++ {
				if num%i == 0 {
					return false
				}
			}
			return num > 1
		})
		expected := Array[int]{2, 3, 5, 7, 11, 13}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 4. Original array is not mutated
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		_ = a.Filter(func(n int) bool {
			return n%2 == 0
		})
		expected := Array[int]{1, 2, 3, 4}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})

	// 5. No elements pass
	t.Run("no elements pass", func(t *testing.T) {
		a := Array[int]{1, 3, 5}
		result := a.Filter(func(n int) bool {
			return n%2 == 0
		})
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 6. All elements pass
	t.Run("all elements pass", func(t *testing.T) {
		a := Array[int]{2, 4, 6}
		result := a.Filter(func(n int) bool {
			return n%2 == 0
		})
		expected := Array[int]{2, 4, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 7. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Filter(func(_ int) bool {
			return true
		})
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 8. Nil array pointer will fail
	//t.Run("nil array pointer", func(t *testing.T) {
	//	var nilArr *Array[int]
	//	result := nilArr.Filter(func(_ int) bool {
	//		return true
	//	})
	//	if len(result) != 0 {
	//		t.Errorf("expected empty array for nil receiver, got %v", result)
	//	}
	//})

	// 9. Result is a shallow copy (distinct backing array)
	t.Run("result is a new array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Filter(func(_ int) bool {
			return true
		})
		result[0] = 99
		if a[0] != 1 {
			t.Errorf("mutating filter result should not change original, got %v", a)
		}
	})
}

func TestArray_Map(t *testing.T) {
	// 1. Map numbers by doubling (MDN example)
	t.Run("double numbers", func(t *testing.T) {
		a := Array[int]{1, 4, 9, 16}
		result := a.Map(func(x int) int {
			return x * 2
		})
		expected := Array[int]{2, 8, 18, 32}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 2. Map strings
	t.Run("uppercase strings", func(t *testing.T) {
		a := Array[string]{"hello", "world"}
		result := a.Map(func(s string) string {
			return s + "!"
		})
		expected := Array[string]{"hello!", "world!"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 3. Original array is not mutated
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Map(func(n int) int {
			return n * 10
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})

	// 4. Identity mapping
	t.Run("identity mapping", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Map(func(n int) int {
			return n
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 5. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Map(func(n int) int {
			return n * 2
		})
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 6. Nil array pointer will fail
	//t.Run("nil array pointer", func(t *testing.T) {
	//	var nilArr *Array[int]
	//	result := nilArr.Map(func(n int) int {
	//		return n
	//	})
	//	if len(result) != 0 {
	//		t.Errorf("expected empty array for nil receiver, got %v", result)
	//	}
	//})

	// 7. Result is a new array
	t.Run("result is a new array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Map(func(n int) int {
			return n
		})
		result[0] = 99
		if a[0] != 1 {
			t.Errorf("mutating map result should not change original, got %v", a)
		}
	})

	// 8. Result is a new type of the array
	t.Run("result is a new type array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Map(func(n int) string {
			return fmt.Sprintf("%d", n)
		})

		expected := Array[string]{"1", "2", "3"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}

	})
}

func TestArray_FlatMap(t *testing.T) {
	// 1. Map each number to [x, x*2] (MDN example)
	t.Run("map to pairs", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		result := a.FlatMap(func(x int) Array[int] {
			return Array[int]{x, x * 2}
		})
		expected := Array[int]{1, 2, 2, 4, 3, 6, 4, 8}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 2. Split strings (MDN example)
	t.Run("split strings", func(t *testing.T) {
		a := Array[string]{"it's Sunny in", "", "California"}
		result := a.FlatMap(func(s string) Array[string] {
			return Array[string](strings.Split(s, " "))
		})
		expected := Array[string]{"it's", "Sunny", "in", "", "California"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 3. Result is a new type of the array
	t.Run("result is a new type array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.FlatMap(func(n int) Array[string] {
			return Array[string]{fmt.Sprintf("%d", n), fmt.Sprintf("%d", n*10)}
		})
		expected := Array[string]{"1", "10", "2", "20", "3", "30"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 4. Empty inner arrays contribute nothing
	t.Run("empty inner arrays", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.FlatMap(func(n int) Array[int] {
			if n%2 == 0 {
				return Array[int]{}
			}
			return Array[int]{n}
		})
		expected := Array[int]{1, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 5. Original array is not mutated
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.FlatMap(func(n int) Array[int] {
			return Array[int]{n, n}
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})

	// 6. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.FlatMap(func(n int) Array[int] {
			return Array[int]{n}
		})
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 7. Nil array pointer will fail
	//t.Run("nil array pointer", func(t *testing.T) {
	//	var nilArr *Array[int]
	//	result := nilArr.FlatMap(func(n int) Array[string] {
	//		return Array[string]{fmt.Sprintf("%d", n)}
	//	})
	//	if len(result) != 0 {
	//		t.Errorf("expected empty array for nil receiver, got %v", result)
	//	}
	//})

	// 8. Result is a new array
	t.Run("result is a new array", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.FlatMap(func(n int) Array[int] {
			return Array[int]{n}
		})
		result[0] = 99
		if a[0] != 1 {
			t.Errorf("mutating flatMap result should not change original, got %v", a)
		}
	})
}

func TestArray_Reduce(t *testing.T) {
	// 1. Sum without initial value (MDN example)
	t.Run("sum without initial value", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		result := a.Reduce(func(acc, val int) int {
			return acc + val
		})
		expected := 10
		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	// 2. Sum with initial value (MDN example)
	t.Run("sum with initial value", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		result := a.Reduce(func(acc, val int) int {
			return acc + val
		}, 10)
		expected := 20
		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	// 3. Different accumulator type R (numbers to string)
	t.Run("different accumulator type", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Reduce(func(acc string, val int) string {
			return acc + fmt.Sprintf("%d", val)
		}, "")
		expected := "123"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	// 4. Count instances of values into map (MDN example)
	t.Run("count instances into map", func(t *testing.T) {
		names := Array[string]{"Alice", "Bob", "Tiff", "Bruce", "Alice"}
		counts := names.Reduce(func(acc map[string]int, name string) map[string]int {
			acc[name]++
			return acc
		}, make(map[string]int))

		expected := map[string]int{
			"Alice": 2,
			"Bob":   1,
			"Tiff":  1,
			"Bruce": 1,
		}
		if !reflect.DeepEqual(counts, expected) {
			t.Errorf("expected %v, got %v", expected, counts)
		}
	})

	// 5. Single element without initial value
	t.Run("single element without initial value", func(t *testing.T) {
		a := Array[int]{42}
		called := false
		result := a.Reduce(func(acc, val int) int {
			called = true
			return acc + val
		})
		if called {
			t.Errorf("callback should not be called for single element without initial value")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	// 6. Single element with initial value
	t.Run("single element with initial value", func(t *testing.T) {
		a := Array[int]{42}
		result := a.Reduce(func(acc, val int) int {
			return acc + val
		}, 8)
		if result != 50 {
			t.Errorf("expected 50, got %d", result)
		}
	})

	// 7. Empty array with initial value
	t.Run("empty array with initial value", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Reduce(func(acc, val int) int {
			return acc + val
		}, 99)
		if result != 99 {
			t.Errorf("expected 99, got %d", result)
		}
	})

	// 8. Empty array without initial value panics
	t.Run("empty array without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic for empty array without initial value")
			}
		}()
		empty := Array[int]{}
		_ = empty.Reduce(func(acc, val int) int {
			return acc + val
		})
	})

	//// 9. Nil array pointer with initial value
	//t.Run("nil array pointer with initial value", func(t *testing.T) {
	//	var nilArr *Array[int]
	//	result := nilArr.Reduce(func(acc, val int) int {
	//		return acc + val
	//	}, 99)
	//	if result != 99 {
	//		t.Errorf("expected 99, got %d", result)
	//	}
	//})

	// 10. Nil array pointer without initial value panics
	t.Run("nil array pointer without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic for nil array without initial value")
			}
		}()
		var nilArr *Array[int]
		_ = nilArr.Reduce(func(acc, val int) int {
			return acc + val
		})
	})

	// 11. Incompatible accumulator type without initial value panics
	t.Run("incompatible accumulator type without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic when accumulator type cannot be initialized from element")
			}
		}()
		a := Array[int]{1, 2, 3}
		_ = a.Reduce(func(acc string, val int) string {
			return acc + fmt.Sprintf("%d", val)
		})
	})

	// 12. Does not mutate original array
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Reduce(func(acc, val int) int {
			return acc + val
		}, 0)
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})
}

func TestArray_ReduceRight(t *testing.T) {
	// 1. String concatenation right-to-left (MDN example)
	t.Run("string concatenation right-to-left", func(t *testing.T) {
		a := Array[string]{"1", "2", "3", "4", "5"}
		result := a.ReduceRight(func(prev, cur string) string {
			return prev + cur
		})
		expected := "54321"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	// 2. Subtraction right-to-left without initial value
	t.Run("subtraction without initial value", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		// 4 - 3 - 2 - 1 = -2
		result := a.ReduceRight(func(acc, cur int) int {
			return acc - cur
		})
		expected := -2
		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	// 3. Subtraction right-to-left with initial value
	t.Run("subtraction with initial value", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		// 10 - 4 - 3 - 2 - 1 = 0
		result := a.ReduceRight(func(acc, cur int) int {
			return acc - cur
		}, 10)
		expected := 0
		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	// 4. Flattening array of arrays right-to-left (MDN example)
	t.Run("flattening right-to-left", func(t *testing.T) {
		a := Array[Array[int]]{{0, 1}, {2, 3}, {4, 5}}
		result := a.ReduceRight(func(prev, cur Array[int]) Array[int] {
			return append(prev, cur...)
		}, Array[int]{})
		expected := Array[int]{4, 5, 2, 3, 0, 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 5. Different accumulator type R (numbers to string in reverse)
	t.Run("different accumulator type", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.ReduceRight(func(acc string, val int) string {
			return acc + fmt.Sprintf("%d", val)
		}, "")
		expected := "321"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	// 6. Single element without initial value
	t.Run("single element without initial value", func(t *testing.T) {
		a := Array[int]{42}
		called := false
		result := a.ReduceRight(func(acc, val int) int {
			called = true
			return acc + val
		})
		if called {
			t.Errorf("callback should not be called for single element without initial value")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	// 7. Single element with initial value
	t.Run("single element with initial value", func(t *testing.T) {
		a := Array[int]{42}
		result := a.ReduceRight(func(acc, val int) int {
			return acc + val
		}, 8)
		if result != 50 {
			t.Errorf("expected 50, got %d", result)
		}
	})

	// 8. Empty array with initial value
	t.Run("empty array with initial value", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.ReduceRight(func(acc, val int) int {
			return acc + val
		}, 99)
		if result != 99 {
			t.Errorf("expected 99, got %d", result)
		}
	})

	// 9. Empty array without initial value panics
	t.Run("empty array without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic for empty array without initial value")
			}
		}()
		empty := Array[int]{}
		_ = empty.ReduceRight(func(acc, val int) int {
			return acc + val
		})
	})

	// 10. Nil array without initial value panics
	t.Run("nil array without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic for nil array without initial value")
			}
		}()
		var nilArr Array[int]
		_ = nilArr.ReduceRight(func(acc, val int) int {
			return acc + val
		})
	})

	// 11. Incompatible accumulator type without initial value panics
	t.Run("incompatible accumulator type without initial value panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic when accumulator type cannot be initialized from element")
			}
		}()
		a := Array[int]{1, 2, 3}
		_ = a.ReduceRight(func(acc string, val int) string {
			return acc + fmt.Sprintf("%d", val)
		})
	})

	// 12. Does not mutate original array
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.ReduceRight(func(acc, val int) int {
			return acc + val
		}, 0)
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})
}

func TestArray_ForEach(t *testing.T) {
	// 1. Iterates over elements in order
	t.Run("iterates in order", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		var collected []string
		a.ForEach(func(s string) {
			collected = append(collected, s)
		})
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(collected, expected) {
			t.Errorf("expected %v, got %v", expected, collected)
		}
	})

	// 2. Side effects (accumulation)
	t.Run("accumulate sum", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4}
		sum := 0
		a.ForEach(func(n int) {
			sum += n
		})
		if sum != 10 {
			t.Errorf("expected sum 10, got %d", sum)
		}
	})

	// 3. Does not mutate original array
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.ForEach(func(n int) {
			n = n * 2
			_ = n
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected %v, got %v", expected, a)
		}
	})

	// 4. Works in method chain
	t.Run("chaining with other methods", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5, 6}
		var evens []int
		a.Filter(func(n int) bool {
			return n%2 == 0
		}).ForEach(func(n int) {
			evens = append(evens, n)
		})
		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(evens, expected) {
			t.Errorf("expected %v, got %v", expected, evens)
		}
	})

	// 5. Empty array
	t.Run("empty array", func(t *testing.T) {
		called := false
		empty := Array[int]{}
		empty.ForEach(func(n int) {
			called = true
		})
		if called {
			t.Errorf("callback should not be called for empty array")
		}
	})

	// 6. Nil array
	t.Run("nil array", func(t *testing.T) {
		called := false
		var nilArr Array[int]
		nilArr.ForEach(func(n int) {
			called = true
		})
		if called {
			t.Errorf("callback should not be called for nil array")
		}
	})

	// 7. Nil callback
	t.Run("nil callback does not panic", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		a.ForEach(nil)
	})
}

func TestArray_Every(t *testing.T) {
	// 1. All elements pass (MDN example)
	t.Run("all elements pass threshold", func(t *testing.T) {
		a := Array[int]{1, 30, 39, 29, 10, 13}
		result := a.Every(func(v int) bool {
			return v < 40
		})
		if !result {
			t.Errorf("expected true, got false")
		}
	})

	// 2. One element fails (MDN example)
	t.Run("one element exceeds threshold", func(t *testing.T) {
		a := Array[int]{1, 30, 39, 29, 10, 43}
		result := a.Every(func(v int) bool {
			return v < 40
		})
		if result {
			t.Errorf("expected false, got true")
		}
	})

	// 3. Short-circuiting verification
	t.Run("short circuits on first failure", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		calls := 0
		result := a.Every(func(v int) bool {
			calls++
			return v < 3
		})
		if result {
			t.Errorf("expected false, got true")
		}
		if calls != 3 {
			t.Errorf("expected 3 calls, got %d", calls)
		}
	})

	// 4. Empty array returns true (vacuous truth)
	t.Run("empty array returns true", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Every(func(v int) bool {
			return v > 100
		})
		if !result {
			t.Errorf("expected true for empty array, got false")
		}
	})

	// 5. Nil array returns true
	t.Run("nil array returns true", func(t *testing.T) {
		var nilArr Array[int]
		result := nilArr.Every(func(v int) bool {
			return false
		})
		if !result {
			t.Errorf("expected true for nil array, got false")
		}
	})

	// 6. Does not mutate original array
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Every(func(v int) bool {
			return v > 0
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})

	// 7. Chaining with other methods
	t.Run("chaining with Filter", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5, 6}
		result := a.Filter(func(v int) bool {
			return v%2 == 0
		}).Every(func(v int) bool {
			return v >= 2
		})
		if !result {
			t.Errorf("expected true, got false")
		}
	})
}

func TestArray_Some(t *testing.T) {
	// 1. At least one element passes (MDN example)
	t.Run("even number exists", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		result := a.Some(func(v int) bool {
			return v%2 == 0
		})
		if !result {
			t.Errorf("expected true, got false")
		}
	})

	// 2. No elements pass
	t.Run("no even number", func(t *testing.T) {
		a := Array[int]{1, 3, 5}
		result := a.Some(func(v int) bool {
			return v%2 == 0
		})
		if result {
			t.Errorf("expected false, got true")
		}
	})

	// 3. Short-circuiting verification
	t.Run("short circuits on first success", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		calls := 0
		result := a.Some(func(v int) bool {
			calls++
			return v == 2
		})
		if !result {
			t.Errorf("expected true, got false")
		}
		if calls != 2 {
			t.Errorf("expected 2 calls, got %d", calls)
		}
	})

	// 4. Empty array returns false
	t.Run("empty array returns false", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Some(func(v int) bool {
			return true
		})
		if result {
			t.Errorf("expected false for empty array, got true")
		}
	})

	// 5. Nil array returns false
	t.Run("nil array returns false", func(t *testing.T) {
		var nilArr Array[int]
		result := nilArr.Some(func(v int) bool {
			return true
		})
		if result {
			t.Errorf("expected false for nil array, got true")
		}
	})

	// 6. Does not mutate original array
	t.Run("does not mutate original", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Some(func(v int) bool {
			return v == 2
		})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(a, expected) {
			t.Errorf("expected original %v, got %v", expected, a)
		}
	})

	// 7. Chaining with other methods
	t.Run("chaining with Map", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Map(func(v int) int {
			return v * 10
		}).Some(func(v int) bool {
			return v == 20
		})
		if !result {
			t.Errorf("expected true, got false")
		}
	})
}

func TestArray_Concat(t *testing.T) {
	// 1. Concat two arrays (MDN example)
	t.Run("concat two arrays", func(t *testing.T) {
		letters := Array[string]{"a", "b", "c"}
		numbers := Array[string]{"d", "e", "f"}
		result := letters.Concat(numbers)
		expected := Array[string]{"a", "b", "c", "d", "e", "f"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 2. Concat multiple arrays (MDN example)
	t.Run("concat three arrays", func(t *testing.T) {
		num1 := Array[int]{1, 2, 3}
		num2 := Array[int]{4, 5, 6}
		num3 := Array[int]{7, 8, 9}
		result := num1.Concat(num2, num3)
		expected := Array[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 3. Concat with no arguments returns a shallow copy
	t.Run("no arguments returns copy", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Concat()
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
		result[0] = 99
		if a[0] != 1 {
			t.Errorf("mutating copy should not change original, got %v", a)
		}
	})

	// 4. Original arrays are not mutated
	t.Run("does not mutate originals", func(t *testing.T) {
		a := Array[int]{1, 2}
		b := Array[int]{3, 4}
		_ = a.Concat(b)
		if !reflect.DeepEqual(a, Array[int]{1, 2}) || !reflect.DeepEqual(b, Array[int]{3, 4}) {
			t.Errorf("expected originals to be unmodified")
		}
	})

	// 5. Empty arrays in arguments
	t.Run("empty arrays in arguments", func(t *testing.T) {
		a := Array[int]{1, 2}
		result := a.Concat(Array[int]{}, Array[int]{3})
		expected := Array[int]{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 6. Empty receiver
	t.Run("empty receiver", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Concat(Array[int]{1, 2})
		expected := Array[int]{1, 2}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 7. Nil receiver
	t.Run("nil receiver", func(t *testing.T) {
		var nilArr Array[int]
		result := nilArr.Concat(Array[int]{1, 2})
		expected := Array[int]{1, 2}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 8. Nil arguments
	t.Run("nil arguments", func(t *testing.T) {
		a := Array[int]{1, 2}
		var nilArg Array[int]
		result := a.Concat(nilArg)
		expected := Array[int]{1, 2}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 9. Chaining with other methods
	t.Run("chaining with Filter", func(t *testing.T) {
		a := Array[int]{1, 2}
		b := Array[int]{3, 4}
		result := a.Concat(b).Filter(func(n int) bool {
			return n%2 == 0
		})
		expected := Array[int]{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestArray_Slice(t *testing.T) {
	animals := Array[string]{"ant", "bison", "camel", "duck", "elephant"}

	// 1. Slice from start (MDN example)
	t.Run("from index 2", func(t *testing.T) {
		result := animals.Slice(2)
		expected := Array[string]{"camel", "duck", "elephant"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 2. Slice from start to end (MDN example)
	t.Run("from 2 to 4", func(t *testing.T) {
		result := animals.Slice(2, 4)
		expected := Array[string]{"camel", "duck"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 3. Slice from 1 to 5 (MDN example)
	t.Run("from 1 to 5", func(t *testing.T) {
		result := animals.Slice(1, 5)
		expected := Array[string]{"bison", "camel", "duck", "elephant"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 4. Slice with negative start (MDN example)
	t.Run("negative start", func(t *testing.T) {
		result := animals.Slice(-2)
		expected := Array[string]{"duck", "elephant"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 5. Slice with positive start and negative end (MDN example)
	t.Run("positive start and negative end", func(t *testing.T) {
		result := animals.Slice(2, -1)
		expected := Array[string]{"camel", "duck"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 6. Slice with no arguments (shallow copy)
	t.Run("no arguments returns copy", func(t *testing.T) {
		result := animals.Slice()
		expected := Array[string]{"ant", "bison", "camel", "duck", "elephant"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
		result[0] = "zebra"
		if animals[0] != "ant" {
			t.Errorf("mutating copy should not modify original, got %v", animals)
		}
	})

	// 7. Start >= length
	t.Run("start exceeding length", func(t *testing.T) {
		result := animals.Slice(10)
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 8. Start >= end
	t.Run("start >= end", func(t *testing.T) {
		result := animals.Slice(3, 2)
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 9. Negative start exceeding length
	t.Run("negative start exceeding length", func(t *testing.T) {
		result := animals.Slice(-10, 2)
		expected := Array[string]{"ant", "bison"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 10. Negative start and end exceeding length
	t.Run("both negative exceeding length", func(t *testing.T) {
		result := animals.Slice(-10, -8)
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 11. End exceeding length
	t.Run("end exceeding length", func(t *testing.T) {
		result := animals.Slice(3, 10)
		expected := Array[string]{"duck", "elephant"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	// 12. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Slice(0, 1)
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 13. Nil array
	t.Run("nil array", func(t *testing.T) {
		var nilArr Array[int]
		result := nilArr.Slice(0, 1)
		if len(result) != 0 {
			t.Errorf("expected empty array, got %v", result)
		}
	})

	// 14. Chaining with other methods
	t.Run("chaining with Map", func(t *testing.T) {
		a := Array[int]{1, 2, 3, 4, 5}
		result := a.Slice(1, 4).Map(func(n int) int {
			return n * 10
		})
		expected := Array[int]{20, 30, 40}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestArray_Flat(t *testing.T) {

	t.Run("data", func(t *testing.T) {
		nested := Array[any]{1, []any{2, 3}, []any{4, []any{5, []any{6}}}}

		tests := []struct {
			name  string
			input Array[any]
			depth []int
			want  Array[any]
		}{
			{
				name:  "default depth is 1",
				input: nested,
				want:  Array[any]{1, 2, 3, 4, []any{5, []any{6}}},
			},
			{
				name:  "depth 2",
				input: nested,
				depth: []int{2},
				want:  Array[any]{1, 2, 3, 4, 5, []any{6}},
			},
			{
				name:  "Infinity flattens everything",
				input: nested,
				depth: []int{-1},
				want:  Array[any]{1, 2, 3, 4, 5, 6},
			},
			{
				name:  "depth 0 does not flatten",
				input: nested,
				depth: []int{0},
				want:  Array[any]{1, []any{2, 3}, []any{4, []any{5, []any{6}}}},
			},
			{
				name:  "empty array",
				input: Array[any]{},
				want:  Array[any]{},
			},
			{
				name:  "empty nested arrays vanish",
				input: Array[any]{[]any{}, 1, []any{[]any{}}},
				depth: []int{-1},
				want:  Array[any]{1},
			},
			{
				name:  "nil elements are preserved",
				input: Array[any]{nil, []any{nil, 1}},
				want:  Array[any]{nil, nil, 1},
			},
			{
				name:  "strings are not flattened",
				input: Array[any]{"ab", []any{"cd"}},
				want:  Array[any]{"ab", "cd"},
			},
			{
				name:  "nested Array type is flattened",
				input: Array[any]{Array[any]{1, 2}, Array[int]{3}},
				want:  Array[any]{1, 2, 3},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.input.Flat(tt.depth...)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Flat(%v) = %#v, want %#v", tt.depth, got, tt.want)
				}
			})
		}
	})

	t.Run("typed slices", func(t *testing.T) {
		a := Array[[]int]{{1, 2}, {3}, {}}
		got := a.Flat()
		want := Array[any]{1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("does not mutate", func(t *testing.T) {
		a := Array[any]{1, []any{2, 3}}
		_ = a.Flat(-1)
		want := Array[any]{1, []any{2, 3}}
		if !reflect.DeepEqual(a, want) {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_ToSorted(t *testing.T) {
	tests := []struct {
		name    string
		input   Array[int]
		compare func(a, b int) int
		want    Array[int]
	}{
		{
			name:    "ascending",
			input:   Array[int]{3, 1, 2},
			compare: cmp.Compare[int],
			want:    Array[int]{1, 2, 3},
		},
		{
			name:    "descending",
			input:   Array[int]{3, 1, 2},
			compare: func(a, b int) int { return cmp.Compare(b, a) },
			want:    Array[int]{3, 2, 1},
		},
		{
			name:    "nil comparator sorts by string form like JS",
			input:   Array[int]{10, 9, 1, 100},
			compare: nil,
			want:    Array[int]{1, 10, 100, 9},
		},
		{
			name:    "empty array",
			input:   Array[int]{},
			compare: cmp.Compare[int],
			want:    Array[int]{},
		},
		{
			name:    "single element",
			input:   Array[int]{7},
			compare: cmp.Compare[int],
			want:    Array[int]{7},
		},
		{
			name:    "duplicates",
			input:   Array[int]{2, 1, 2, 1},
			compare: cmp.Compare[int],
			want:    Array[int]{1, 1, 2, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToSorted(tt.compare)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSorted() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestArray_ToSortedIsStable(t *testing.T) {
	type pair struct {
		key int
		tag string
	}
	a := Array[pair]{{1, "a"}, {0, "b"}, {1, "c"}, {0, "d"}}
	got := a.ToSorted(func(x, y pair) int { return cmp.Compare(x.key, y.key) })
	want := Array[pair]{{0, "b"}, {0, "d"}, {1, "a"}, {1, "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not stable: got %#v, want %#v", got, want)
	}
}

func TestArray_ToSortedDoesNotMutateReceiver(t *testing.T) {
	a := Array[int]{3, 1, 2}
	got := a.ToSorted(cmp.Compare[int])
	if want := (Array[int]{3, 1, 2}); !reflect.DeepEqual(a, want) {
		t.Errorf("receiver mutated: %#v", a)
	}
	got[0] = 99
	if a[0] == 99 {
		t.Error("result shares memory with receiver")
	}
}

func TestArray_ToReversed(t *testing.T) {
	tests := []struct {
		name  string
		input Array[int]
		want  Array[int]
	}{
		{"odd length", Array[int]{1, 2, 3}, Array[int]{3, 2, 1}},
		{"even length", Array[int]{1, 2, 3, 4}, Array[int]{4, 3, 2, 1}},
		{"single element", Array[int]{1}, Array[int]{1}},
		{"empty array", Array[int]{}, Array[int]{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToReversed()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToReversed() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestArray_ToReversedDoesNotMutateReceiver(t *testing.T) {
	a := Array[int]{1, 2, 3}
	got := a.ToReversed()
	if want := (Array[int]{1, 2, 3}); !reflect.DeepEqual(a, want) {
		t.Errorf("receiver mutated: %#v", a)
	}
	got[0] = 99
	if a[2] == 99 {
		t.Error("result shares memory with receiver")
	}
}

func TestArray_ToSpliced(t *testing.T) {
	tests := []struct {
		name        string
		input       Array[int]
		start       int
		deleteCount int
		items       []int
		want        Array[int]
	}{
		{"insert only", Array[int]{1, 2, 3}, 1, 0, []int{9}, Array[int]{1, 9, 2, 3}},
		{"delete only", Array[int]{1, 2, 3}, 1, 1, nil, Array[int]{1, 3}},
		{"replace", Array[int]{1, 2, 3}, 1, 1, []int{9, 8}, Array[int]{1, 9, 8, 3}},
		{"replace with fewer items", Array[int]{1, 2, 3, 4}, 0, 3, []int{9}, Array[int]{9, 4}},
		{"negative start counts from end", Array[int]{1, 2, 3}, -1, 1, nil, Array[int]{1, 2}},
		{"negative start before beginning is clamped to 0", Array[int]{1, 2, 3}, -10, 1, nil, Array[int]{2, 3}},
		{"start past end is clamped to length", Array[int]{1, 2, 3}, 10, 0, []int{4}, Array[int]{1, 2, 3, 4}},
		{"deleteCount past end is clamped", Array[int]{1, 2, 3}, 1, 100, nil, Array[int]{1}},
		{"MaxInt deletes to the end", Array[int]{1, 2, 3}, 1, math.MaxInt, nil, Array[int]{1}},
		{"negative deleteCount is treated as 0", Array[int]{1, 2, 3}, 1, -5, []int{9}, Array[int]{1, 9, 2, 3}},
		{"insert into empty array", Array[int]{}, 0, 0, []int{1, 2}, Array[int]{1, 2}},
		{"delete from empty array", Array[int]{}, 0, 3, nil, Array[int]{}},
		{"nothing to do returns a copy", Array[int]{1, 2, 3}, 0, 0, nil, Array[int]{1, 2, 3}},
		{"most negative start", Array[int]{1, 2, 3}, math.MinInt, 1, nil, Array[int]{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToSpliced(tt.start, tt.deleteCount, tt.items...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSpliced(%d, %d, %v) = %#v, want %#v",
					tt.start, tt.deleteCount, tt.items, got, tt.want)
			}
		})
	}
}

func TestToSplicedDoesNotMutateReceiver(t *testing.T) {
	a := Array[int]{1, 2, 3}
	_ = a.ToSpliced(1, 1, 9, 9)
	if want := (Array[int]{1, 2, 3}); !reflect.DeepEqual(a, want) {
		t.Errorf("receiver mutated: %#v", a)
	}

	// Even when nothing changes, the result must not share memory.
	same := a.ToSpliced(0, 0)
	same[0] = 99
	if a[0] == 99 {
		t.Error("result shares memory with receiver")
	}
}

func TestArray_With(t *testing.T) {
	tests := []struct {
		name  string
		input Array[int]
		index int
		value int
		want  Array[int]
	}{
		{"middle", Array[int]{1, 2, 3}, 1, 9, Array[int]{1, 9, 3}},
		{"first", Array[int]{1, 2, 3}, 0, 9, Array[int]{9, 2, 3}},
		{"last", Array[int]{1, 2, 3}, 2, 9, Array[int]{1, 2, 9}},
		{"negative index counts from end", Array[int]{1, 2, 3}, -1, 9, Array[int]{1, 2, 9}},
		{"most negative valid index", Array[int]{1, 2, 3}, -3, 9, Array[int]{9, 2, 3}},
		{"single element", Array[int]{1}, 0, 9, Array[int]{9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.With(tt.index, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("With(%d, %d) = %#v, want %#v", tt.index, tt.value, got, tt.want)
			}
		})
	}
}

func TestWithPanicsOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		input Array[int]
		index int
	}{
		{"index equal to length", Array[int]{1, 2, 3}, 3},
		{"index past length", Array[int]{1, 2, 3}, 10},
		{"negative index before start", Array[int]{1, 2, 3}, -4},
		{"most negative index", Array[int]{1, 2, 3}, math.MinInt},
		{"empty array", Array[int]{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("With(%d, 0) on length %d did not panic", tt.index, len(tt.input))
				}
			}()
			tt.input.With(tt.index, 0)
		})
	}
}

func TestWithDoesNotMutateReceiver(t *testing.T) {
	a := Array[int]{1, 2, 3}
	got := a.With(1, 9)
	if want := (Array[int]{1, 2, 3}); !reflect.DeepEqual(a, want) {
		t.Errorf("receiver mutated: %#v", a)
	}
	got[0] = 99
	if a[0] == 99 {
		t.Error("result shares memory with receiver")
	}
}

func TestArray_Find(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name      string
		input     Array[int]
		predicate func(int) bool
		want      int
		wantOK    bool
	}{
		{"first match is returned", Array[int]{1, 3, 4, 6}, isEven, 4, true},
		{"match at start", Array[int]{2, 3, 5}, isEven, 2, true},
		{"match at end", Array[int]{1, 3, 8}, isEven, 8, true},
		{"no match returns zero value and false", Array[int]{1, 3, 5}, isEven, 0, false},
		{"empty array", Array[int]{}, isEven, 0, false},
		{"nil array", nil, isEven, 0, false},
		{"zero value can be a real match", Array[int]{1, 0, 2}, func(v int) bool { return v == 0 }, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.input.Find(tt.predicate)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("Find() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}

	t.Run("stops at first match", func(t *testing.T) {
		calls := 0
		a := Array[int]{1, 2, 3, 4}
		got, ok := a.Find(func(v int) bool {
			calls++
			return v == 2
		})
		if !ok || got != 2 {
			t.Fatalf("Find() = (%v, %v), want (2, true)", got, ok)
		}
		if calls != 2 {
			t.Errorf("predicate called %d times, want 2", calls)
		}
	})

	t.Run("structs", func(t *testing.T) {
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got, ok := users.Find(func(u user) bool { return u.age == 25 })
		if !ok || got.name != "bob" {
			t.Errorf("Find() = (%v, %v), want bob", got, ok)
		}
	})

	t.Run("pointers return the same pointer", func(t *testing.T) {
		// Like JS returning the same object reference.
		a, b := &user{"ann", 30}, &user{"bob", 25}
		ptrs := Array[*user]{a, b}
		got, ok := ptrs.Find(func(u *user) bool { return u.name == "bob" })
		if !ok || got != b {
			t.Errorf("Find() returned %p, want %p", got, b)
		}
	})

	t.Run("pointers not found yields nil", func(t *testing.T) {
		ptrs := Array[*user]{{"ann", 30}, {"bob", 25}}
		got, ok := ptrs.Find(func(u *user) bool { return u.name == "zed" })
		if ok || got != nil {
			t.Errorf("Find() = (%v, %v), want (nil, false)", got, ok)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_, _ = a.Find(func(v int) bool { return v == 2 })
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_FindLast(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name      string
		input     Array[int]
		predicate func(int) bool
		want      int
		wantOK    bool
	}{
		{"last match is returned", Array[int]{1, 2, 3, 4, 5}, isEven, 4, true},
		{"match at start", Array[int]{2, 3, 5}, isEven, 2, true},
		{"match at end", Array[int]{1, 3, 8}, isEven, 8, true},
		{"single match in the middle", Array[int]{1, 4, 5}, isEven, 4, true},
		{"no match returns zero value and false", Array[int]{1, 3, 5}, isEven, 0, false},
		{"empty array", Array[int]{}, isEven, 0, false},
		{"nil array", nil, isEven, 0, false},
		{"zero value can be a real match", Array[int]{1, 0, 2}, func(v int) bool { return v == 0 }, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.input.FindLast(tt.predicate)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("FindLast() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}

	t.Run("searches from the end and stops at first match", func(t *testing.T) {
		var visited []int
		a := Array[int]{1, 2, 3, 4}
		got, ok := a.FindLast(func(v int) bool {
			visited = append(visited, v)
			return v == 3
		})
		if !ok || got != 3 {
			t.Fatalf("FindLast() = (%v, %v), want (3, true)", got, ok)
		}
		if len(visited) != 2 || visited[0] != 4 || visited[1] != 3 {
			t.Errorf("visited %v, want [4 3]", visited)
		}
	})

	t.Run("differs from Find when several elements match", func(t *testing.T) {
		a := Array[int]{2, 4, 6}
		first, _ := a.Find(isEven)
		last, _ := a.FindLast(isEven)
		if first != 2 || last != 6 {
			t.Errorf("Find = %d, FindLast = %d, want 2 and 6", first, last)
		}
	})

	t.Run("structs", func(t *testing.T) {
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got, ok := users.FindLast(func(u user) bool { return u.age == 25 })
		if !ok || got.name != "cy" {
			t.Errorf("FindLast() = (%v, %v), want cy", got, ok)
		}
	})

	t.Run("pointers return the same pointer", func(t *testing.T) {
		a, b, c := &user{"ann", 25}, &user{"bob", 30}, &user{"cy", 25}
		ptrs := Array[*user]{a, b, c}
		got, ok := ptrs.FindLast(func(u *user) bool { return u.age == 25 })
		if !ok || got != c {
			t.Errorf("FindLast() returned %p, want %p", got, c)
		}
	})

	t.Run("pointers not found yields nil", func(t *testing.T) {
		ptrs := Array[*user]{{"ann", 30}, {"bob", 25}}
		got, ok := ptrs.FindLast(func(u *user) bool { return u.name == "zed" })
		if ok || got != nil {
			t.Errorf("FindLast() = (%v, %v), want (nil, false)", got, ok)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_, _ = a.FindLast(func(v int) bool { return v == 2 })
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_FindIndex(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name      string
		input     Array[int]
		predicate func(int) bool
		want      int
	}{
		{"first match is returned", Array[int]{1, 3, 4, 6}, isEven, 2},
		{"match at start", Array[int]{2, 3, 5}, isEven, 0},
		{"match at end", Array[int]{1, 3, 8}, isEven, 2},
		{"no match returns -1", Array[int]{1, 3, 5}, isEven, -1},
		{"empty array", Array[int]{}, isEven, -1},
		{"nil array", nil, isEven, -1},
		{"zero value can be a real match", Array[int]{1, 0, 2}, func(v int) bool { return v == 0 }, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.FindIndex(tt.predicate)
			if got != tt.want {
				t.Errorf("FindIndex() = %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("stops at first match", func(t *testing.T) {
		calls := 0
		a := Array[int]{1, 2, 3, 4}
		got := a.FindIndex(func(v int) bool {
			calls++
			return v == 2
		})
		if got != 1 {
			t.Fatalf("FindIndex() = %d, want 1", got)
		}
		if calls != 2 {
			t.Errorf("predicate called %d times, want 2", calls)
		}
	})

	t.Run("structs", func(t *testing.T) {
		type user struct {
			name string
			age  int
		}
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got := users.FindIndex(func(u user) bool { return u.age == 25 })
		if got != 1 {
			t.Errorf("FindIndex() = %d, want 1", got)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.FindIndex(func(v int) bool { return v == 2 })
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_FindLastIndex(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name      string
		input     Array[int]
		predicate func(int) bool
		want      int
	}{
		{"last match is returned", Array[int]{1, 2, 3, 4, 5}, isEven, 3},
		{"match at start", Array[int]{2, 3, 5}, isEven, 0},
		{"match at end", Array[int]{1, 3, 8}, isEven, 2},
		{"single match in the middle", Array[int]{1, 4, 5}, isEven, 1},
		{"no match returns -1", Array[int]{1, 3, 5}, isEven, -1},
		{"empty array", Array[int]{}, isEven, -1},
		{"nil array", nil, isEven, -1},
		{"zero value can be a real match", Array[int]{1, 0, 2}, func(v int) bool { return v == 0 }, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.FindLastIndex(tt.predicate)
			if got != tt.want {
				t.Errorf("FindLastIndex() = %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("searches from the end and stops at first match", func(t *testing.T) {
		var visited []int
		a := Array[int]{1, 2, 3, 4}
		got := a.FindLastIndex(func(v int) bool {
			visited = append(visited, v)
			return v == 3
		})
		if got != 2 {
			t.Fatalf("FindLastIndex() = %d, want 2", got)
		}
		if len(visited) != 2 || visited[0] != 4 || visited[1] != 3 {
			t.Errorf("visited %v, want [4 3]", visited)
		}
	})

	t.Run("differs from FindIndex when several elements match", func(t *testing.T) {
		a := Array[int]{2, 4, 6}
		first := a.FindIndex(isEven)
		last := a.FindLastIndex(isEven)
		if first != 0 || last != 2 {
			t.Errorf("FindIndex = %d, FindLastIndex = %d, want 0 and 2", first, last)
		}
	})

	t.Run("structs", func(t *testing.T) {
		type user struct {
			name string
			age  int
		}
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got := users.FindLastIndex(func(u user) bool { return u.age == 25 })
		if got != 2 {
			t.Errorf("FindLastIndex() = %d, want 2", got)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.FindLastIndex(func(v int) bool { return v == 2 })
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_Includes(t *testing.T) {
	// 1. Basic inclusion (MDN examples)
	t.Run("basic inclusion", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		if !a.Includes(2) {
			t.Errorf("expected true, got false")
		}
		if a.Includes(4) {
			t.Errorf("expected false, got true")
		}
	})

	// 2. fromIndex positive
	t.Run("fromIndex positive", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		if !a.Includes(2, 1) {
			t.Errorf("expected true, got false")
		}
		if a.Includes(2, 2) {
			t.Errorf("expected false, got true")
		}
	})

	// 3. fromIndex negative (counts from end)
	t.Run("fromIndex negative", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		if !a.Includes(3, -1) {
			t.Errorf("expected true, got false")
		}
		if !a.Includes(2, -2) {
			t.Errorf("expected true, got false")
		}
		if a.Includes(1, -2) {
			t.Errorf("expected false, got true")
		}
	})

	// 4. fromIndex greater than length
	t.Run("fromIndex >= length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		if a.Includes(1, 3) {
			t.Errorf("expected false, got true")
		}
		if a.Includes(1, 100) {
			t.Errorf("expected false, got true")
		}
	})

	// 5. fromIndex negative exceeding length
	t.Run("negative fromIndex exceeding length", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		if !a.Includes(1, -10) {
			t.Errorf("expected true, got false")
		}
	})

	// 6. Zero value match
	t.Run("zero value match", func(t *testing.T) {
		a := Array[int]{1, 0, 2}
		if !a.Includes(0) {
			t.Errorf("expected true, got false")
		}
	})

	// 7. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		if empty.Includes(1) {
			t.Errorf("expected false, got true")
		}
	})

	// 8. Nil array
	t.Run("nil array", func(t *testing.T) {
		var nilArr Array[int]
		if nilArr.Includes(1) {
			t.Errorf("expected false, got true")
		}
	})

	// 9. String elements
	t.Run("string elements", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		if !a.Includes("b") {
			t.Errorf("expected true, got false")
		}
		if a.Includes("d") {
			t.Errorf("expected false, got true")
		}
	})

	// 10. Struct elements
	t.Run("struct elements", func(t *testing.T) {
		type user struct {
			name string
			age  int
		}
		a := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		if !a.Includes(user{"bob", 25}) {
			t.Errorf("expected true, got false")
		}
		if a.Includes(user{"zed", 99}) {
			t.Errorf("expected false, got true")
		}
	})

	// 11. Does not mutate receiver
	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Includes(2)
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_At(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	tests := []struct {
		name   string
		input  Array[int]
		index  int
		want   int
		wantOK bool
	}{
		{"index 0", Array[int]{1, 2, 3}, 0, 1, true},
		{"index 1", Array[int]{1, 2, 3}, 1, 2, true},
		{"index at end", Array[int]{1, 2, 3}, 2, 3, true},
		{"index equal to length (out of bounds)", Array[int]{1, 2, 3}, 3, 0, false},
		{"index past length (out of bounds)", Array[int]{1, 2, 3}, 10, 0, false},
		{"negative index -1", Array[int]{1, 2, 3}, -1, 3, true},
		{"negative index -2", Array[int]{1, 2, 3}, -2, 2, true},
		{"negative index -3", Array[int]{1, 2, 3}, -3, 1, true},
		{"negative index before start", Array[int]{1, 2, 3}, -4, 0, false},
		{"negative index past length", Array[int]{1, 2, 3}, -10, 0, false},
		{"zero value can be a real match", Array[int]{1, 0, 2}, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.input.At(tt.index)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("At(%d) = (%v, %v), want (%v, %v)", tt.index, got, ok, tt.want, tt.wantOK)
			}
		})
	}

	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		got, ok := empty.At(0)
		if got != 0 || ok {
			t.Errorf("At(0) = (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("nil array", func(t *testing.T) {
		var nilArr Array[int]
		got, ok := nilArr.At(0)
		if got != 0 || ok {
			t.Errorf("At(0) = (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("struct elements", func(t *testing.T) {
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got, ok := users.At(1)
		if !ok || got.name != "bob" || got.age != 25 {
			t.Errorf("At(1) = (%v, %v), want (bob, 25, true)", got, ok)
		}
	})

	t.Run("negative struct index", func(t *testing.T) {
		users := Array[user]{{"ann", 30}, {"bob", 25}, {"cy", 25}}
		got, ok := users.At(-1)
		if !ok || got.name != "cy" {
			t.Errorf("At(-1) = (%v, %v), want (cy, true)", got, ok)
		}
	})

	t.Run("string elements", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		got, ok := a.At(1)
		if !ok || got != "b" {
			t.Errorf("At(1) = (%v, %v), want (b, true)", got, ok)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_, _ = a.At(1)
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_Join(t *testing.T) {
	// 1. Basic join with default separator (MDN example)
	t.Run("default separator", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Join()
		if result != "1,2,3" {
			t.Errorf("expected 1,2,3, got %s", result)
		}
	})

	// 2. Custom separator
	t.Run("custom separator", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		result := a.Join("-")
		if result != "1-2-3" {
			t.Errorf("expected 1-2-3, got %s", result)
		}
	})

	// 3. Empty array
	t.Run("empty array", func(t *testing.T) {
		empty := Array[int]{}
		result := empty.Join()
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	// 4. Nil array
	t.Run("nil array", func(t *testing.T) {
		var nilArr Array[int]
		result := nilArr.Join()
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	// 5. String elements
	t.Run("string elements", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		result := a.Join()
		if result != "a,b,c" {
			t.Errorf("expected a,b,c, got %s", result)
		}
	})

	// 6. Struct elements
	t.Run("struct elements", func(t *testing.T) {
		type user struct {
			name string
			age  int
		}
		a := Array[user]{{"ann", 30}, {"bob", 25}}
		result := a.Join("|")
		want := "{ann 30}|{bob 25}"
		if result != want {
			t.Errorf("expected %s, got %s", want, result)
		}
	})

	// 7. Single element
	t.Run("single element", func(t *testing.T) {
		a := Array[int]{1}
		result := a.Join()
		if result != "1" {
			t.Errorf("expected 1, got %s", result)
		}
	})

	// 8. Does not mutate receiver
	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[int]{1, 2, 3}
		_ = a.Join()
		if len(a) != 3 || a[0] != 1 || a[1] != 2 || a[2] != 3 {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_Keys(t *testing.T) {
	collect := func(a Array[string]) []int {
		out := []int{}
		for i := range a.Keys() {
			out = append(out, i)
		}
		return out
	}

	tests := []struct {
		name  string
		input Array[string]
		want  []int
	}{
		{"yields indices in order", Array[string]{"a", "b", "c"}, []int{0, 1, 2}},
		{"single element", Array[string]{"a"}, []int{0}},
		{"empty array yields nothing", Array[string]{}, []int{}},
		{"nil array yields nothing", nil, []int{}},
		{"zero-value elements still have keys", Array[string]{"", "", ""}, []int{0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collect(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("break stops the iteration", func(t *testing.T) {
		a := Array[string]{"a", "b", "c", "d"}
		var seen []int
		for i := range a.Keys() {
			seen = append(seen, i)
			if i == 1 {
				break
			}
		}
		if want := []int{0, 1}; !reflect.DeepEqual(seen, want) {
			t.Errorf("visited %v, want %v", seen, want)
		}
	})

	t.Run("can be iterated more than once", func(t *testing.T) {
		a := Array[string]{"a", "b"}
		seq := a.Keys()
		count := 0
		for range seq {
			count++
		}
		for range seq {
			count++
		}
		if count != 4 {
			t.Errorf("counted %d keys over two passes, want 4", count)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		for range a.Keys() {
		}
		if want := (Array[string]{"a", "b", "c"}); !reflect.DeepEqual(a, want) {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_Values(t *testing.T) {
	type user struct {
		name string
		age  int
	}

	collect := func(a Array[string]) []string {
		out := []string{}
		for v := range a.Values() {
			out = append(out, v)
		}
		return out
	}

	tests := []struct {
		name  string
		input Array[string]
		want  []string
	}{
		{"yields values in order", Array[string]{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"single element", Array[string]{"a"}, []string{"a"}},
		{"empty array yields nothing", Array[string]{}, []string{}},
		{"nil array yields nothing", nil, []string{}},
		{"duplicate values are all yielded", Array[string]{"x", "x"}, []string{"x", "x"}},
		{"zero-value elements are yielded", Array[string]{"", "a", ""}, []string{"", "a", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collect(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("break stops the iteration", func(t *testing.T) {
		a := Array[string]{"a", "b", "c", "d"}
		var seen []string
		for v := range a.Values() {
			seen = append(seen, v)
			if v == "b" {
				break
			}
		}
		if want := []string{"a", "b"}; !reflect.DeepEqual(seen, want) {
			t.Errorf("visited %v, want %v", seen, want)
		}
	})

	t.Run("can be iterated more than once", func(t *testing.T) {
		a := Array[string]{"a", "b"}
		seq := a.Values()
		count := 0
		for range seq {
			count++
		}
		for range seq {
			count++
		}
		if count != 4 {
			t.Errorf("counted %d values over two passes, want 4", count)
		}
	})

	t.Run("is lazy and sees element changes made during iteration", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		var got []string
		first := true
		for v := range a.Values() {
			if first {
				a[2] = "changed"
				first = false
			}
			got = append(got, v)
		}
		if want := []string{"a", "b", "changed"}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("structs", func(t *testing.T) {
		users := Array[user]{{"ann", 30}, {"bob", 25}}
		var names []string
		for u := range users.Values() {
			names = append(names, u.name)
		}
		if want := []string{"ann", "bob"}; !reflect.DeepEqual(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("pointers are the same pointers", func(t *testing.T) {
		a, b := &user{"ann", 30}, &user{"bob", 25}
		ptrs := Array[*user]{a, b}
		var got []*user
		for p := range ptrs.Values() {
			got = append(got, p)
		}
		if len(got) != 2 || got[0] != a || got[1] != b {
			t.Errorf("Values() did not yield the original pointers: %v", got)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		for range a.Values() {
		}
		if want := (Array[string]{"a", "b", "c"}); !reflect.DeepEqual(a, want) {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}

func TestArray_Entries(t *testing.T) {
	collect := func(a Array[string]) []Entry[int, string] {
		out := []Entry[int, string]{}
		for e := range a.Entries() {
			out = append(out, e)
		}
		return out
	}

	tests := []struct {
		name  string
		input Array[string]
		want  []Entry[int, string]
	}{
		{"yields index and value in order", Array[string]{"a", "b", "c"}, []Entry[int, string]{{0, "a"}, {1, "b"}, {2, "c"}}},
		{"single element", Array[string]{"a"}, []Entry[int, string]{{0, "a"}}},
		{"empty array yields nothing", Array[string]{}, []Entry[int, string]{}},
		{"nil array yields nothing", nil, []Entry[int, string]{}},
		{"duplicate values keep their own indices", Array[string]{"x", "x"}, []Entry[int, string]{{0, "x"}, {1, "x"}}},
		{"zero-value elements are yielded", Array[string]{"", "a"}, []Entry[int, string]{{0, ""}, {1, "a"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collect(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Key is the index and Val is the value", func(t *testing.T) {
		a := Array[string]{"a", "b"}
		for e := range a.Entries() {
			if e.Val != a[e.Key] {
				t.Errorf("entry %+v does not match a[%d] = %q", e, e.Key, a[e.Key])
			}
		}
	})

	t.Run("break stops the iteration", func(t *testing.T) {
		a := Array[string]{"a", "b", "c", "d"}
		var seen []int
		for e := range a.Entries() {
			seen = append(seen, e.Key)
			if e.Key == 1 {
				break
			}
		}
		if want := []int{0, 1}; !reflect.DeepEqual(seen, want) {
			t.Errorf("visited %v, want %v", seen, want)
		}
	})

	t.Run("can be iterated more than once", func(t *testing.T) {
		a := Array[string]{"a", "b"}
		seq := a.Entries()
		count := 0
		for range seq {
			count++
		}
		for range seq {
			count++
		}
		if count != 4 {
			t.Errorf("counted %d entries over two passes, want 4", count)
		}
	})

	t.Run("is lazy and sees element changes made during iteration", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		var got []string
		for e := range a.Entries() {
			if e.Key == 0 {
				a[2] = "changed"
			}
			got = append(got, e.Val)
		}
		if want := []string{"a", "b", "changed"}; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("structs", func(t *testing.T) {
		type user struct {
			name string
			age  int
		}
		users := Array[user]{{"ann", 30}, {"bob", 25}}
		var got []Entry[int, user]
		for e := range users.Entries() {
			got = append(got, e)
		}
		want := []Entry[int, user]{{0, user{"ann", 30}}, {1, user{"bob", 25}}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("entries are copies, so editing one does not touch the array", func(t *testing.T) {
		a := Array[string]{"a", "b"}
		for e := range a.Entries() {
			e.Val = "changed"
		}
		if want := (Array[string]{"a", "b"}); !reflect.DeepEqual(a, want) {
			t.Errorf("array changed through an entry: %#v", a)
		}
	})

	t.Run("does not mutate receiver", func(t *testing.T) {
		a := Array[string]{"a", "b", "c"}
		for range a.Entries() {
		}
		if want := (Array[string]{"a", "b", "c"}); !reflect.DeepEqual(a, want) {
			t.Errorf("receiver mutated: %#v", a)
		}
	})
}
