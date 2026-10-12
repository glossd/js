package js

import (
	"cmp"
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
)

// Go prevents developers from taking the memory address (&) of certain temporary or dynamic values.
// To allow chain call functions the immutable methods have value receivers,
// as the array should not be copied, just the SliceHeader.

type Array[T any] []T

func (a Array[T]) Length() int {
	if a == nil {
		return 0
	}
	return len(a)
}

// Push adds one or more elements to the end of the array and returns the new length.
func (a *Array[T]) Push(values ...T) int {
	*a = append(*a, values...)
	return len(*a)
}

// Pop removes the last element from the array and returns it.
// If the array is empty, it returns the zero value of T and false.
func (a *Array[T]) Pop() (val T, ok bool) {
	if a == nil || len(*a) == 0 {
		return val, false
	}
	lastIdx := len(*a) - 1
	val = (*a)[lastIdx]
	var zero T
	(*a)[lastIdx] = zero // avoid memory leak
	*a = (*a)[:lastIdx]
	return val, true
}

// Shift removes the first element from the array and returns it.
// If the array is empty, it returns the zero value of T and false.
func (a *Array[T]) Shift() (val T, ok bool) {
	if a == nil || len(*a) == 0 {
		return val, false
	}
	val = (*a)[0]
	var zero T
	(*a)[0] = zero // avoid memory leak
	*a = (*a)[1:]
	return val, true
}

// Unshift adds one or more elements to the beginning of the array and returns the new length.
func (a *Array[T]) Unshift(values ...T) int {
	if len(values) == 0 {
		if a == nil {
			return 0
		}
		return len(*a)
	}
	newArr := make(Array[T], len(values)+len(*a))
	copy(newArr, values)
	copy(newArr[len(values):], *a)
	*a = newArr
	return len(*a)
}

// Splice changes the contents of an array by removing or replacing existing elements
// and/or adding new elements in place, and returns an Array containing the deleted elements.
func (a *Array[T]) Splice(start int, deleteCount int, items ...T) Array[T] {
	if a == nil {
		return Array[T]{}
	}
	l := len(*a)

	actualStart := start
	if actualStart < 0 {
		actualStart = l + actualStart
		if actualStart < 0 {
			actualStart = 0
		}
	} else if actualStart > l {
		actualStart = l
	}

	actualDeleteCount := deleteCount
	if actualDeleteCount < 0 {
		actualDeleteCount = 0
	} else if actualDeleteCount > l-actualStart {
		actualDeleteCount = l - actualStart
	}

	deleted := make(Array[T], actualDeleteCount)
	copy(deleted, (*a)[actualStart:actualStart+actualDeleteCount])

	itemCount := len(items)
	newLength := l - actualDeleteCount + itemCount

	newArr := make(Array[T], newLength)
	copy(newArr, (*a)[:actualStart])
	copy(newArr[actualStart:], items)
	copy(newArr[actualStart+itemCount:], (*a)[actualStart+actualDeleteCount:])

	// Avoid memory leak by zeroing removed slots in old slice
	var zero T
	for i := actualStart; i < actualStart+actualDeleteCount; i++ {
		(*a)[i] = zero
	}

	*a = newArr
	return deleted
}

// Sort sorts the elements of an array in place and returns the reference to the same array.
// If compareFn is not supplied or is nil, elements are sorted by converting them to strings
// and comparing them in UTF-8 / Unicode code point order (lexicographically).
// If compareFn is supplied, it determines the order of elements based on its return value:
// negative if a < b, positive if a > b, and zero if a == b.
// The sort is stable, preserving the relative order of elements that compare equal.
func (a *Array[T]) Sort(compareFn ...func(a, b T) int) *Array[T] {
	if a == nil || len(*a) <= 1 {
		return a
	}

	if len(compareFn) > 0 && compareFn[0] != nil {
		slices.SortStableFunc(*a, compareFn[0])
	} else {
		slices.SortStableFunc(*a, func(x, y T) int {
			return cmp.Compare(fmt.Sprint(x), fmt.Sprint(y))
		})
	}

	return a
}

