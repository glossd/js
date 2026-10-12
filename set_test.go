package js

import "testing"

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
