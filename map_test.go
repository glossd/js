package js

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewMap(t *testing.T) {
	m := NewMap[string, int]()

	m = NewMap(Entry[string, int]{"mykey", 123})
	if 123 != m["mykey"] {
		t.Errorf("expected mykey to be 123, got %d", m["mykey"])
	}

	m = NewMap(E("mykey", 123))
	if 123 != m["mykey"] {
		t.Errorf("expected mykey to be 123, got %d", m["mykey"])
	}

	m = NewMap([]Entry[string, int]{{"k1", 1}, {"k2", 2}}...)
	if 2 != m["k2"] {
		t.Errorf("expected mykey to be 2, got %d", m["k2"])
	}

	m = NewMap(E2("k1", 1, "k2", 2)...)
	if 2 != m["k2"] {
		t.Errorf("expected mykey to be 2, got %d", m["k2"])
	}

	myMap := NewMap(E2("key1", 12, "key2", 123)...)
	sum := ArrayFrom(myMap.Values()).Reduce(func(acc int, e int) int { return acc + e })
	if sum != 135 {
		t.Errorf("expected sum to be 135, got %d", sum)
	}

}

func TestSet(t *testing.T) {
	m := NewMap[string, int]()

	// Test setting a new key
	m.Set("newKey", 42)
	if val, ok := m.Get("newKey"); !ok || val != 42 {
		t.Errorf("expected newKey to be 42, got %d", val)
	}

	// Test updating an existing key
	m.Set("existing", 100)
	if val, ok := m.Get("existing"); !ok || val != 100 {
		t.Errorf("expected existing to be 100, got %d", val)
	}

	// Test setting with different value
	m.Set("existing", 999)
	if val, ok := m.Get("existing"); !ok || val != 999 {
		t.Errorf("expected existing to be 999, got %d", val)
	}
}

func TestGet(t *testing.T) {
	m := NewMap[string, int]()

	// Test getting an existing key
	m.Set("key1", 10)
	if val, ok := m.Get("key1"); !ok || val != 10 {
		t.Errorf("expected key1 to be 10, got %v, %v", val, ok)
	}

	// Test getting a non-existing key
	if val, ok := m.Get("nonexistent"); ok {
		t.Errorf("expected nonexistent to be false, got %v", val)
	}
}

func TestHas(t *testing.T) {
	m := NewMap[string, int]()

	// Test checking an existing key
	m.Set("key", 5)
	if m.Has("key") != true {
		t.Errorf("expected key to be present, got %v", m.Has("key"))
	}

	// Test checking a non-existing key
	if m.Has("missing") != false {
		t.Errorf("expected missing to be absent, got %v", m.Has("missing"))
	}

	// Test on empty map
	emptyMap := NewMap[string, int]()
	if emptyMap.Has("anything") != false {
		t.Errorf("expected empty map to not have anything, got %v", emptyMap.Has("anything"))
	}
}

func TestDelete(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Test deleting an existing key
	if !m.Delete("a") {
		t.Error("expected Delete to return true for existing key")
	}
	if m.Has("a") {
		t.Error("expected key a to be deleted")
	}

	// Test deleting a non-existing key
	if m.Delete("nonexistent") {
		t.Error("expected Delete to return false for missing key")
	}

	// Test deleting from nil map
	var nilMap Map[string, int]
	if nilMap.Delete("x") {
		t.Error("expected Delete to return false for nil map")
	}
}

func TestClear(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Test clearing a populated map
	m.Clear()
	if m.Size() != 0 {
		t.Errorf("expected map size to be 0 after Clear, got %d", m.Size())
	}
	if m.Has("a") {
		t.Error("expected key a to be gone after Clear")
	}

	// Test clearing an empty map
	emptyMap := NewMap[string, int]()
	emptyMap.Clear()
	if emptyMap.Size() != 0 {
		t.Errorf("expected empty map size to remain 0, got %d", emptyMap.Size())
	}

	// Test clearing a nil map
	var nilMap Map[string, int]
	nilMap.Clear()
	if nilMap.Size() != 0 {
		t.Errorf("expected nil map size to be 0, got %d", nilMap.Size())
	}
}

func TestKeys(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Collect keys using range over the iterator
	keys := make([]string, 0, m.Size())
	for k := range m.Keys() {
		keys = append(keys, k)
	}

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	// Verify all keys are present (order not guaranteed in Go maps)
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}
	for _, expected := range []string{"a", "b", "c"} {
		if !keySet[expected] {
			t.Errorf("expected key %q to be present", expected)
		}
	}

	// Test keys on empty map
	emptyMap := NewMap[string, int]()
	count := 0
	for range emptyMap.Keys() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 keys for empty map, got %d", count)
	}

	// Test keys on nil map
	var nilMap2 Map[string, int]
	count = 0
	for range nilMap2.Keys() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 keys for nil map, got %d", count)
	}

	// Test early termination with break (yield returns false)
	count = 0
	for range m.Keys() {
		count++
		if count >= 2 {
			break
		}
	}
	// Test early termination with break (yield returns false)
	count = 0
	for range m.Keys() {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("expected early break after 2, got %d", count)
	}
}