// Reverse reverses an array in place and returns the reference to the same array.
func (a *Array[T]) Reverse() *Array[T] {
	if a == nil || len(*a) <= 1 {
		return a
	}
	slices.Reverse(*a)
	return a
}

// Fill changes all elements in an array to a static value, from a start index
// (default 0) to an end index (default array length), and returns the reference to the same array.
func (a *Array[T]) Fill(value T, bounds ...int) *Array[T] {
	if a == nil {
		return nil
	}
	l := len(*a)
	if l == 0 {
		return a
	}

	startIdx := 0
	endIdx := l

	if len(bounds) > 0 {
		startIdx = bounds[0]
	}
	if len(bounds) > 1 {
		endIdx = bounds[1]
	}

	if startIdx < 0 {
		startIdx = l + startIdx
		if startIdx < 0 {
			startIdx = 0
		}
	} else if startIdx > l {
		startIdx = l
	}

	if endIdx < 0 {
		endIdx = l + endIdx
		if endIdx < 0 {
			endIdx = 0
		}
	} else if endIdx > l {
		endIdx = l
	}

	for i := startIdx; i < endIdx; i++ {
		(*a)[i] = value
	}

	return a
}

// CopyWithin shallow copies part of an array to another location in the same array
// without modifying its length, and returns the reference to the same array.
// target is the zero-based index at which to copy the sequence to.
// Optional bounds can specify start (default 0) and end (default array length).
func (a *Array[T]) CopyWithin(target int, bounds ...int) *Array[T] {
	if a == nil {
		return nil
	}
	l := len(*a)
	if l == 0 {
		return a
	}

	to := target
	if to < 0 {
		to = l + to
		if to < 0 {
			to = 0
		}
	} else if to > l {
		to = l
	}

	from := 0
	if len(bounds) > 0 {
		from = bounds[0]
	}
	if from < 0 {
		from = l + from
		if from < 0 {
			from = 0
		}
	} else if from > l {
		from = l
	}

	final := l
	if len(bounds) > 1 {
		final = bounds[1]
	}
	if final < 0 {
		final = l + final
		if final < 0 {
			final = 0
		}
	} else if final > l {
		final = l
	}

	count := min(final-from, l-to)
	if count > 0 {
		copy((*a)[to:to+count], (*a)[from:from+count])
	}

	return a
}

// Here begin the Immutable methods.

