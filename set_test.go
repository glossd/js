package js

import (
	"slices"
	"strings"
	"testing"
)

func TestSet_Readme(t *testing.T) {
	fruits := NewSet("apple", "banana", "avocado")

	aFruits := ArrayFrom(fruits.Values()).Filter(func(v string) bool { return strings.HasPrefix(v, "a") })

	if !slices.Equal(ArrayOf("apple", "avocado"), aFruits) {
		t.Errorf("aFruits not equal")
	}
}

func TestSet_Size(t *testing.T) {
	if NewSet(1, 2, 3).Size() != 3 {
		t.Fail()
	}
}

func TestSet_AddAndHas(t *testing.T) {
	s := Set[int]{}

	if s.Has(1) {
		t.Fatal("empty set should not contain 1")
	}

	s.Add(1)
	if !s.Has(1) {
		t.Fatal("set should contain 1 after Add")
	}
	if s.Has(2) {
		t.Fatal("set should not contain 2")
	}
}

func TestSet_AddIsIdempotent(t *testing.T) {
	s := Set[string]{}
	s.Add("a")
	s.Add("a")
	s.Add("a")

	if len(s) != 1 {
		t.Fatalf("expected len 1 after adding duplicates, got %d", len(s))
	}
}

func TestSet_AddChaining(t *testing.T) {
	s := Set[int]{}
	s.Add(1).Add(2).Add(3)

	if len(s) != 3 {
		t.Fatalf("expected len 3 after chained Add, got %d", len(s))
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Has(v) {
			t.Fatalf("expected set to contain %d", v)
		}
	}
}

func TestSet_AddReturnsSameSet(t *testing.T) {
	s := Set[int]{}
	r := s.Add(1)

	r.Add(2)
	if !s.Has(2) {
		t.Fatal("Add should return the same underlying set")
	}
}

func TestSet_Delete(t *testing.T) {
	s := Set[int]{}
	s.Add(1).Add(2)

	if !s.Delete(1) {
		t.Fatal("Delete of existing value should return true")
	}
	if s.Has(1) {
		t.Fatal("value should be gone after Delete")
	}
	if !s.Has(2) {
		t.Fatal("Delete should not affect other values")
	}
}

func TestSet_DeleteMissing(t *testing.T) {
	s := Set[int]{}
	s.Add(1)

	if s.Delete(99) {
		t.Fatal("Delete of missing value should return false")
	}
	if len(s) != 1 {
		t.Fatalf("set should be unchanged, got len %d", len(s))
	}
}

func TestSet_DeleteTwice(t *testing.T) {
	s := Set[int]{}
	s.Add(1)

	if !s.Delete(1) {
		t.Fatal("first Delete should return true")
	}
	if s.Delete(1) {
		t.Fatal("second Delete should return false")
	}
}

func TestSet_Clear(t *testing.T) {
	s := Set[int]{}
	s.Add(1).Add(2).Add(3)

	s.Clear()

	if len(s) != 0 {
		t.Fatalf("expected empty set after Clear, got len %d", len(s))
	}
	for _, v := range []int{1, 2, 3} {
		if s.Has(v) {
			t.Fatalf("set should not contain %d after Clear", v)
		}
	}
}

func TestSet_ClearEmpty(t *testing.T) {
	s := Set[int]{}
	s.Clear() // must not panic

	if len(s) != 0 {
		t.Fatal("empty set should remain empty")
	}
}

func TestSet_ReuseAfterClear(t *testing.T) {
	s := Set[string]{}
	s.Add("a")
	s.Clear()
	s.Add("b")

	if s.Has("a") || !s.Has("b") || len(s) != 1 {
		t.Fatalf("unexpected state after reuse: %v", s)
	}
}

func TestSet_StructKey(t *testing.T) {
	type point struct{ X, Y int }

	s := Set[point]{}
	s.Add(point{1, 2})

	if !s.Has(point{1, 2}) {
		t.Fatal("struct keys should compare by value")
	}
	if s.Has(point{2, 1}) {
		t.Fatal("different struct value should not be present")
	}
}

func TestSet_NaNKey(t *testing.T) {
	// Go maps treat NaN != NaN, so each Add inserts a new entry and
	// Has/Delete can never find it. (JS Set uses SameValueZero, where
	// NaN equals NaN.) This test documents the difference.
	nan := 0.0
	nan = nan / nan

	s := Set[float64]{}
	s.Add(nan)

	if s.Has(nan) {
		t.Fatal("Go map semantics: NaN is never found")
	}
}