func TestValues(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Collect values using range over the iterator
	values := make([]int, 0, m.Size())
	for v := range m.Values() {
		values = append(values, v)
	}

	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}

	// Verify all values are present
	valueSet := make(map[int]bool)
	for _, v := range values {
		valueSet[v] = true
	}
	for _, expected := range []int{1, 2, 3} {
		if !valueSet[expected] {
			t.Errorf("expected value %d to be present", expected)
		}
	}

	// Test values on empty map
	emptyMap := NewMap[string, int]()
	count := 0
	for range emptyMap.Values() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 values for empty map, got %d", count)
	}

	// Test values on nil map
	var nilMap2 Map[string, int]
	count = 0
	for range nilMap2.Values() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 values for nil map, got %d", count)
	}

	// Test early termination with break (yield returns false)
	count = 0
	for range m.Values() {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("expected early break after 2, got %d", count)
	}
}

func TestEntries(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Collect entries using range over the iterator
	entries := make([]Entry[string, int], 0, m.Size())
	for e := range m.Entries() {
		entries = append(entries, e)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	// Verify all entries are present (order not guaranteed)
	entrySet := make(map[string]int)
	for _, e := range entries {
		entrySet[e.Key] = e.Val
	}
	for _, key := range []string{"a", "b", "c"} {
		if _, ok := entrySet[key]; !ok {
			t.Errorf("expected entry with key %q to be present", key)
		}
	}

	// Test specific entry
	if val, ok := entrySet["a"]; !ok || val != 1 {
		t.Errorf("expected entry 'a' with value 1, got %v, %v", val, ok)
	}

	// Test entries on empty map
	emptyMap := NewMap[string, int]()
	count := 0
	for range emptyMap.Entries() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries for empty map, got %d", count)
	}

	// Test entries on nil map
	var nilMap2 Map[string, int]
	count = 0
	for range nilMap2.Entries() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries for nil map, got %d", count)
	}

	// Test early termination with break (yield returns false)
	count = 0
	for range m.Entries() {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("expected early break after 2, got %d", count)
	}
}

func TestMapForEach(t *testing.T) {
	t.Run("visits every entry exactly once", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2, "c": 3}

		visited := map[string]int{}
		calls := 0
		m.ForEach(func(v int, k string) {
			calls++
			visited[k] = v
		})

		if calls != len(m) {
			t.Fatalf("expected %d calls, got %d", len(m), calls)
		}
		if !reflect.DeepEqual(visited, map[string]int(m)) {
			t.Fatalf("expected %v, got %v", m, visited)
		}
	})

	t.Run("passes value first and key second", func(t *testing.T) {
		m := Map[string, int]{"only": 42}

		var gotKey string
		var gotValue int
		m.ForEach(func(v int, k string) {
			gotValue, gotKey = v, k
		})

		if gotKey != "only" || gotValue != 42 {
			t.Fatalf("expected (42, %q), got (%d, %q)", "only", gotValue, gotKey)
		}
	})

	t.Run("empty map never calls fn", func(t *testing.T) {
		m := Map[string, int]{}

		m.ForEach(func(v int, k string) {
			t.Fatalf("callback should not be called, got (%d, %q)", v, k)
		})
	})

	t.Run("nil map does not panic or call fn", func(t *testing.T) {
		var m Map[string, int]

		m.ForEach(func(v int, k string) {
			t.Fatalf("callback should not be called, got (%d, %q)", v, k)
		})
	})

	t.Run("works with non-string keys and struct values", func(t *testing.T) {
		type point struct{ x, y int }
		m := Map[int, point]{1: {1, 2}, 2: {3, 4}}

		var got []Entry[int, point]
		m.ForEach(func(v point, k int) {
			got = append(got, Entry[int, point]{k, v})
		})
		sort.Slice(got, func(i, j int) bool { return got[i].Key < got[j].Key })

		want := []Entry[int, point]{{1, point{1, 2}}, {2, point{3, 4}}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})

	t.Run("does not modify the map", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2}
		before := map[string]int{"a": 1, "b": 2}

		m.ForEach(func(v int, k string) {})

		if !reflect.DeepEqual(map[string]int(m), before) {
			t.Fatalf("map changed: expected %v, got %v", before, m)
		}
	})

	t.Run("callback can update values of existing keys", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2}

		m.ForEach(func(v int, k string) {
			m[k] = v * 10
		})

		want := Map[string, int]{"a": 10, "b": 20}
		if !reflect.DeepEqual(m, want) {
			t.Fatalf("expected %v, got %v", want, m)
		}
	})

	t.Run("callback can delete entries during iteration", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2, "c": 3}

		m.ForEach(func(v int, k string) {
			delete(m, k)
		})

		if len(m) != 0 {
			t.Fatalf("expected empty map, got %v", m)
		}
	})
}

