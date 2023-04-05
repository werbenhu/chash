// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"sort"
)

type Circle []uint32

func (idx Circle) Len() int {
	return len(idx)
}

func (idx Circle) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

func (idx Circle) Less(i, j int) bool {
	return idx[i] < (idx[j])
}

func (idx Circle) Sort() {
	sort.Sort(idx)
}

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

// Match returns an element close to where key hashes to in the circle.
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
