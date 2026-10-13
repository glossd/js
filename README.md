# JS ES6 Data Structures in Go

Go 1.27 introduced generic methods, allowing existence of the most popular methods such as **Map**, **FlatMap**, and **Reduce**.

#### Here is a simple example of processing an array of text inputs
```go
convertToNumber := func(s string) *float64 {
    value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
    if err != nil {
        return nil
    }
    return &value
}
total := js.ArrayOf("10", "abc", "25.5", "  42  ").
    Map(convertToNumber).
    Filter(func(n *float64) bool { return n != nil }).
    Reduce(func(acc float64, v *float64) float64 { return acc + *v }, 0)

// total == 77.5
```

#### If you want, you can still do normal loops
```go
for i, v := range js.ArrayOf(1, 2, 3) {
	...
}
```

#### You can cast Go slices and maps to access the methods.
```go
yourSlice := []string{"hello", "world"}
joined := js.Array[string](yourSlice).Join("-")
// joined == "hello-world" 
```

#### Map can be converted to Entries, Values, Keys the same way as in JS
```go 
myMap := NewMap[string, int]()
myMap.Set("key1", 12)
myMap.Set("key2", 123)

sum := ArrayFrom(myMap.Values()).Reduce(func(acc int, e int) int { return acc + e })
// sum == 135

// the same map could be initiated with the helper functions.
myMap := NewMap(E2("key1", 12, "key2", 123)...)
// To keep the Map strongly typed, type Entry{Key, Val} had to be introduced to represent tuples in JS
```

#### Set can be converted to Values, the same way as in JS
```go 
fruits := NewSet("apple", "banana", "avocado")
aFruits := ArrayFrom(fruits.Values()).Filter(func(v string) bool { return strings.HasPrefix(v, "a") })

// aFruies == ["apple", "avocado"]
```

---

## Array

| Method | Signature | Description |
|--------|-----------|-------------|
| Length | `Length() int` | Returns the number of elements |
| Push | `Push(values ...T) int` | Adds elements to the end, returns new length |
| Pop | `Pop() (val T, ok bool)` | Removes and returns the last element |
| Shift | `Shift() (val T, ok bool)` | Removes and returns the first element |
| Unshift | `Unshift(values ...T) int` | Adds elements to the beginning, returns new length |
| Splice | `Splice(start, deleteCount int, items ...T) Array[T]` | Removes/inserts elements in place |
| Sort | `Sort(compareFn ...func(a, b T) int) *Array[T]` | Sorts in place, returns the array |
| Reverse | `Reverse() *Array[T]` | Reverses in place, returns the array |
| Fill | `Fill(value T, bounds ...int) *Array[T]` | Fills a range with a value, returns the array |
| CopyWithin | `CopyWithin(target int, bounds ...int) *Array[T]` | Copies within the array, returns the array |
| Filter | `Filter(callbackFn func(element T) bool) Array[T]` | Returns a new array of filtered elements |
| Map | `Map[R any](callbackFn func(element T) R) Array[R]` | Returns a new array of transformed elements |
| FlatMap | `FlatMap[R any](callbackFn func(element T) Array[R]) Array[R]` | Returns a new array formed by applying and flattening |
| Reduce | `Reduce[R any](callbackFn func(accumulator R, element T) R, initialValue ...R) R` | Reduces the array to a single value |
| ReduceRight | `ReduceRight[R any](callbackFn func(accumulator R, element T) R, initialValue ...R) R` | Reduces from right to left |
| ForEach | `ForEach(callbackFn func(element T))` | Executes a function for each element |
| Every | `Every(callbackFn func(element T) bool) bool` | Tests whether all elements pass a test |
| Some | `Some(callbackFn func(element T) bool) bool` | Tests whether at least one element passes |
| Concat | `Concat(items ...Array[T]) Array[T]` | Merges arrays into a new array |
| Slice | `Slice(bounds ...int) Array[T]` | Returns a shallow copy of a portion |
| Flat | `Flat(depth ...int) Array[any]` | Flattens nested arrays |
| ToSorted | `ToSorted(compare func(a, b T) int) Array[T]` | Returns a sorted copy |
| ToReversed | `ToReversed() Array[T]` | Returns a reversed copy |
| ToSpliced | `ToSpliced(start, deleteCount int, items ...T) Array[T]` | Returns a copy with removed/inserted elements |
| With | `With(index int, value T) Array[T]` | Returns a copy with element at index replaced |
| Find | `Find(predicate func(value T) bool) (T, bool)` | Returns the first element matching the predicate |
| FindLast | `FindLast(predicate func(value T) bool) (T, bool)` | Returns the last element matching the predicate |
| Includes | `Includes(searchElement T, fromIndex ...int) bool` | Tests whether an array contains a value |
| At | `At(index int) (T, bool)` | Returns the element at the given index |
| FindIndex | `FindIndex(predicate func(value T) bool) int` | Returns the index of the first match, or -1 |
| FindLastIndex | `FindLastIndex(predicate func(value T) bool) int` | Returns the index of the last match, or -1 |
| IndexOf | `IndexOf(searchElement T, fromIndex ...int) int` | Returns the first index of an element, or -1 |
| LastIndexOf | `LastIndexOf(searchElement T, fromIndex ...int) int` | Returns the last index of an element, or -1 |
| Join | `Join(separator ...string) string` | Joins all elements into a string |
| Keys | `Keys() iter.Seq[int]` | Returns an iterator over the indices |
| Values | `Values() iter.Seq[T]` | Returns an iterator over the elements |
| Entries | `Entries() iter.Seq[Entry[int, T]]` | Returns an iterator over index/value pairs |
| ToString | `ToString() string` | Joins all elements into a string |