func TestMapForEachUntil(t *testing.T) {
	t.Run("visits all entries when fn always returns true", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2, "c": 3}

		calls := 0
		m.ForEachUntil(func(v int, k string) bool {
			calls++
			return true
		})

		if calls != len(m) {
			t.Fatalf("expected %d calls, got %d", len(m), calls)
		}
	})

	t.Run("stops as soon as fn returns false", func(t *testing.T) {
		m := Map[string, int]{"a": 1, "b": 2, "c": 3, "d": 4}

		calls := 0
		m.ForEachUntil(func(v int, k string) bool {
			calls++
			return false
		})

		if calls != 1 {
			t.Fatalf("expected exactly 1 call, got %d", calls)
		}
	})

	t.Run("stops after the entry that returns false", func(t *testing.T) {
		m := Map[int, int]{}
		for i := 0; i < 100; i++ {
			m[i] = i
		}

		const limit = 10
		calls := 0
		m.ForEachUntil(func(v int, k int) bool {
			calls++
			return calls < limit
		})

		if calls != limit {
			t.Fatalf("expected %d calls, got %d", limit, calls)
		}
	})

	t.Run("passes value first and key second", func(t *testing.T) {
		m := Map[string, int]{"only": 7}

		var gotKey string
		var gotValue int
		m.ForEachUntil(func(v int, k string) bool {
			gotValue, gotKey = v, k
			return true
		})

		if gotKey != "only" || gotValue != 7 {
			t.Fatalf("expected (7, %q), got (%d, %q)", "only", gotValue, gotKey)
		}
	})

	t.Run("empty and nil maps never call fn", func(t *testing.T) {
		var nilMap Map[string, int]
		for name, m := range map[string]Map[string, int]{
			"empty": {},
			"nil":   nilMap,
		} {
			m.ForEachUntil(func(v int, k string) bool {
				t.Fatalf("%s map: callback should not be called", name)
				return true
			})
		}
	})
}

func TestGetOrInsert(t *testing.T) {
	t.Run("returns existing value", func(t *testing.T) {
		m := NewMap[string, int]()
		m.Set("key", 42)

		if val := m.GetOrInsert("key", 100); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("inserts and returns new value for missing key", func(t *testing.T) {
		m := NewMap[string, int]()

		if val := m.GetOrInsert("key", 42); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}

		if !m.Has("key") {
			t.Error("expected key to exist after GetOrInsert")
		}
		if val := m.GetOrInsert("key", 100); val != 42 {
			t.Errorf("expected 42 (existing), got %d", val)
		}
	})

	t.Run("initializes nil map", func(t *testing.T) {
		var m Map[string, int]
		if val := m.GetOrInsert("key", 42); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
		if val := m.GetOrInsert("key", 100); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("preserves order of insertions", func(t *testing.T) {
		m := NewMap[string, int]()
		m.GetOrInsert("first", 1)
		m.GetOrInsert("second", 2)
		m.GetOrInsert("third", 3)

		if m.Size() != 3 {
			t.Errorf("expected size 3, got %d", m.Size())
		}
	})
}

func TestGetOrInsertComputed(t *testing.T) {
	t.Run("returns existing value", func(t *testing.T) {
		m := NewMap[string, int]()
		m.Set("key", 42)

		if val := m.GetOrInsertComputed("key", func() int { return 100 }); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("computes and inserts new value for missing key", func(t *testing.T) {
		m := NewMap[string, int]()

		if val := m.GetOrInsertComputed("key", func() int { return 42 }); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}

		if val := m.GetOrInsertComputed("key", func() int { return 100 }); val != 42 {
			t.Errorf("expected 42 (existing), got %d", val)
		}
	})

	t.Run("initializes nil map", func(t *testing.T) {
		var m Map[string, int]
		if val := m.GetOrInsertComputed("key", func() int { return 42 }); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
		if val := m.GetOrInsertComputed("key", func() int { return 100 }); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("uses fn only once per key", func(t *testing.T) {
		count := 0
		m := NewMap[string, int]()

		m.GetOrInsertComputed("key", func() int {
			count++
			return 42
		})

		if count != 1 {
			t.Errorf("expected fn to be called once, got %d", count)
		}
	})
}
