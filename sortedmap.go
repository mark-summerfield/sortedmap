// Copyright © 2024 Mark Summerfield. All rights reserved.

// This package provides a generic red-black tree implementation. It is in
// effect a < ordered key-value map. ([TOC])
//
// [TOC]: file:///home/mark/app/golib/doc/index.html
package sortedmap

import "iter"

// Comparable allows only string or integer keys.
type Comparable interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// An SortedMap zero value is usable.
// Create it with statements like these:
//
//	var tree SortedMap[string, int]
//	tree := SortedMap[int, int]{}
type SortedMap[K Comparable, V any] struct {
	root *node[K, V]
	size int
}

type node[K Comparable, V any] struct {
	key         K
	value       V
	red         bool
	left, right *node[K, V]
}

// Insert inserts a new key-value item into the tree and
// returns true; or replaces an existing key-value pair’s
// value if the keys are equal and returns false. For example:
//
//	ok := tree.Insert(key, value).
func (me *SortedMap[K, V]) Insert(key K, value V) bool {
	size := me.size
	me.root = me.insert(me.root, key, value)
	me.root.red = false
	return size == me.size
}

func (me *SortedMap[K, V]) insert(root *node[K, V], key K,
	value V) *node[K, V] {
	if root == nil { // If key was present it would go here
		me.size++
		return &node[K, V]{key: key, value: value, red: true}
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	if key < root.key {
		root.left = me.insert(root.left, key, value)
	} else if key > root.key {
		root.right = me.insert(root.right, key, value)
	} else { // Key already in tree so just replace value
		root.value = value
	}
	return insertRotation(root)
}

func isRed[K Comparable, V any](root *node[K, V]) bool {
	return root != nil && root.red
}

func colorFlip[K Comparable, V any](root *node[K, V]) {
	root.red = !root.red
	if root.left != nil {
		root.left.red = !root.left.red
	}
	if root.right != nil {
		root.right.red = !root.right.red
	}
}

func insertRotation[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	return root
}

func rotateLeft[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	x := root.right
	root.right = x.left
	x.left = root
	x.red = root.red
	root.red = true
	return x
}

func rotateRight[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	x := root.left
	root.left = x.right
	x.right = root
	x.red = root.red
	root.red = true
	return x
}

// Len returns the number of items in the tree.
func (me *SortedMap[K, V]) Len() int { return me.size }

// All is a range function for use as an iterable in a
// for … range loop that returns all of the tree’s
// keys and values, e.g.,
//
//	for key, value := range tree.All()
//
// See also [Keys] and [Values]
func (me *SortedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		all(me.root, yield)
	}
}

func all[K Comparable, V any](root *node[K, V],
	yield func(K, V) bool) bool {
	if root != nil {
		return all(root.left, yield) &&
			yield(root.key, root.value) &&
			all(root.right, yield)
	}
	return true
}

// Keys is a range function for use as an iterable in a
// for … range loop that returns all of the tree’s keys:
//
//	for key := range tree.Keys()
//
// See also [All] and [Values]
func (me *SortedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		keys(me.root, yield)
	}
}

func keys[K Comparable, V any](root *node[K, V],
	yield func(K) bool) bool {
	if root != nil {
		return keys(root.left, yield) &&
			yield(root.key) &&
			keys(root.right, yield)
	}
	return true
}

// Values is a range function for use as an iterable in a
// for … range loop that returns all of the tree’s values:
//
//	for value := range tree.Values()
//
// See also [All] and [Keys]
func (me *SortedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		values(me.root, yield)
	}
}

func values[K Comparable, V any](root *node[K, V],
	yield func(V) bool) bool {
	if root != nil {
		return values(root.left, yield) &&
			yield(root.value) &&
			values(root.right, yield)
	}
	return true
}

// Contains returns true if the key is in the tree and false otherwise.
func (me *SortedMap[K, V]) Contains(key K) bool {
	_, found := me.Find(key)
	return found
}

// Find returns the value and true if the key is in the tree
// or V’s zero value and false otherwise. For example:
//
//	value, ok := tree.Find(key)
func (me *SortedMap[K, V]) Find(key K) (V, bool) {
	var zero V
	root := me.root
	for root != nil {
		if key < root.key {
			root = root.left
		} else if key > root.key {
			root = root.right
		} else {
			return root.value, true
		}
	}
	return zero, false
}

// Delete deletes the key-value item with the given key from the
// tree and returns true, or does nothing and returns false if
// there is no key-value with the given key. For example:
//
//	ok := tree.Delete(key).
//
// See also [Clear]
func (me *SortedMap[K, V]) Delete(key K) bool {
	deleted := false
	if me.root != nil {
		if me.root, deleted = delete_(me.root,
			key); me.root != nil {
			me.root.red = false
		}
	}
	if deleted {
		me.size--
	}
	return deleted
}

func delete_[K Comparable, V any](root *node[K, V], key K) (
	*node[K, V], bool) {
	deleted := false
	if key < root.key {
		if root.left != nil {
			if !isRed(root.left) && !isRed(root.left.left) {
				root = moveRedLeft(root)
			}
			root.left, deleted = delete_(root.left, key)
		}
	} else {
		if isRed(root.left) {
			root = rotateRight(root)
		}
		if key == root.key && root.right == nil {
			// free(root)
			return nil, true
		}
		if root.right != nil {
			root, deleted = deleteRight(root, key)
		}
	}
	return fixUp(root), deleted
}

func moveRedLeft[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	colorFlip(root)
	if root.right != nil && isRed(root.right.left) {
		root.right = rotateRight(root.right)
		root = rotateLeft(root)
		colorFlip(root)
	}
	return root
}

func deleteRight[K Comparable, V any](root *node[K, V], key K) (
	*node[K, V], bool) {
	deleted := false
	if !isRed(root.right) && !isRed(root.right.left) {
		root = moveRedRight(root)
	}
	if key == root.key {
		smallest := first(root.right)
		root.key = smallest.key
		root.value = smallest.value
		root.right = deleteMinimum(root.right)
		deleted = true
	} else {
		root.right, deleted = delete_(root.right, key)
	}
	return root, deleted
}

func moveRedRight[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	colorFlip(root)
	if root.left != nil && isRed(root.left.left) {
		root = rotateRight(root)
		colorFlip(root)
	}
	return root
}

// We do not provide an exported First() method because this
// is an implementation detail.
func first[K Comparable, V any](root *node[K, V]) *node[K, V] {
	for root.left != nil {
		root = root.left
	}
	return root
}

func deleteMinimum[K Comparable, V any](
	root *node[K, V]) *node[K, V] {
	if root.left == nil {
		// free(root)
		return nil
	}
	if !isRed(root.left) && !isRed(root.left.left) {
		root = moveRedLeft(root)
	}
	root.left = deleteMinimum(root.left)
	return fixUp(root)
}

func fixUp[K Comparable, V any](root *node[K, V]) *node[K, V] {
	if isRed(root.right) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	return root
}

// Clear deletes all the tree’s key-value items.
// See also [Delete]
func (me *SortedMap[K, V]) Clear() {
	me.root = nil
	me.size = 0
}