// Filter creates a shallow copy of a portion of the array, filtered down to just
// the elements that pass the test implemented by the provided function.
// Filter does not mutate the original array.
func (a Array[T]) Filter(callbackFn func(element T) bool) Array[T] {
	if a == nil {
		return Array[T]{}
	}

	src := a
	result := make(Array[T], 0, len(src))
	for _, v := range src {
		if callbackFn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map creates a new array populated with the results of calling the provided
// function on every element in the calling array.
// Map does not mutate the original array.
func (a Array[T]) Map[R any](callbackFn func(element T) R) Array[R] {
	if a == nil {
		return Array[R]{}
	}

	src := a
	result := make(Array[R], len(src))
	for i, v := range src {
		result[i] = callbackFn(v)
	}
	return result
}

// FlatMap returns a new array formed by applying the provided function to each
// element and then concatenating the results by one level.
// FlatMap does not mutate the original array.
func (a Array[T]) FlatMap[R any](callbackFn func(element T) Array[R]) Array[R] {
	if a == nil {
		return Array[R]{}
	}

	src := a
	result := make(Array[R], 0)
	for _, v := range src {
		result = append(result, callbackFn(v)...)
	}
	return result
}

// Reduce executes a user-supplied "reducer" callback function on each element of the array,
// in order, passing in the return value from the calculation on the preceding element.
// The final result of running the reducer across all elements is a single value.
// If initialValue is provided, accumulator is initialized to initialValue, and callbackFn
// runs over all elements.
// If initialValue is not provided, accumulator is initialized to the first element (which must
// be of type R), and callbackFn runs over the remaining elements.
// If the array is empty and no initialValue is provided, Reduce panics with a TypeError.
// Reduce does not mutate the original array.
func (a Array[T]) Reduce[R any](callbackFn func(accumulator R, element T) R, initialValue ...R) R {
	return a.reduce(callbackFn, false, initialValue...)
}

// ReduceRight applies a function against an accumulator and each value of the array (from right-to-left)
// to reduce it to a single value.
// If initialValue is provided, accumulator is initialized to initialValue, and callbackFn
// runs over all elements from right to left.
// If initialValue is not provided, accumulator is initialized to the last element (which must
// be of type R), and callbackFn runs over the remaining elements from right to left.
// If the array is empty and no initialValue is provided, ReduceRight panics with a TypeError.
// ReduceRight does not mutate the original array.
func (a Array[T]) ReduceRight[R any](callbackFn func(accumulator R, element T) R, initialValue ...R) R {
	return a.reduce(callbackFn, true, initialValue...)
}

func (a Array[T]) reduce[R any](callbackFn func(accumulator R, element T) R, fromRight bool, initialValue ...R) R {
	if a == nil || len(a) == 0 {
		if len(initialValue) > 0 {
			return initialValue[0]
		}
		panic("TypeError: Reduce of empty array with no initial value")
	}

	src := a
	var accumulator R
	var startIndex int

	if len(initialValue) > 0 {
		accumulator = initialValue[0]
		if fromRight {
			startIndex = len(src) - 1
		} else {
			startIndex = 0
		}
	} else {
		initIdx := 0
		if fromRight {
			initIdx = len(src) - 1
		}
		if acc, ok := any(src[initIdx]).(R); ok {
			accumulator = acc
			if fromRight {
				startIndex = len(src) - 2
			} else {
				startIndex = 1
			}
		} else {
			panic("TypeError: Reduce of array with no initial value and incompatible accumulator type")
		}
	}

	if fromRight {
		for i := startIndex; i >= 0; i-- {
			accumulator = callbackFn(accumulator, src[i])
		}
	} else {
		for i := startIndex; i < len(src); i++ {
			accumulator = callbackFn(accumulator, src[i])
		}
	}

	return accumulator
}

// ForEach executes a provided function once for each array element.
func (a Array[T]) ForEach(callbackFn func(element T)) {
	if a == nil || len(a) == 0 || callbackFn == nil {
		return
	}

	for _, v := range a {
		callbackFn(v)
	}
}

// Every tests whether all elements in the array pass the test implemented by the provided function.
// It returns true for an empty array.
// Every does not mutate the original array.
func (a Array[T]) Every(callbackFn func(element T) bool) bool {
	if a == nil || len(a) == 0 {
		return true
	}

	for _, v := range a {
		if !callbackFn(v) {
			return false
		}
	}
	return true
}

// Some tests whether at least one element in the array passes the test implemented by the provided function.
// It returns false for an empty array.
// Some does not mutate the original array.
func (a Array[T]) Some(callbackFn func(element T) bool) bool {
	if a == nil || len(a) == 0 {
		return false
	}

	for _, v := range a {
		if callbackFn(v) {
			return true
		}
	}
	return false
}

// Concat is used to merge two or more arrays. This method does not change
// the existing arrays, but instead returns a new array.
func (a Array[T]) Concat(items ...Array[T]) Array[T] {
	totalLen := len(a)
	for _, item := range items {
		totalLen += len(item)
	}

	result := make(Array[T], 0, totalLen)
	result = append(result, a...)
	for _, item := range items {
		result = append(result, item...)
	}
	return result
}

// Slice returns a shallow copy of a portion of an array into a new array object
// selected from start to end (end not included) where start and end represent
// the index of items in that array.
// Optional bounds can specify start (default 0) and end (default array length).
// The original array is not modified.
func (a Array[T]) Slice(bounds ...int) Array[T] {
	l := len(a)
	if a == nil || l == 0 {
		return Array[T]{}
	}

	startIdx := 0
	endIdx := l

	if len(bounds) > 0 {
		startIdx = bounds[0]
	}
	if len(bounds) > 1 {
		endIdx = bounds[1]
	}

	if startIdx < 0 {
		startIdx = l + startIdx
		if startIdx < 0 {
			startIdx = 0
		}
	} else if startIdx > l {
		startIdx = l
	}

	if endIdx < 0 {
		endIdx = l + endIdx
		if endIdx < 0 {
			endIdx = 0
		}
	} else if endIdx > l {
		endIdx = l
	}

	if startIdx >= endIdx {
		return Array[T]{}
	}

	result := make(Array[T], endIdx-startIdx)
	copy(result, a[startIdx:endIdx])
	return result
}

// Flat flattens the array. The element type of a
// flattened result can't be expressed in terms of T for arbitrary nesting
// depth, so the result is an Array[any].
//
// depth is optional and defaults to 1, like in JavaScript. Pass 0 to copy
// without flattening, or -1 to flatten completely (acts as Infinity). Any element whose
// kind is slice or array is treated as a nested array (strings are not).
// The receiver is never modified.
func (a Array[T]) Flat(depth ...int) Array[any] {
	d := 1
	if len(depth) > 0 {
		d = depth[0]
	}
	out := make(Array[any], 0, len(a))
	for _, item := range a {
		out = flatten(out, item, d)
	}
	return out
}

func flatten(out Array[any], item any, depth int) Array[any] {
	if depth != 0 {
		rv := reflect.ValueOf(item)
		if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
			next := depth
			if depth > 0 {
				next = depth - 1
			}
			for i := 0; i < rv.Len(); i++ {
				out = flatten(out, rv.Index(i).Interface(), next)
			}
			return out
		}
	}
	return append(out, item)
}

// ToSorted returns a sorted copy and leaves the receiver untouched. The sort is
// stable, like in JavaScript.
//
// Go has no built-in ordering for an arbitrary T, so the comparator is
// passed in and follows the usual Go convention: negative if a < b, zero if
// a == b, positive if a > b.
//
// If compare is nil, elements are ordered by their string form (fmt.Sprint),
// which mirrors JavaScript's default of sorting by string conversion. That
// means 10 sorts before 9, exactly as `[10, 9].toSorted()` does in JS.
func (a Array[T]) ToSorted(compare func(a, b T) int) Array[T] {
	out := make(Array[T], len(a))
	copy(out, a)
	if compare == nil {
		compare = func(x, y T) int {
			return cmp.Compare(fmt.Sprint(x), fmt.Sprint(y))
		}
	}
	slices.SortStableFunc(out, compare)
	return out
}

// ToReversed is a port of JavaScript's Array.prototype.toReversed.
//
// It returns a reversed copy and leaves the receiver untouched.
func (a Array[T]) ToReversed() Array[T] {
	out := make(Array[T], len(a))
	copy(out, a)
	slices.Reverse(out)
	return out
}

// ToSpliced is a port of JavaScript's Array.prototype.toSpliced.
//
// It returns a copy of the array with deleteCount elements removed starting
// at start, and items inserted in their place. The receiver is never
// modified.
//
// Index handling follows JavaScript:
//   - a negative start counts from the end, and is clamped to 0
//   - a start past the end is clamped to len(a)
//   - a negative deleteCount is treated as 0
//   - a deleteCount larger than the remaining length is clamped, so pass
//     math.MaxInt to delete everything from start onwards (JS does this when
//     deleteCount is omitted, which Go can't express)
func (a Array[T]) ToSpliced(start, deleteCount int, items ...T) Array[T] {
	n := len(a)

	if start < 0 {
		start = max(n+start, 0)
	} else if start > n {
		start = n
	}
	deleteCount = min(max(deleteCount, 0), n-start)

	out := make(Array[T], 0, n-deleteCount+len(items))
	out = append(out, a[:start]...)
	out = append(out, items...)
	out = append(out, a[start+deleteCount:]...)
	return out
}

// With is a port of JavaScript's Array.prototype.with.
//
// It returns a copy of the array with the element at index replaced by
// value. A negative index counts from the end. The receiver is never
// modified.
//
// Like JavaScript, which throws a RangeError, With panics if the index is out
// of range (index >= len(a) or index < -len(a)).
func (a Array[T]) With(index int, value T) Array[T] {
	n := len(a)
	i := index
	if i < 0 {
		i += n
	}
	if i < 0 || i >= n {
		panic(fmt.Sprintf("array: With index %d out of range for length %d", index, n))
	}

	out := make(Array[T], n)
	copy(out, a)
	out[i] = value
	return out
}

// Find is a port of JavaScript's Array.prototype.find.
//
// It returns the first element for which predicate reports true, and true.
// If no element matches, it returns the zero value of T and false. JavaScript
// signals "not found" with undefined, which Go can't do for an arbitrary T,
// and the zero value alone would be ambiguous (0 may be a real match), so
// the second result tells the two cases apart.
//
// Elements are visited in order and the search stops at the first match.
// The receiver is never modified.
func (a Array[T]) Find(predicate func(value T) bool) (T, bool) {
	for _, v := range a {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FindLast is a port of JavaScript's Array.prototype.findLast.
//
// It returns the last element for which predicate reports true, and true.
// If no element matches, it returns the zero value of T and false (see Find
// for why a second result is needed).
//
// Elements are visited from the end towards the start and the search stops
// at the first match found that way. The receiver is never modified.
func (a Array[T]) FindLast(predicate func(value T) bool) (T, bool) {
	for i := len(a) - 1; i >= 0; i-- {
		if predicate(a[i]) {
			return a[i], true
		}
	}
	var zero T
	return zero, false
}

// Includes is a port of JavaScript's Array.prototype.includes.
//
// It determines whether an array includes a certain value among its entries,
// returning true or false as appropriate.
//
// The method accepts a second argument, fromIndex, which is the index to start
// the search at. Defaults to 0. A negative index counts from the end of the
// array. If the index is greater than or equal to the array length, false is
// returned. If the index is less than 0, the entire array is searched.
//
// The receiver is never modified.
func (a Array[T]) Includes(searchElement T, fromIndex ...int) bool {
	if len(a) == 0 {
		return false
	}

	idx := 0
	if len(fromIndex) > 0 {
		idx = fromIndex[0]
		if idx < 0 {
			idx += len(a)
			if idx < 0 {
				idx = 0
			}
		} else if idx >= len(a) {
			return false
		}
	}

	for i := idx; i < len(a); i++ {
		if reflect.DeepEqual(a[i], searchElement) {
			return true
		}
	}
	return false
}

// At is a port of JavaScript's Array.prototype.at.
//
// It returns the element at the given index. A negative index counts back from
// the end of the array. If the index is out of bounds (index >= len(a) or
// index < -len(a)), it returns the zero value of T and false.
//
// The receiver is never modified.
func (a Array[T]) At(index int) (T, bool) {
	n := len(a)
	if n == 0 {
		var zero T
		return zero, false
	}

	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		var zero T
		return zero, false
	}

	return a[index], true
}

// FindIndex is a port of JavaScript's Array.prototype.findIndex.
//
// It returns the index of the first element for which predicate reports true,
// or -1 if no element matches. This mirrors JavaScript's behavior where
// findIndex returns -1 (not undefined) when nothing is found.
//
// Elements are visited in order and the search stops at the first match.
// The receiver is never modified.
func (a Array[T]) FindIndex(predicate func(value T) bool) int {
	for i, v := range a {
		if predicate(v) {
			return i
		}
	}
	return -1
}

// FindLastIndex is a port of JavaScript's Array.prototype.findLastIndex.
//
// It returns the index of the last element for which predicate reports true,
// or -1 if no element matches.
//
// Elements are visited from the end towards the start and the search stops
// at the first match found that way. The receiver is never modified.
func (a Array[T]) FindLastIndex(predicate func(value T) bool) int {
	for i := len(a) - 1; i >= 0; i-- {
		if predicate(a[i]) {
			return i
		}
	}
	return -1
}

// Join is a port of JavaScript's Array.prototype.join.
//
// It joins all elements of an array into a string, separated by the specified
// separator string. If separator is not provided, elements are separated with
// a comma (","). Elements are converted to strings using fmt.Sprint.
//
// If the array is empty, the result is an empty string.
//
// The receiver is never modified.
func (a Array[T]) Join(separator ...string) string {
	sep := ","
	if len(separator) > 0 {
		sep = separator[0]
	}
	if len(a) == 0 {
		return ""
	}
	var b strings.Builder
	for i, v := range a {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(fmt.Sprint(v))
	}
	return b.String()
}

// Keys is a port of JavaScript's Array.prototype.keys.
//
// It returns an iterator over the indices of the array, in order, as an
// iter.Seq[int] that can be ranged over directly:
//
//	for i := range a.Keys() {
//		fmt.Println(i)
//	}
//
// Like in JavaScript, every index is yielded, including those of zero-value
// elements. Breaking out of the loop stops the iteration. The length is
// checked on each step, as in Entries. The receiver is never modified.
func (a Array[T]) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < len(a); i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Values is a port of JavaScript's Array.prototype.values.
//
// It returns an iterator over the elements of the array, in order, as an
// iter.Seq[T] that can be ranged over directly:
//
//	for v := range a.Values() {
//		fmt.Println(v)
//	}
//
// Breaking out of the loop stops the iteration. Like in JavaScript the
// iterator is lazy and reads each element when it is reached, so changes
// made to existing elements during the loop are seen by later iterations.
// The receiver is never modified.
func (a Array[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(a); i++ {
			if !yield(a[i]) {
				return
			}
		}
	}
}

// Entries is a port of JavaScript's Array.prototype.entries.
//
// It returns an iterator over the index/value pairs of the array, in order,
// as an iter.Seq[Entry[int, T]] that can be ranged over directly:
//
//	for e := range a.Entries() {
//		fmt.Println(e.Key, e.Val)
//	}
//
// Breaking out of the loop stops the iteration. Like in JavaScript the
// iterator is lazy and reads each element when it is reached, so changes made
// to existing elements during the loop are seen by later iterations.
// The receiver is never modified by Entries itself.
func (a Array[T]) Entries() iter.Seq[Entry[int, T]] {
	return func(yield func(Entry[int, T]) bool) {
		for i := 0; i < len(a); i++ {
			if !yield(Entry[int, T]{Key: i, Val: a[i]}) {
				return
			}
		}
	}
}

func (a Array[T]) ToString() string {
	return a.Join()
}

// ArrayOf is a port of JavaScript's Array.of.
//
// It creates an Array from its arguments, whatever their number:
//
//	ArrayOf(1, 2, 3) // Array[int]{1, 2, 3}
//	ArrayOf[string]() // empty, non-nil Array[string]
//
// Go has no static methods on a generic type, so this is a package-level
// function, and the element type is usually inferred from the arguments.
//
// The result always has its own backing array. Calling Of(s...) with an
// existing slice therefore copies it rather than aliasing it, which matches
// JavaScript, where Array.of(...arr) builds a new array.
func ArrayOf[T any](values ...T) Array[T] {
	out := make(Array[T], len(values))
	copy(out, values)
	return out
}

// ArrayFrom is a port of JavaScript's Array.from.
func ArrayFrom[T any](seq iter.Seq[T]) Array[T] {
	arr := ArrayOf[T]()
	for v := range seq {
		arr = append(arr, v)
	}
	return arr
}
