// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"sort"
)

type Indexes []uint32

func (idx Indexes) Len() int {
	return len(idx)
}

func (idx Indexes) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

func (idx Indexes) Less(i, j int) bool {
	return idx[i] < (idx[j])
}

func (idx Indexes) Sort() {
	sort.Sort(idx)
}

func (idx Indexes) Search(target uint32) (int, bool) {
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
func (idx Indexes) Match(target uint32) (int, bool) {
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
