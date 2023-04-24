// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"sort"
)

type Circle []uint32

// Len returns the length of the circle.
func (idx Circle) Len() int {
	return len(idx)
}

// Swap swaps the elements at positions i and j in the circle.
func (idx Circle) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

// Less compares the elements at positions i and j in the slice
// and returns true if the element at position i is less than the element at position j.
func (idx Circle) Less(i, j int) bool {
	return idx[i] < (idx[j])
}

// Sort sorts the elements in the slice in ascending order.
func (idx Circle) Sort() {
	sort.Sort(idx)
}

// Search searches for the index of target in the sorted slice.
// It returns the index and true if target is found, or 0 and false otherwise.
func (idx Circle) Search(target uint32) (int, bool) {
	if len(idx) == 0 {
		return 0, false
	}
	f := func(x int) bool {
		return idx[x] >= target
	}
	i := sort.Search(len(idx), f)
	if i >= idx.Len() {
		return 0, false
	}
	if idx[i] != target {
		return 0, false
	}
	return i, true
}

// Match returns the index of the element in the slice that is closest to target.
// If target is greater than all elements in the slice, it returns the last index;
// if target is less than all elements in the slice, it returns the last index.
func (idx Circle) Match(target uint32) (int, bool) {
	if len(idx) == 0 {
		return 0, false
	}
	length := len(idx)
	f := func(x int) bool {
		return idx[x] > target
	}
	i := sort.Search(length, f)

	if i >= length || i == 0 {
		return length - 1, true
	}
	return i - 1, true
}