func TestSet_IsSubsetOf(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected bool
	}{
		{"proper subset", NewSet(1, 2), NewSet(1, 2, 3), true},
		{"equal sets", NewSet(1, 2, 3), NewSet(1, 2, 3), true},
		{"empty of non-empty", NewSet[int](), NewSet(1), true},
		{"empty of empty", NewSet[int](), NewSet[int](), true},
		{"non-empty of empty", NewSet(1), NewSet[int](), false},
		{"superset is not subset", NewSet(1, 2, 3), NewSet(1, 2), false},
		{"partial overlap", NewSet(1, 4), NewSet(1, 2, 3), false},
		{"same size, different values", NewSet(1, 2), NewSet(3, 4), false},
		{"disjoint", NewSet(1), NewSet(2, 3), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsSubsetOf(tt.b); got != tt.expected {
				t.Fatalf("%v.IsSubsetOf(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestSet_IsSupersetOf(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected bool
	}{
		{"proper superset", NewSet(1, 2, 3), NewSet(1, 2), true},
		{"equal sets", NewSet(1, 2, 3), NewSet(1, 2, 3), true},
		{"non-empty of empty", NewSet(1), NewSet[int](), true},
		{"empty of empty", NewSet[int](), NewSet[int](), true},
		{"empty of non-empty", NewSet[int](), NewSet(1), false},
		{"subset is not superset", NewSet(1, 2), NewSet(1, 2, 3), false},
		{"partial overlap", NewSet(1, 2, 3), NewSet(1, 4), false},
		{"same size, different values", NewSet(1, 2), NewSet(3, 4), false},
		{"disjoint", NewSet(2, 3), NewSet(1), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsSupersetOf(tt.b); got != tt.expected {
				t.Fatalf("%v.IsSupersetOf(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestSet_IsDisjointFrom(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected bool
	}{
		{"no common values", NewSet(1, 2), NewSet(3, 4), true},
		{"one common value", NewSet(1, 2), NewSet(2, 3), false},
		{"equal sets", NewSet(1, 2), NewSet(1, 2), false},
		{"subset", NewSet(1), NewSet(1, 2, 3), false},
		{"superset", NewSet(1, 2, 3), NewSet(2), false},
		{"empty and non-empty", NewSet[int](), NewSet(1), true},
		{"non-empty and empty", NewSet(1), NewSet[int](), true},
		{"empty and empty", NewSet[int](), NewSet[int](), true},
		{"small disjoint from large", NewSet(100), NewSet(1, 2, 3, 4, 5), true},
		{"large disjoint from small", NewSet(1, 2, 3, 4, 5), NewSet(100), true},
		{"small overlaps large", NewSet(5), NewSet(1, 2, 3, 4, 5), false},
		{"large overlaps small", NewSet(1, 2, 3, 4, 5), NewSet(5), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsDisjointFrom(tt.b); got != tt.expected {
				t.Fatalf("%v.IsDisjointFrom(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestSet_SetRelationsAreSymmetricWhereExpected(t *testing.T) {
	a, b := NewSet(1, 2), NewSet(1, 2, 3)

	if !a.IsSubsetOf(b) || !b.IsSupersetOf(a) {
		t.Fatal("a ⊂ b should imply b ⊃ a")
	}
	if a.IsDisjointFrom(b) != b.IsDisjointFrom(a) {
		t.Fatal("IsDisjointFrom should be symmetric")
	}
}

func TestSet_RelationsDoNotMutate(t *testing.T) {
	a, b := NewSet(1, 2), NewSet(2, 3)

	a.IsSubsetOf(b)
	a.IsSupersetOf(b)
	a.IsDisjointFrom(b)

	if len(a) != 2 || !a.Has(1) || !a.Has(2) {
		t.Fatalf("receiver was mutated: %v", a)
	}
	if len(b) != 2 || !b.Has(2) || !b.Has(3) {
		t.Fatalf("argument was mutated: %v", b)
	}
}

func TestSet_RelationsWithNilSet(t *testing.T) {
	var nilSet Set[int]
	s := NewSet(1, 2)

	if !nilSet.IsSubsetOf(s) {
		t.Fatal("nil set is a subset of any set")
	}
	if !s.IsSupersetOf(nilSet) {
		t.Fatal("any set is a superset of nil set")
	}
	if !nilSet.IsDisjointFrom(s) || !s.IsDisjointFrom(nilSet) {
		t.Fatal("nil set is disjoint from any set")
	}
	if s.IsSubsetOf(nilSet) {
		t.Fatal("non-empty set is not a subset of nil set")
	}
}

func equalSets[T comparable](a, b Set[T]) bool {
	return a.IsSubsetOf(b) && b.IsSubsetOf(a)
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected Set[int]
	}{
		{"disjoint", NewSet(1, 2), NewSet(3, 4), NewSet(1, 2, 3, 4)},
		{"overlapping", NewSet(1, 2, 3), NewSet(3, 4), NewSet(1, 2, 3, 4)},
		{"equal sets", NewSet(1, 2), NewSet(1, 2), NewSet(1, 2)},
		{"subset", NewSet(1), NewSet(1, 2, 3), NewSet(1, 2, 3)},
		{"empty and non-empty", NewSet[int](), NewSet(1, 2), NewSet(1, 2)},
		{"non-empty and empty", NewSet(1, 2), NewSet[int](), NewSet(1, 2)},
		{"empty and empty", NewSet[int](), NewSet[int](), NewSet[int]()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Union(tt.b)
			if !equalSets(got, tt.expected) {
				t.Fatalf("%v.Union(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
			if rev := tt.b.Union(tt.a); !equalSets(rev, tt.expected) {
				t.Fatalf("Union should be commutative: got %v, want %v", rev, tt.expected)
			}
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected Set[int]
	}{
		{"disjoint", NewSet(1, 2), NewSet(3, 4), NewSet[int]()},
		{"overlapping", NewSet(1, 2, 3), NewSet(3, 4, 2), NewSet(2, 3)},
		{"equal sets", NewSet(1, 2), NewSet(1, 2), NewSet(1, 2)},
		{"subset", NewSet(1), NewSet(1, 2, 3), NewSet(1)},
		{"superset", NewSet(1, 2, 3), NewSet(1), NewSet(1)},
		{"empty and non-empty", NewSet[int](), NewSet(1, 2), NewSet[int]()},
		{"non-empty and empty", NewSet(1, 2), NewSet[int](), NewSet[int]()},
		{"empty and empty", NewSet[int](), NewSet[int](), NewSet[int]()},
		{"small in large", NewSet(5), NewSet(1, 2, 3, 4, 5), NewSet(5)},
		{"large with small", NewSet(1, 2, 3, 4, 5), NewSet(5), NewSet(5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Intersection(tt.b)
			if !equalSets(got, tt.expected) {
				t.Fatalf("%v.Intersection(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
			if rev := tt.b.Intersection(tt.a); !equalSets(rev, tt.expected) {
				t.Fatalf("Intersection should be commutative: got %v, want %v", rev, tt.expected)
			}
		})
	}
}

func TestSet_UnionIntersectionReturnNewSets(t *testing.T) {
	a, b := NewSet(1, 2), NewSet(2, 3)

	u := a.Union(b)
	i := a.Intersection(b)

	// Inputs are untouched.
	if !equalSets(a, NewSet(1, 2)) || !equalSets(b, NewSet(2, 3)) {
		t.Fatalf("inputs were mutated: a=%v b=%v", a, b)
	}

	// Results are independent of the inputs and of each other.
	u.Add(99)
	i.Add(98)
	if a.Has(99) || b.Has(99) || a.Has(98) || b.Has(98) {
		t.Fatal("mutating a result must not affect the inputs")
	}
	if u.Has(98) || i.Has(99) {
		t.Fatal("results must be independent of each other")
	}
}

func TestSet_UnionIntersectionWithSelfAndNil(t *testing.T) {
	var nilSet Set[int]
	s := NewSet(1, 2)

	if !equalSets(s.Union(s), s) || !equalSets(s.Intersection(s), s) {
		t.Fatal("union and intersection with self should equal self")
	}

	u := nilSet.Union(s)
	if !equalSets(u, s) {
		t.Fatalf("nil.Union(s) = %v, want %v", u, s)
	}
	// The result must be a usable, writable set even when built from nil.
	u.Add(3)

	n := nilSet.Union(nilSet)
	if n == nil || len(n) != 0 {
		t.Fatalf("nil.Union(nil) should be a non-nil empty set, got %#v", n)
	}
	n.Add(1) // must not panic

	if i := nilSet.Intersection(s); i == nil || len(i) != 0 {
		t.Fatalf("nil.Intersection(s) should be a non-nil empty set, got %#v", i)
	}
}

func TestSet_AlgebraicLaws(t *testing.T) {
	a, b := NewSet(1, 2, 3), NewSet(3, 4)

	u, i := a.Union(b), a.Intersection(b)

	if !a.IsSubsetOf(u) || !b.IsSubsetOf(u) {
		t.Fatal("both operands should be subsets of their union")
	}
	if !i.IsSubsetOf(a) || !i.IsSubsetOf(b) {
		t.Fatal("the intersection should be a subset of both operands")
	}
	if len(u)+len(i) != len(a)+len(b) {
		t.Fatalf("|A∪B| + |A∩B| should equal |A| + |B|: %d + %d != %d + %d", len(u), len(i), len(a), len(b))
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected Set[int]
	}{
		{"disjoint", NewSet(1, 2), NewSet(3, 4), NewSet(1, 2)},
		{"overlapping", NewSet(1, 2, 3), NewSet(2, 3, 4), NewSet(1)},
		{"equal sets", NewSet(1, 2), NewSet(1, 2), NewSet[int]()},
		{"subset minus superset", NewSet(1, 2), NewSet(1, 2, 3), NewSet[int]()},
		{"superset minus subset", NewSet(1, 2, 3), NewSet(2), NewSet(1, 3)},
		{"empty minus non-empty", NewSet[int](), NewSet(1), NewSet[int]()},
		{"non-empty minus empty", NewSet(1, 2), NewSet[int](), NewSet(1, 2)},
		{"empty minus empty", NewSet[int](), NewSet[int](), NewSet[int]()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Difference(tt.b)
			if !equalSets(got, tt.expected) {
				t.Fatalf("%v.Difference(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestSet_DifferenceIsNotCommutative(t *testing.T) {
	a, b := NewSet(1, 2, 3), NewSet(3, 4)

	if got := a.Difference(b); !equalSets(got, NewSet(1, 2)) {
		t.Fatalf("a-b = %v, want {1 2}", got)
	}
	if got := b.Difference(a); !equalSets(got, NewSet(4)) {
		t.Fatalf("b-a = %v, want {4}", got)
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Set[int]
		expected Set[int]
	}{
		{"disjoint", NewSet(1, 2), NewSet(3, 4), NewSet(1, 2, 3, 4)},
		{"overlapping", NewSet(1, 2, 3), NewSet(3, 4), NewSet(1, 2, 4)},
		{"equal sets", NewSet(1, 2), NewSet(1, 2), NewSet[int]()},
		{"subset", NewSet(1, 2), NewSet(1, 2, 3), NewSet(3)},
		{"superset", NewSet(1, 2, 3), NewSet(2), NewSet(1, 3)},
		{"empty and non-empty", NewSet[int](), NewSet(1, 2), NewSet(1, 2)},
		{"non-empty and empty", NewSet(1, 2), NewSet[int](), NewSet(1, 2)},
		{"empty and empty", NewSet[int](), NewSet[int](), NewSet[int]()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.SymmetricDifference(tt.b)
			if !equalSets(got, tt.expected) {
				t.Fatalf("%v.SymmetricDifference(%v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
			if rev := tt.b.SymmetricDifference(tt.a); !equalSets(rev, tt.expected) {
				t.Fatalf("SymmetricDifference should be commutative: got %v, want %v", rev, tt.expected)
			}
		})
	}
}

func TestSet_DifferenceReturnNewSets(t *testing.T) {
	a, b := NewSet(1, 2, 3), NewSet(3, 4)

	d := a.Difference(b)
	sd := a.SymmetricDifference(b)

	if !equalSets(a, NewSet(1, 2, 3)) || !equalSets(b, NewSet(3, 4)) {
		t.Fatalf("inputs were mutated: a=%v b=%v", a, b)
	}

	d.Add(99)
	sd.Add(98)
	if a.Has(99) || b.Has(99) || a.Has(98) || b.Has(98) {
		t.Fatal("mutating a result must not affect the inputs")
	}
	if d.Has(98) || sd.Has(99) {
		t.Fatal("results must be independent of each other")
	}
}

func TestSet_DifferenceWithSelfAndNil(t *testing.T) {
	var nilSet Set[int]
	s := NewSet(1, 2)

	if d := s.Difference(s); d == nil || len(d) != 0 {
		t.Fatalf("s.Difference(s) should be a non-nil empty set, got %#v", d)
	}
	if sd := s.SymmetricDifference(s); sd == nil || len(sd) != 0 {
		t.Fatalf("s.SymmetricDifference(s) should be a non-nil empty set, got %#v", sd)
	}

	if d := nilSet.Difference(s); d == nil || len(d) != 0 {
		t.Fatalf("nil.Difference(s) should be a non-nil empty set, got %#v", d)
	}
	if d := s.Difference(nilSet); !equalSets(d, s) {
		t.Fatalf("s.Difference(nil) = %v, want %v", d, s)
	}

	sd := nilSet.SymmetricDifference(s)
	if !equalSets(sd, s) {
		t.Fatalf("nil.SymmetricDifference(s) = %v, want %v", sd, s)
	}
	sd.Add(3) // result must be writable

	if sd := nilSet.SymmetricDifference(nilSet); sd == nil || len(sd) != 0 {
		t.Fatalf("nil.SymmetricDifference(nil) should be a non-nil empty set, got %#v", sd)
	}
}

func TestSet_DifferenceIdentities(t *testing.T) {
	a, b := NewSet(1, 2, 3), NewSet(3, 4)

	// A △ B == (A \ B) ∪ (B \ A) == (A ∪ B) \ (A ∩ B)
	sd := a.SymmetricDifference(b)
	viaDiffs := a.Difference(b).Union(b.Difference(a))
	viaUnion := a.Union(b).Difference(a.Intersection(b))

	if !equalSets(sd, viaDiffs) {
		t.Fatalf("A△B = %v, but (A\\B)∪(B\\A) = %v", sd, viaDiffs)
	}
	if !equalSets(sd, viaUnion) {
		t.Fatalf("A△B = %v, but (A∪B)\\(A∩B) = %v", sd, viaUnion)
	}
	if !a.Difference(b).IsDisjointFrom(b) {
		t.Fatal("A\\B should be disjoint from B")
	}
}

func TestSet_Values(t *testing.T) {
	s := NewSet(3, 1, 2)

	var got []int
	for v := range s.Values() {
		got = append(got, v)
	}
	slices.Sort(got) // iteration order is unspecified

	if want := []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Fatalf("Values() yielded %v, want %v", got, want)
	}
}

func TestSet_ValuesEmptyAndNil(t *testing.T) {
	var nilSet Set[int]

	for _, s := range []Set[int]{nilSet, NewSet[int]()} {
		for v := range s.Values() {
			t.Fatalf("empty set yielded %v", v)
		}
		if got := slices.Collect(s.Values()); len(got) != 0 {
			t.Fatalf("Collect on empty set = %v, want empty", got)
		}
	}
}

func TestSet_ValuesYieldsEachValueOnce(t *testing.T) {
	s := Set[int]{}
	for i := range 1000 {
		s.Add(i % 100) // duplicates collapse
	}

	seen := map[int]int{}
	for v := range s.Values() {
		seen[v]++
	}

	if len(seen) != 100 {
		t.Fatalf("yielded %d distinct values, want 100", len(seen))
	}
	for v, n := range seen {
		if n != 1 {
			t.Fatalf("value %d yielded %d times, want 1", v, n)
		}
	}
}

func TestSet_ValuesEarlyBreak(t *testing.T) {
	s := NewSet(1, 2, 3, 4, 5)

	count := 0
	for range s.Values() {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Fatalf("expected loop to stop after 2 values, ran %d", count)
	}

	// Calling the iterator function directly: it must stop once yield
	// returns false.
	calls := 0
	s.Values()(func(int) bool {
		calls++
		return false
	})
	if calls != 1 {
		t.Fatalf("yield called %d times after returning false, want 1", calls)
	}
}

func TestSet_ValuesIsReusable(t *testing.T) {
	s := NewSet(1, 2, 3)
	seq := s.Values()

	first := slices.Sorted(seq)
	second := slices.Sorted(seq)

	if !slices.Equal(first, second) || len(first) != 3 {
		t.Fatalf("iterator should be re-iterable: %v then %v", first, second)
	}
}

func TestSet_ValuesSeesLiveSet(t *testing.T) {
	s := NewSet(1, 2)
	seq := s.Values()

	s.Add(3)

	if got := slices.Sorted(seq); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("iterator should reflect the current set, got %v", got)
	}
}

func TestSet_ValuesDeleteDuringIteration(t *testing.T) {
	s := NewSet(1, 2, 3, 4, 5)

	// Deleting every visited value while iterating is safe and visits
	// each value exactly once.
	visited := 0
	for v := range s.Values() {
		visited++
		s.Delete(v)
	}

	if visited != 5 {
		t.Fatalf("visited %d values, want 5", visited)
	}
	if len(s) != 0 {
		t.Fatalf("set should be empty after deleting during iteration, got %v", s)
	}
}

func TestSet_ValuesRoundTrip(t *testing.T) {
	a := NewSet("x", "y", "z")

	b := Set[string]{}
	for v := range a.Values() {
		b.Add(v)
	}

	if !equalSets(a, b) {
		t.Fatalf("rebuilding from Values() gave %v, want %v", b, a)
	}
}
