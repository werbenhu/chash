// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircleSearch(t *testing.T) {

	items := []struct {
		idx    Circle
		target uint32
		ret    int
		ok     bool
	}{
		{
			//No.0
			idx:    Circle{},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.1
			idx:    Circle{10},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.2
			idx:    Circle{10},
			target: uint32(10),
			ret:    0,
			ok:     true,
		},
		{
			//No.3
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.4
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(3),
			ret:    1,
			ok:     true,
		},
		{
			//No.5
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(5),
			ret:    2,
			ok:     true,
		},
		{
			//No.6
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.7
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(0),
			ret:    0,
			ok:     true,
		},
		{
			//No.8
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(10),
			ret:    5,
			ok:     true,
		},
		{
			//No.9
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(6),
			ret:    3,
			ok:     true,
		},
		{
			//No.10
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(2),
			ret:    1,
			ok:     true,
		},
		{
			//No.11
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(12),
			ret:    0,
			ok:     false,
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

func TestCircleMatch(t *testing.T) {

	items := []struct {
		idx    Circle
		target uint32
		ret    int
		ok     bool
	}{
		{
			//No.0
			idx:    Circle{},
			target: uint32(1),
			ret:    0,
			ok:     false,
		},
		{
			//No.1
			idx:    Circle{10},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.2
			idx:    Circle{10},
			target: uint32(10),
			ret:    0,
			ok:     true,
		},
		{
			//No.3
			idx:    Circle{10},
			target: uint32(100),
			ret:    0,
			ok:     true,
		},
		{
			//No.4
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.5
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(3),
			ret:    1,
			ok:     true,
		},
		{
			//No.6
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(4),
			ret:    1,
			ok:     true,
		},
		{
			//No.7
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.8
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(100),
			ret:    4,
			ok:     true,
		},
		{
			//No.9
			idx:    Circle{1, 3, 5, 7, 9},
			target: uint32(0),
			ret:    4,
			ok:     true,
		},
		{
			//No.10
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(0),
			ret:    0,
			ok:     true,
		},
		{
			//No.11
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(10),
			ret:    5,
			ok:     true,
		},
		{
			//No.12
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(6),
			ret:    3,
			ok:     true,
		},
		{
			//No.13
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(1),
			ret:    0,
			ok:     true,
		},
		{
			//No.14
			idx:    Circle{0, 2, 4, 6, 8, 10},
			target: uint32(9),
			ret:    4,
			ok:     true,
		},
		{
			//No.15
			idx:    Circle{0, 2, 4, 6, 8, 10},
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

func TestCircleSort(t *testing.T) {
	expert := Circle{0, 1, 4, 5, 6, 8, 9}
	idx := Circle{9, 1, 5, 4, 8, 0, 6}
	idx.Sort()
	assert.Equal(t, idx, expert)
}

func TestCircleLen(t *testing.T) {
	idx := Circle{9, 1, 5, 4, 8, 0, 6}
	assert.Equal(t, 7, idx.Len())
}

func TestCircleLess(t *testing.T) {
	idx := Circle{9, 1, 5, 4, 8, 0, 6}
	assert.Equal(t, false, idx.Less(0, 1))
	assert.Equal(t, true, idx.Less(1, 2))
	assert.Equal(t, false, idx.Less(0, 6))
}

func TestCircleSwap(t *testing.T) {
	idx := Circle{9, 1, 5, 4, 8, 0, 6}
	idx.Swap(0, 1)
	assert.Equal(t, uint32(1), idx[0])
	assert.Equal(t, uint32(9), idx[1])

	idx.Swap(0, 6)
	assert.Equal(t, uint32(6), idx[0])
	assert.Equal(t, uint32(1), idx[6])
}

func BenchmarkCircleMatch(b *testing.B) {
	idx := make(Circle, 0)
	for i := 0; i < 20000; i++ {
		idx = append(idx, uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Match(uint32(i))
	}
}

func BenchmarkCircleSearch(b *testing.B) {
	idx := make(Circle, 0)
	for i := 0; i < 20000; i++ {
		idx = append(idx, uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := b.N % 20000
		idx.Search(uint32(val))
	}
}
