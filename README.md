# JS Data Structures in Go

[![Go Version](https://img.shields.io/badge/Go-1.27%2B-blue)](https://golang.org/) [![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go library that mirrors JavaScript's Array, Map, and Set data structures with their native methods.

## Why?

Go 1.27 introduced generic methods, enabling a faithful port of JavaScript's most popular array/map methods like **Map**, **FlatMap**, and **Reduce**.

```go
// Convert string inputs to numbers and sum them
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

## Installation

```shell
go get github.com/glossd/js
```

Requires Go 1.27 or higher.

---

## Quick Start

### Cast and use with native Go slices/maps

```go
// Convert a slice to Array
yourSlice := []string{"hello", "world"}
joined := js.Array[string](yourSlice).Join("-")
// joined == "hello-world"
```

### For...of style iteration

```go
// Normal Go range loop is still supported
for i, v := range js.ArrayOf(1, 2, 3) {
    fmt.Printf("index %d: %d\n", i, v)
}
```

### Create a Map

```go
myMap := js.NewMap[string, int]()
myMap.Set("key1", 12).Set("key2", 123)

// Get a value (returns zero value and false if not found)
if val, ok := myMap.Get("key1"); ok {
    fmt.Println("key1:", val)
}

// Check existence
if myMap.Has("key1") {
    fmt.Println("key1 exists!")
}
```

### Create a Set

```go
fruits := js.NewSet("apple", "banana", "avocado")

if fruits.Has("apple") {
    fmt.Println("apple is in the set!")
}
```

### Convert iterators to Arrays/Sets

```go
// Map keys/values/entries to Array
sum := js.ArrayFrom(myMap.Values()).Reduce(func(acc, v int) int {
    return acc + v
}, 0)

// Array from Map entries
for entry := range myMap.Entries() {
    fmt.Println(entry.Key, entry.Val)
}

// Set from values
aFruits := js.ArrayFrom(fruits.Values()).
    Filter(func(v string) bool { return strings.HasPrefix(v, "a") })
// ["apple", "avocado"]
```

---

## API Reference

### Array

| Method | Signature | Description |
|--------|-----------|-------------|
| **Length** | `Length() int` | Returns the number of elements |
| **Push** | `Push(values ...T) int` | Adds elements to the end |
| **Pop** | `Pop() (T, bool)` | Removes and returns last element |
| **Shift** | `Shift() (T, bool)` | Removes and returns first element |
| **Unshift** | `Unshift(values ...T) int` | Adds elements to the beginning |
| **Splice** | `Splice(start, deleteCount int, items ...T) Array[T]` | Removes/inserts elements in place |
| **Sort** | `Sort(compareFn ...func(a, b T) int) *Array[T]` | Sorts in place |
| **Reverse** | `Reverse() *Array[T]` | Reverses in place |
| **Fill** | `Fill(value T, bounds ...int) *Array[T]` | Fills range with value |
| **CopyWithin** | `CopyWithin(target int, bounds ...int) *Array[T]` | Copies within array |
| **Filter** | `Filter(predicate func(T) bool) Array[T]` | Returns filtered copy |
| **Map** | `Map[R any](fn func(T) R) Array[R]` | Returns transformed copy |
| **FlatMap** | `FlatMap[R any](fn func(T) Array[R]) Array[R]` | Maps and flattens |
| **Reduce** | `Reduce[R any](fn func(R, T) R, init ...R) R` | Reduces to single value |
| **ReduceRight** | `ReduceRight[R any](fn func(R, T) R, init ...R) R` | Reduces from right |
| **ForEach** | `ForEach(fn func(T))` | Executes for each element |
| **Every** | `Every(fn func(T) bool) bool` | Tests if all pass |
| **Some** | `Some(fn func(T) bool) bool` | Tests if any pass |
| **Concat** | `Concat(items ...Array[T]) Array[T]` | Merges arrays |
| **Slice** | `Slice(bounds ...int) Array[T]` | Returns portion |
| **Flat** | `Flat(depth ...int) Array[any]` | Flattens nested arrays |
| **At** | `At(index int) (T, bool)` | Gets element at index (supports negative) |
| **Find** | `Find(predicate func(T) bool) (T, bool)` | Finds first matching element |
| **FindLast** | `FindLast(predicate func(T) bool) (T, bool)` | Finds last matching element |
| **FindIndex** | `FindIndex(predicate func(T) bool) int` | Index of first match, -1 if none |
| **FindLastIndex** | `FindLastIndex(predicate func(T) bool) int` | Index of last match, -1 if none |
| **Includes** | `Includes(search T, fromIndex ...int) bool` | Tests if value exists |
| **IndexOf** | `IndexOf(search T, fromIndex ...int) int` | First index of value, -1 if none |
| **LastIndexOf** | `LastIndexOf(search T, fromIndex ...int) int` | Last index of value, -1 if none |
| **Join** | `Join(sep ...string) string` | Joins elements with separator |
| **ToSorted** | `ToSorted(compare func(T, T) int) Array[T]` | Returns sorted copy |
| **ToReversed** | `ToReversed() Array[T]` | Returns reversed copy |
| **ToSpliced** | `ToSpliced(start, deleteCount int, items ...T) Array[T]` | Returns spliced copy |
| **With** | `With(index int, value T) Array[T]` | Returns copy with element replaced |
| **Keys** | `Keys() iter.Seq[int]` | Iterator over indices |
| **Values** | `Values() iter.Seq[T]` | Iterator over elements |
| **Entries** | `Entries() iter.Seq[Entry[int, T]]` | Iterator over [index, value] pairs |
| **ToString** | `ToString() string` | Joins with comma separator |

#### Package-level functions

| Function | Description |
|----------|-------------|
| `ArrayOf[T](values ...T) Array[T]` | Creates array from arguments |
| `ArrayFrom[T](seq iter.Seq[T]) Array[T]` | Creates array from iterator |

---

### Map

| Method | Signature | Description |
|--------|-----------|-------------|
| **Set** | `Set(key K, value V) *Map[K, V]` | Adds/updates entry, returns map |
| **Get** | `Get(key K) (V, bool)` | Gets value, returns found flag |
| **Has** | `Has(key K) bool` | Checks key existence |
| **Delete** | `Delete(key K) bool` | Removes entry |
| **Clear** | `Clear() *Map[K, V]` | Removes all entries |
| **GetOrInsert** | `GetOrInsert(key K, value V) V` | Get or insert with given value |
| **GetOrInsertComputed** | `GetOrInsertComputed(key K, fn func() V) V` | Get or compute and insert |
| **Keys** | `Keys() iter.Seq[K]` | Iterator over keys |
| **Values** | `Values() iter.Seq[V]` | Iterator over values |
| **Entries** | `Entries() iter.Seq[Entry[K, V]]` | Iterator over entries |
| **Size** | `Size() int` | Returns element count |
| **ForEach** | `ForEach(fn func(V, K))` | Iterates entries |
| **ForEachUntil** | `ForEachUntil(fn func(V, K) bool)` | Iterates until fn returns false |

#### Package-level functions

| Function | Description |
|----------|-------------|
| `NewMap[K, V](entries ...Entry[K, V]) Map[K, V]` | Creates map from entries |

#### Entry helpers

| Function | Description |
|----------|-------------|
| `E(key, value)` | Creates single entry |
| `E2` \| `E3` \| ... \| `E9` | Creates multiple entries |

---

### Set

| Method | Signature | Description |
|--------|-----------|-------------|
| **Add** | `Add(v T) Set[T]` | Inserts value |
| **Delete** | `Delete(v T) bool` | Removes value |
| **Has** | `Has(v T) bool` | Checks value existence |
| **Clear** | `Clear()` | Removes all values |
| **Size** | `Size() int` | Returns element count |
| **IsSubsetOf** | `IsSubsetOf(other Set[T]) bool` | Checks subset relation |
| **IsSupersetOf** | `IsSupersetOf(other Set[T]) bool` | Checks superset relation |
| **IsDisjointFrom** | `IsDisjointFrom(other Set[T]) bool` | Checks if disjoint |
| **Union** | `Union(other Set[T]) Set[T]` | Returns union |
| **Intersection** | `Intersection(other Set[T]) Set[T]` | Returns intersection |
| **Difference** | `Difference(other Set[T]) Set[T]` | Returns difference |
| **SymmetricDifference** | `SymmetricDifference(other Set[T]) Set[T]` | Returns symmetric diff |

#### Package-level functions

| Function | Description |
|----------|-------------|
| `NewSet[T](values ...T) Set[T]` | Creates set from values |

---

## Comparison with Standard Go

| JavaScript | This Library | Standard Library |
|------------|--------------|------------------|
| `arr.map(fn)` | `arr.Map(fn)` | N/A |
| `obj.get(key)` | `m.Get(key)` | `m[key]` |
| `set.has(val)` | `s.Has(val)` | N/A |
| `map.entries()` | `m.Entries()` | N/A |

---

<sub><sup>Caution: This library was vibe-coded—use it at your own risk!</sup></sub>