// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexesSearch(t *testing.T) {

	items := []struct {
		idx    Indexes
		target uint32
		ret    int
		ok     bool
	}{
		{
			//No.0
			idx:    Indexes{},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.1
			idx:    Indexes{10},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.2
			idx:    Indexes{10},
			target: uint32(10),
			ret:    0,
			ok:     true,
		},
		{
			//No.3
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.4
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(3),
			ret:    1,
			ok:     true,
		},
		{
			//No.5
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(5),
			ret:    2,
			ok:     true,
		},
		{
			//No.6
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.7
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(0),
			ret:    0,
			ok:     true,
		},
		{
			//No.8
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(10),
			ret:    5,
			ok:     true,
		},
		{
			//No.9
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(6),
			ret:    3,
			ok:     true,
		},
		{
			//No.10
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(2),
			ret:    1,
			ok:     true,
		},
	}

	for i, item := range items {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ret, ok := item.idx.Search(item.target)
			assert.Equal(t, item.ret, ret)
			assert.Equal(t, item.ok, ok)
		})
	}
}

func TestIndexesMatch(t *testing.T) {

	items := []struct {
		idx    Indexes
		target uint32
		ret    int
		ok     bool
	}{
		{
			//No.0
			idx:    Indexes{},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.1
			idx:    Indexes{10},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.2
			idx:    Indexes{10},
			target: uint32(10),
			ret:    0,
			ok:     true,
		},
		{
			//No.3
			idx:    Indexes{10},
			target: uint32(100),
			ret:    0,
			ok:     true,
		},
		{
			//No.4
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.5
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(3),
			ret:    1,
			ok:     true,
		},
		{
			//No.6
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(4),
			ret:    1,
			ok:     true,
		},
		{
			//No.7
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.8
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(100),
			ret:    4,
			ok:     true,
		},
		{
			//No.9
			idx:    Indexes{1, 3, 5, 7, 9},
			target: uint32(0),
			ret:    4,
			ok:     true,
		},
		{
			//No.10
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(0),
			ret:    0,
			ok:     true,
		},
		{
			//No.11
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(10),
			ret:    5,
			ok:     true,
		},
		{
			//No.12
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(6),
			ret:    3,
			ok:     true,
		},
		{
			//No.13
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.14
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.15
			idx:    Indexes{0, 2, 4, 6, 8, 10},
			target: uint32(10000),
			ret:    5,
			ok:     true,
		},
	}

	for i, item := range items {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ret, ok := item.idx.Match(item.target)
			assert.Equal(t, item.ret, ret)
			assert.Equal(t, item.ok, ok)
		})
	}
}

func TestIndexesSort(t *testing.T) {
	expert := Indexes{0, 1, 4, 5, 6, 8, 9}
	idx := Indexes{9, 1, 5, 4, 8, 0, 6}
	idx.Sort()
	assert.Equal(t, idx, expert)
}

func TestIndexesLen(t *testing.T) {
	idx := Indexes{9, 1, 5, 4, 8, 0, 6}
	assert.Equal(t, 7, idx.Len())
}

func TestIndexesLess(t *testing.T) {
	idx := Indexes{9, 1, 5, 4, 8, 0, 6}
	assert.Equal(t, false, idx.Less(0, 1))
	assert.Equal(t, true, idx.Less(1, 2))
	assert.Equal(t, false, idx.Less(0, 6))
}

func TestIndexesSwap(t *testing.T) {
	idx := Indexes{9, 1, 5, 4, 8, 0, 6}
	idx.Swap(0, 1)
	assert.Equal(t, uint32(1), idx[0])
	assert.Equal(t, uint32(9), idx[1])

	idx.Swap(0, 6)
	assert.Equal(t, uint32(6), idx[0])
	assert.Equal(t, uint32(1), idx[6])
}

func BenchmarkIndexesMatch(b *testing.B) {
	idx := make(Indexes, 0)
	for i := 0; i < 20000; i++ {
		idx = append(idx, uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Match(uint32(i))
	}
}

func BenchmarkIndexesSearch(b *testing.B) {
	idx := make(Indexes, 0)
	for i := 0; i < 20000; i++ {
		idx = append(idx, uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := b.N % 20000
		idx.Search(uint32(val))
	}
}
