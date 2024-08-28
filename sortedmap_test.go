// Copyright © 2024 Mark Summerfield. All rights reserved.
package sortedmap

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	//tag::api1[]
	//                 012345678
	letters := []rune("ZENZEBRAS")
	tree := SortedMap[rune, int]{}
	for i, letter := range letters {
		tree.Insert(letter, i)
	}
	var out strings.Builder
	for letter, value := range tree.All() {
		out.WriteString(fmt.Sprintf("%c:%d ", letter, value))
	}
	text := strings.TrimSpace(out.String())
	// text: A:7 B:5 E:4 N:2 R:6 S:8 Z:3
	//end::api1[]
	expected := "A:7 B:5 E:4 N:2 R:6 S:8 Z:3"
	if expected != text {
		t.Errorf("expected %q; got %q", expected, text)
	}
	//tag::api2[]
	size := tree.Len()          // 7
	value, ok := tree.Find('Y') // garbage, false
	//end::api2[]
	if size != 7 {
		t.Errorf("expected size 7; got %d (%d)", size, value)
	}
	if ok {
		t.Error("expected false; got true")
	}
	//tag::api3[]
	value, ok = tree.Find('B') // 5, true
	//end::api3[]
	if !ok || value != 5 {
		t.Errorf("expected 5, true; got %t, %d", ok, value)
	}
	size = tree.Len() // size: 7
	//tag::api4[]
	deleted := tree.Delete('Y') // false
	//end::api4[]
	if size != 7 {
		t.Errorf("expected size 7; got %d", size)
	}
	if deleted {
		t.Error("expected false; got true")
	}
	//tag::api5[]
	deleted = tree.Delete('B') // true
	//end::api5[]
	if !deleted {
		t.Error("expected true; got false")
	}
	size = tree.Len() // size: 6
	if size != 6 {
		t.Errorf("expected 6; got %d", size)
	}
	n := 0
	//tag::api6[]
	for letter := range tree.Keys() { // A E N R S Z
		//end::api6[]
		n += int(letter)
	}
	//tag::api7[]
	for value := range tree.Values() { // 7 4 2 6 8 3
		//end::api7[]
		n -= value
	}
	tree.Clear()
	size = tree.Len() // size: 0
	if size != 0 {
		t.Errorf("expected 0; got %d", size)
	}
}

func Test1(t *testing.T) {
	data := []struct {
		S string
		I int
	}{
		{"can", 3}, {"in", 8}, {"a", 1}, {"ebony", 5}, {"go", 7},
		{"be", 2}, {"dent", 4}, {"for", 6},
	}
	expected := []string{"a", "be", "can", "dent", "ebony", "for", "go",
		"in"}
	var tree SortedMap[string, int] // <1>
	for _, datum := range data {
		tree.Insert(datum.S, datum.I)
	}
	i := 1
	for word, n := range tree.All() {
		if word != expected[i-1] || n != i {
			t.Errorf("expected %q %d; got %q %d", expected[i-1], i, word, n)
		}
		i++
	}
	i = 1
	for word := range tree.Keys() {
		if word != expected[i-1] {
			t.Errorf("expected %q; got %q", expected[i-1], word)
		}
		i++
	}
}

func Test2(t *testing.T) {
	data := []struct{ K, V int }{
		{3, 3}, {8, 8}, {1, 1}, {5, 5}, {7, 7}, {2, 2}, {4, 4}, {6, 6},
	}
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	tree := SortedMap[int, int]{} // <1>
	for _, datum := range data {
		tree.Insert(datum.K, datum.V)
	}
	i := 1
	for key, n := range tree.All() {
		if key != expected[i-1] || n != i {
			t.Errorf("expected %d %d; got %d %d", expected[i-1], i, key, n)
		}
		i++
	}
	i = 1
	for key := range tree.Keys() {
		if key != expected[i-1] {
			t.Errorf("expected %d; got %d", expected[i-1], key)
		}
		i++
	}
}

func TestStringKeyInsertion(t *testing.T) {
	var wordForWord SortedMap[string, string]
	for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
		wordForWord.Insert(strings.ToLower(word), word)
	}
	var words []string
	for word := range wordForWord.Values() {
		words = append(words, word)
	}
	actual, expected := strings.Join(words, ""), "FivefouroneTHREETwo"
	if actual != expected {
		t.Errorf("%q != %q", actual, expected)
	}
}

