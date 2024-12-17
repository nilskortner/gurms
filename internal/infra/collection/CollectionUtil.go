package collection

import (
	"gurms/internal/infra/0goutil/concurrentset"
	"gurms/internal/infra/0goutil/entry"
	"gurms/internal/infra/0goutil/immutablemap"
	"gurms/internal/infra/0goutil/immutableset"
	"reflect"
	"sort"
)

// map (hashmap)
// iterator
// linked list
// map
// list
// slice (random access)
// set
// sync map (concurrent hashmap)

// const slice (immutable collection)
// read only access map(immutable map)
// map(T)bool (unified set) mhm

func init() {

}

// region new instance

func getMapCapability(expectedSize int) int {
	return (int)(float32(expectedSize)/0.75 + 1.0)
}

func newSetBySlice[K comparable](values []K) map[K]struct{} {
	set := make(map[K]struct{})
	for _, key := range values {
		set[key] = struct{}{}
	}
	return set
}

// new set collections T ... collections
func newSetByTwoSets[K comparable](set1 map[K]struct{}, set2 map[K]struct{}) map[K]struct{} {
	set := make(map[K]struct{})
	for key := range set1 {
		set[key] = struct{}{}
	}
	for key := range set2 {
		set[key] = struct{}{}
	}
	return set
}

func newSetByManySets[K comparable](sets ...map[K]struct{}) map[K]struct{} {
	set := make(map[K]struct{})
	for _, manyset := range sets {
		for key := range manyset {
			set[key] = struct{}{}
		}
	}
	return set
}

func newConcurrentSet() *concurrentset.ConcurrentSet {
	set := concurrentset.NewConcurrentSet()
	return set
}

func newMapBySlice[K comparable, V any](keys []K, valueMapper func(K) V) map[K]V {
	newmap := make(map[K]V)
	for _, key := range keys {
		newmap[key] = valueMapper(key)
	}
	return newmap
}

func newMapBySet[K comparable, V any](keys map[K]struct{}, valueMapper func(K) V) map[K]V {
	newmap := make(map[K]V)
	for key := range keys {
		newmap[key] = valueMapper(key)
	}
	return newmap
}

func newImmutableMap[K comparable, V any](entrys []entry.Entry) *immutablemap.ImmutableMap {
	imap := immutablemap.NewImmutableMap(entrys)
	return imap
}

func newImmutableSet[K comparable](slice []K) *immutableset.ImmutableSet {
	iset := immutableset.NewImmutableSet(slice)
	return iset
}

// endregion

// region introspection
func isEmpty[K comparable, V any](collection map[K]V) bool {
	return len(collection) == 0
}

func isNotEmpty[K comparable, V any](collection map[K]V) bool {
	return len(collection) != 0
}

func isImmutableSet(set any) bool {
	if _, ok := set.(immutableset.ImmutableSet); ok {
		return true
	}
	return false
}

// endregion

// region contains
func containsAll[K comparable, V any](map1, map2 map[K]V) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if map2Value, ok := map2[key]; !ok {
			return false
		} else {
			if !reflect.DeepEqual(value, map2Value) {
				return false
			}
		}
	}
	return true
}

func containsAllLooseComparison(map1, map2 map[string]any) bool {
	for key, value := range map1 {
		if map2Value, ok := map2[key]; !ok {
			return false
		} else {
			if !areTwoInterfacesLooselyEqual(value, map2Value) {
				return false
			}
		}
	}
	return true
}

func areTwoInterfacesLooselyEqual(actualValue, expectedValue any) bool {
	if actualValue == nil {
		return nil == expectedValue
	} else if expectedValue == nil {
		return false
	}
	if reflect.DeepEqual(expectedValue, actualValue) {
		return true
	}
	// Compare for slices and collections
	if _, ok := expectedValue.([]any); ok {
		return areSliceInterfaceLooselyEqual(getSlice(expectedValue), actualValue)
	} else if expectedValueSet, ok := expectedValue.(map[any]struct{}); ok {
		// Compare for arrays and collections
		return areSetInterfaceLooselyEqual(expectedValueSet, actualValue)
	} else if expectedValueMap, ok := expectedValue.(map[string]any); ok {
		// Compare for maps
		if actualValueMap, ok := actualValue.(map[string]any); ok {
			return containsAllLooseComparison(actualValueMap, expectedValueMap)
		}
	}
	return false
}

