package js

import "iter"

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	s := make(Set[T], len(values))
	for _, v := range values {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[T]) Size() int {
	return len(s)
}

// Add inserts v into the set and returns the set to allow chaining,
// like JS `set.add(v)`. Adding an existing value is a no-op.
func (s Set[T]) Add(v T) Set[T] {
	s[v] = struct{}{}
	return s
}

// Delete removes v from the set. It reports whether v was present,
// like JS `set.delete(v)`.
func (s Set[T]) Delete(v T) bool {
	_, ok := s[v]
	delete(s, v)
	return ok
}

// Has reports whether v is in the set, like JS `set.has(v)`.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Clear removes all values from the set, like JS `set.clear()`.
func (s Set[T]) Clear() {
	clear(s)
}

// ForEach executes a provided function once for each value in the set,
// like JS Set.prototype.forEach. The receiver is never modified.
func (s Set[T]) ForEach(callbackFn func(value T)) {
	if s == nil || len(s) == 0 || callbackFn == nil {
		return
	}

	for v := range s {
		callbackFn(v)
	}
}

// IsSubsetOf reports whether every value of s is also in other,
// like JS `set.isSubsetOf(other)`. The empty set is a subset of every set,
// and every set is a subset of itself.
func (s Set[T]) IsSubsetOf(other Set[T]) bool {
	if len(s) > len(other) {
		return false
	}
	for v := range s {
		if !other.Has(v) {
			return false
		}
	}
	return true
}

// IsSupersetOf reports whether every value of other is also in s,
// like JS `set.isSupersetOf(other)`. Every set is a superset of the empty
// set and of itself.
func (s Set[T]) IsSupersetOf(other Set[T]) bool {
	return other.IsSubsetOf(s)
}

// IsDisjointFrom reports whether s and other have no values in common,
// like JS `set.isDisjointFrom(other)`. The empty set is disjoint from
// every set.
func (s Set[T]) IsDisjointFrom(other Set[T]) bool {
	// Iterate the smaller set and probe the larger one.
	small, large := s, other
	if len(small) > len(large) {
		small, large = large, small
	}
	for v := range small {
		if large.Has(v) {
			return false
		}
	}
	return true
}

// Union returns a new set containing every value that is in s or in other,
// like JS `set.union(other)`. Neither s nor other is modified.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := make(Set[T], len(s)+len(other))
	for v := range s {
		result.Add(v)
	}
	for v := range other {
		result.Add(v)
	}
	return result
}

// Intersection returns a new set containing only the values that are in both
// s and other, like JS `set.intersection(other)`. Neither s nor other is
// modified.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	// Iterate the smaller set and probe the larger one.
	small, large := s, other
	if len(small) > len(large) {
		small, large = large, small
	}
	result := Set[T]{}
	for v := range small {
		if large.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns a new set containing the values that are in s but not in
// other, like JS `set.difference(other)`. Neither s nor other is modified.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := Set[T]{}
	for v := range s {
		if !other.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// SymmetricDifference returns a new set containing the values that are in
// exactly one of s and other, like JS `set.symmetricDifference(other)`.
// Neither s nor other is modified.
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	result := Set[T]{}
	for v := range s {
		if !other.Has(v) {
			result.Add(v)
		}
	}
	for v := range other {
		if !s.Has(v) {
			result.Add(v)
		}
	}
	return result
}

// Values returns an iterator over the values of the set, like JS
// `set.values()`. Use it with range-over-func:
//
//	for v := range s.Values() { ... }
//
// Unlike JS, which iterates in insertion order, the iteration order is
// unspecified (it follows Go map iteration order) and may differ between
// calls. Stopping early with break is safe, and values deleted from the set
// during iteration that have not been reached yet will not be produced.
func (s Set[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
