package js

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