func areSliceInterfaceLooselyEqual(value1 []any, value2 any) bool {
	if _, ok := value2.([]any); ok {
		values := getSlice(value2)
		return areSlicesLooselyEqual(value1, values)
	} else if values, ok := value2.(map[any]struct{}); ok {
		return areSetInterfaceLooselyEqual(values, value1)
	}
	return false
}

func areSetInterfaceLooselyEqual(values1 map[any]struct{}, values2 any) bool {
	if _, ok := values2.([]any); ok {
		values := getSlice(values2)
		return areSetSliceLooselyEqual(values1, values)
	} else if values, ok := values2.(map[any]struct{}); ok {
		if len(values1) != len(values) {
			return false
		}
		for key := range values1 {
			if values1[key] != values[key] {
				return false
			}
		}
		return true
	}
	return false
}

func areSetSliceLooselyEqual(values1 map[any]struct{}, values2 []any) bool {
	if len(values1) != len(values2) {
		return false
	}
	var i = 0
	for key := range values1 {
		i++
		if !areTwoInterfacesLooselyEqual(values2[i], key) {
			return false
		}
	}
	return true
}

func areSlicesLooselyEqual(values1, values2 []any) bool {
	if len(values1) != len(values2) {
		return false
	}
	for i := range values1 {
		if !areTwoInterfacesLooselyEqual(values1[i], values2[i]) {
			return false
		}
	}
	return true
}

// endregion

// region conversion
func sliceToSet[K comparable](slice []K) map[K]struct{} {
	if len(slice) == 0 {
		set := make(map[K]struct{})
		return set
	}
	return newSetBySlice(slice)
}

func sliceToImmutableSet[K comparable](slice []K) *immutableset.ImmutableSet {
	return newImmutableSet(slice)
}

// endregion

// region transform
func transformAsSliceBySet[K comparable, V any](set map[K]struct{}, mapper func(K) V) []V {
	var slice []V
	for key := range set {
		slice = append(slice, mapper(key))
	}
	return slice
}

func transformValuesAsMap[K comparable, V any, R any](values map[K]V, supplier func(V) R) map[K]R {
	result := make(map[K]R)
	for key, value := range values {
		result[key] = supplier(value)
	}
	return result
}

func transformValuesAsMapWithValue[K comparable, V any, R any](values map[K]V, value R) map[K]R {
	result := make(map[K]R)
	for key := range values {
		result[key] = value
	}
	return result
}

func sortSlice(values []int64) []int64 {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}

// endregion

// region add/remove

func addSetToSlice(list []int64, values map[int64]struct{}) []int64 {
	for key := range values {
		list = append(list, key)
	}
	return list
}

func addSetToSet(set map[int64]struct{}, values map[int64]struct{}) map[int64]struct{} {
	for key := range values {
		set[key] = struct{}{}
	}
	return set
}

func removeFromSet(set map[int64]struct{}, value int64) map[int64]struct{} {
	if len(set) == 0 {
		return set
	}
	delete(set, value)
	return set
}

func removeSetFromSet(set map[int64]struct{}, values map[int64]struct{}) map[int64]struct{} {
	if len(set) == 0 {
		return set
	}
	for key := range values {
		delete(set, key)
	}
	return set
}

func addToMap() {

}

// endregion

// region merge
func merge[K comparable, V any](map1, map2 map[K]V) map[K]V {
	result := make(map[K]V)
	for key, value := range map1 {
		result[key] = value
	}
	for key, value := range map2 {
		result[key] = value
	}
	return result
}

func deepMerge() {

}

// endregion

// region intersection/union
func intersection[K comparable](set1, set2 map[K]struct{}) map[K]struct{} {
	result := make(map[K]struct{})
	for key := range set2 {
		if _, ok := set1[key]; ok {
			result[key] = struct{}{}
		}
	}
	return result
}

func unionTwoSlices[K comparable](list1, list2 []K) []K {
	var length = len(list1)
	result := make([]K, length+len(list2))
	copy(list1, result)
	copy(list2, result[length:])
	return result
}

func UnionThreeSlices[K comparable](list1, list2, list3 []K) []K {
	var length = len(list1) + len(list2)
	result := make([]K, length+len(list3))
	copy(list1, result)
	copy(list2, result[len(list1):])
	copy(list3, result[length:])
	return result
}

func unionSet() {

}

// region slice

func RemoveByValue[T comparable](slice []T, value T) []T {
	result := []T{}
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
