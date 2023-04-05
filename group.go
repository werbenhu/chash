// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"hash/crc32"
	"strconv"
	"sync"
)

type Element struct {
	Key     string `json:"key"`
	Payload []byte `json:"payload"`
}

type Group struct {
	sync.RWMutex
	Name             string              `json:"name"`
	NumberOfReplicas int                 `json:"numberOfReplicas"`
	Elements         map[string]*Element `json:"elements"`

	circle Circle
	rows   map[uint32]*Element
}

func NewGroup(name string, replicas int) *Group {
	return &Group{
		Name:             name,
		NumberOfReplicas: replicas,
		Elements:         make(map[string]*Element),
		circle:           make(Circle, 0),
		rows:             make(map[uint32]*Element),
	}
}

func (b *Group) Init() {
	if b.Elements == nil {
		b.Elements = make(map[string]*Element)
	}
	if b.circle == nil {
		b.circle = make(Circle, 0)
	}
	if b.rows == nil {
		b.rows = make(map[uint32]*Element)
	}
}

func (b *Group) hash(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (b *Group) virtualKey(key string, idx int) string {
	return strconv.Itoa(idx) + key
}

func (b *Group) hashElement(element *Element) {
	for i := 0; i < b.NumberOfReplicas; i++ {
		virtualKey := b.virtualKey(element.Key, i)
		crc := b.hash(virtualKey)
		b.rows[crc] = element
		b.circle = append(b.circle, crc)
	}
	b.circle.Sort()
}

func (b *Group) Insert(key string, payload []byte) error {
	element := &Element{Key: key, Payload: payload}
	b.Lock()
	defer b.Unlock()

	if _, ok := b.Elements[element.Key]; ok {
		return ErrKeyExisted
	}

	b.Elements[element.Key] = element
	b.hashElement(element)
	return nil
}

func (b *Group) Delete(key string) error {
	element := &Element{Key: key, Payload: nil}
	b.Lock()
	defer b.Unlock()

	delete(b.Elements, element.Key)
	for i := 0; i < b.NumberOfReplicas; i++ {
		virtualKey := b.virtualKey(key, i)
		crc := b.hash(virtualKey)
		delete(b.rows, crc)

		if val, ok := b.circle.Search(crc); ok {
			b.circle = append(b.circle[:val], b.circle[val+1:]...)
		}
	}
	return nil
}

func (b *Group) Match(key string) (string, []byte, error) {
	crc := b.hash(key)
	b.RLock()
	defer b.RUnlock()

	if point, ok := b.circle.Match(crc); ok {
		index := uint32(b.circle[point])
		return b.rows[index].Key, b.rows[index].Payload, nil
	}
	return "", nil, ErrNoResultMatched
}
