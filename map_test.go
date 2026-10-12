package jss

import (
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