func TestIntKeyFind(t *testing.T) {
	var intMap SortedMap[int, int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intMap.Insert(number, number*10)
	}
	for _, number := range []int{0, 1, 5, 8, 9} {
		value, ok := intMap.Find(number)
		if !ok {
			t.Errorf("failed to find %d", number)
		}
		actual, expected := value, number*10
		if actual != expected {
			t.Errorf("value is %d should be %d", actual, expected)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		_, ok := intMap.Find(number)
		if ok {
			t.Errorf("should not have found %d", number)
		}
	}
}

func TestIntKeyDelete(t *testing.T) {
	var intMap SortedMap[int, int]
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intMap.Insert(number, number*10)
	}
	if intMap.Len() != 10 {
		t.Errorf("map len %d should be 10", intMap.Len())
	}
	length := 9
	for i, number := range []int{0, 1, 5, 8, 9} {
		if deleted := intMap.Delete(number); !deleted {
			t.Errorf("failed to delete %d", number)
		}
		if intMap.Len() != length-i {
			t.Errorf("map len %d should be %d", intMap.Len(), length-i)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		if deleted := intMap.Delete(number); deleted {
			t.Errorf("should not have deleted nonexistent %d", number)
		}
	}
	if intMap.Len() != 5 {
		t.Errorf("map len %d should be 5", intMap.Len())
	}
}

func TestPassing(t *testing.T) {
	var intMap SortedMap[int, int]
	intMap.Insert(7, 7)
	passTree(&intMap, t)
}

func passTree(tree *SortedMap[int, int], t *testing.T) {
	for _, number := range []int{9, 3, 6, 4, 5, 0} {
		tree.Insert(number, number)
	}
	if tree.Len() != 7 {
		t.Errorf("should have %d items", 7)
	}
}

// Thanks to Russ Cox for improving these benchmarks
func BenchmarkFindSuccess(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var intMap SortedMap[int, int]
	for i := range 1000000 {
		intMap.Insert(i, i)
	}
	b.StartTimer() // Time the Find() method succeeding
	for i := range b.N {
		intMap.Find(i % 1e6)
	}
}

func BenchmarkFindFailure(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	intMap := SortedMap[int, int]{}
	for i := range 1000000 {
		intMap.Insert(2*i, i)
	}
	b.StartTimer() // Time the Find() method failing
	for i := range b.N {
		intMap.Find(2*(i%1e6) + 1)
	}
}

func BenchmarkMapInsertion(b *testing.B) {
	m := map[int]int{}
	for i := range 1000000 {
		m[i] = i
	}

}

func BenchmarkMapSortedIteration(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	m := map[int]int{}
	for i := range 1000000 {
		m[i] = i
	}
	b.StartTimer() // Time sort & iterate
	total := 0
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		total += key
	}
	b.StopTimer() // Don't time check
	if total != 499999500000 {
		panic(total)
	}
}

func BenchmarkBTreeInsertion(b *testing.B) {
	var m SortedMap[int, int]
	for i := range 1000000 {
		m.Insert(i, i)
	}
}

func BenchmarkBTreeIteration(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	var m SortedMap[int, int]
	for i := range 1000000 {
		m.Insert(i, i)
	}
	b.StartTimer() // Time sort & iterate
	total := 0
	for key := range m.Keys() {
		total += key
	}
	b.StopTimer() // Don't time check
	if total != 499999500000 {
		panic(total)
	}
}
func Test_DeleteValue(t *testing.T) {
	var tree SortedMap[int, string]
	for _, n := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		tree.Insert(n, strconv.Itoa(n))
	}
	for key, value := range tree.All() {
		k := strconv.Itoa(key)
		if k != value {
			t.Errorf("expected [%d]=%q; got %q", key, k, value)
		}
	}
	if tree.Len() != 10 {
		t.Errorf("Len expected 10; got %d", tree.Len())
	}
	tree.Delete(99)
	if tree.Len() != 10 {
		t.Errorf("Len expected 10; got %d", tree.Len())
	}
	for _, i := range []int{6, 9, 3, 1} {
		tree.Delete(i)
	}
	if tree.Len() != 6 {
		t.Errorf("Len expected 6; got %d", tree.Len())
	}
	for key, value := range tree.All() {
		k := strconv.Itoa(key)
		if k != value {
			t.Errorf("expected [%d]=%q; got %q", key, k, value)
		}
	}
	for key, value := range tree.All() {
		if key == 5 {
			break
		}
		k := strconv.Itoa(key)
		if k != value {
			t.Errorf("expected [%d]=%q; got %q", key, k, value)
		}
	}
	for key := range tree.Keys() {
		if key == 5 {
			break
		}
	}
	for value := range tree.Values() {
		if value == "5" {
			break
		}
	}
}
