package jss

import "iter"

// In Go, the map key must be comparable

type Map[K comparable, V any] map[K]V

func NewMap[K comparable, V any](entries ...Entry[K, V]) Map[K, V] {
	var m Map[K, V]
	if len(entries) == 0 {
		m = make(Map[K, V])
	} else {
		m = make(Map[K, V], len(entries))
	}

	for _, e := range entries {
		m[e.Key] = e.Val
	}
	return m
}

// Set adds a new element with (key, value) specified or updates an existing element.
// The receiver is modified and returned, allowing chaining, and mirrors JavaScript's
// return of the Map object itself.
func (m *Map[K, V]) Set(key K, value V) *Map[K, V] {
	if *m == nil {
		*m = make(Map[K, V])
	}
	(*m)[key] = value
	return m
}

// Get returns the value associated with the key, or the zero value of V and false
// if the key is not present. It returns false instead of undefined to distinguish
// "not found" from a stored zero value, mirroring the (value, ok) idiom used by
// Array.Find and Array.At.
func (m Map[K, V]) Get(key K) (V, bool) {
	val, ok := m[key]
	return val, ok
}

// Has returns true if the map contains an element with the specified key, and false otherwise.
// The receiver is never modified.
func (m Map[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

// Delete removes the element with the specified key and returns true if an element was removed.
// Returns false if the key was not found. The receiver is modified.
func (m *Map[K, V]) Delete(key K) bool {
	if *m == nil {
		return false
	}
	_, ok := (*m)[key]
	if ok {
		delete(*m, key)
	}
	return ok
}

// Clear removes all elements from the map. The receiver is modified and returned
// for chaining, mirroring JavaScript's Map.prototype.clear (which returns undefined).
func (m *Map[K, V]) Clear() *Map[K, V] {
	if *m != nil {
		clear(*m)
	}
	return m
}

// Keys returns an iterator over the keys of the map.
// The receiver is never modified.
func (m Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over the values of the map.
// The receiver is never modified.
func (m Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				return
			}
		}
	}
}

// Entries returns an iterator over the key-value entries of the map.
// The receiver is never modified.
func (m Map[K, V]) Entries() iter.Seq[Entry[K, V]] {
	return func(yield func(Entry[K, V]) bool) {
		for k, v := range m {
			if !yield(Entry[K, V]{Key: k, Val: v}) {
				return
			}
		}
	}
}

func (m Map[K, V]) Size() int {
	return len(m)
}
