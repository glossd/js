package jss

// Entry represents entry key-value object in Javascript.
type Entry[K comparable, V any] struct {
	Key K
	Val V
}

//type Entries[Key comparable, Val any] = Array[E[Key, Val]]

func E[K comparable, V any](key K, value V) Entry[K, V] {
	return Entry[K, V]{Key: key, Val: value}
}

func E2[K comparable, V any](k1 K, v1 V, k2 K, v2 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2))
}

func E3[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3))
}

func E4[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4))
}

func E5[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V, k5 K, v5 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4), E(k5, v5))
}

func E6[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V, k5 K, v5 V, k6 K, v6 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4), E(k5, v5), E(k6, v6))
}

func E7[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V, k5 K, v5 V, k6 K, v6 V, k7 K, v7 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4), E(k5, v5), E(k6, v6), E(k7, v7))
}

func E8[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V, k5 K, v5 V, k6 K, v6 V, k7 K, v7 V, k8 K, v8 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4), E(k5, v5), E(k6, v6), E(k7, v7), E(k8, v8))
}

func E9[K comparable, V any](k1 K, v1 V, k2 K, v2 V, k3 K, v3 V, k4 K, v4 V, k5 K, v5 V, k6 K, v6 V, k7 K, v7 V, k8 K, v8 V, k9 K, v9 V) Array[Entry[K, V]] {
	return ArrayOf(E(k1, v1), E(k2, v2), E(k3, v3), E(k4, v4), E(k5, v5), E(k6, v6), E(k7, v7), E(k8, v8), E(k9, v9))
}