### Package-level functions

| Function | Description |
|----------|-------------|
| `ArrayOf[T any](values ...T) Array[T]` | Creates an Array from its arguments |
| `ArrayFrom[T any](seq iter.Seq[T]) Array[T]` | Creates an Array from a sequence |

---

## Map

| Method | Signature | Description |
|--------|-----------|-------------|
| Set | `Set(key K, value V) *Map[K, V]` | Adds or updates a key-value pair, returns the Map for chaining |
| Get | `Get(key K) (V, bool)` | Returns the value and whether the key exists |
| Has | `Has(key K) bool` | Returns whether the key exists |
| Delete | `Delete(key K) bool` | Removes the key, returns whether it was deleted |
| Clear | `Clear() *Map[K, V]` | Removes all elements, returns the Map for chaining |
| Keys | `Keys() iter.Seq[K]` | Returns an iterator over the keys |
| Values | `Values() iter.Seq[V]` | Returns an iterator over the values |
| Entries | `Entries() iter.Seq[Entry[K, V]]` | Returns an iterator over key-value entries |
| Size | `Size() int` | Returns the number of elements |
| GetOrInsert | `GetOrInsert(key K, value V) V` | Returns the value for the key if it exists, otherwise inserts and returns the given value |
| GetOrInsertComputed | `GetOrInsertComputed(key K, fn func() V) V` | Returns the value for the key if it exists, otherwise computes, inserts, and returns the value |
| ForEach | `ForEach(fn func(value V, key K))` | Calls fn once for each key/value pair |
| ForEachUntil | `ForEachUntil(fn func(value V, key K) bool)` | Calls fn once for each key/value pair, stopping when fn returns false |

### Package-level functions

| Function | Description |
|----------|-------------|
| `NewMap[K comparable, V any](entries ...Entry[K, V]) Map[K, V]` | Creates a Map from entries |

### Entry helpers

| Function | Description |
|----------|-------------|
| `E[K comparable, V any](key K, value V) Entry[K, V]` | Creates a single entry |
| `E2`, `E3`, `E4`, `E5`, `E6`, `E7`, `E8`, `E9` | Create multiple entries at once |

### Set

| Method | Signature | Description |
|--------|-----------|-------------|
| `NewSet[T comparable](values ...T) Set[T]` | Creates a Set from values |
| `Size()` | Returns the number of elements |
| `Add(v T) Set[T]` | Inserts v into the set (no-op if already present) |
| `Delete(v T) bool` | Removes v from the set, returns whether it was present |
| `Has(v T) bool` | Reports whether v is in the set |
| `Clear()` | Removes all values from the set |
| `IsSubsetOf(other Set[T]) bool` | Reports whether every value of s is also in other |
| `IsSupersetOf(other Set[T]) bool` | Reports whether every value of other is also in s |
| `IsDisjointFrom(other Set[T]) bool` | Reports whether s and other have no values in common |
| `Union(other Set[T]) Set[T]` | Returns a new set containing every value that is in s or in other |
| `Intersection(other Set[T]) Set[T]` | Returns a new set containing only the values that are in both s and other |
| `Difference(other Set[T]) Set[T]` | Returns a new set containing the values that are in s but not in other |
| `SymmetricDifference(other Set[T]) Set[T]` | Returns a new set containing the values that are in exactly one of s and other |

### Package-level functions

| Function                                        | Description               |
|-------------------------------------------------|---------------------------|
| `NewSet[T comparable](values ...T) Set[T]` | Creates a Set from values |

---

<sub><sup>Caution: the library was vibe-coded, use it at your own risk.</sub></sub>