// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"hash/crc32"
	"strconv"
	"sync"
)

type Agent struct {
	Key     string `json:"key"`
	Payload []byte `json:"payload"`
}

type Bucket struct {
	sync.RWMutex
	Name             string            `json:"name"`
	NumberOfReplicas int               `json:"numberOfReplicas"`
	Agents           map[string]*Agent `json:"agents"`

	circle Indexes
	rows   map[uint32]*Agent
}

func NewBucket(name string, replicas int) *Bucket {
	return &Bucket{
		Name:             name,
		NumberOfReplicas: replicas,
		Agents:           make(map[string]*Agent),
		circle:           make(Indexes, 0),
		rows:             make(map[uint32]*Agent),
	}
}

func (b *Bucket) Init() {
	if b.Agents == nil {
		b.Agents = make(map[string]*Agent)
	}
	if b.circle == nil {
		b.circle = make(Indexes, 0)
	}
	if b.rows == nil {
		b.rows = make(map[uint32]*Agent)
	}
}

func (b *Bucket) hash(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (b *Bucket) virtualKey(key string, idx int) string {
	return strconv.Itoa(idx) + key
}

func (b *Bucket) hashAgent(agent *Agent) {
	for i := 0; i < b.NumberOfReplicas; i++ {
		virtualKey := b.virtualKey(agent.Key, i)
		crc := b.hash(virtualKey)
		b.rows[crc] = agent
		b.circle = append(b.circle, crc)
	}
	b.circle.Sort()
}

func (b *Bucket) Insert(key string, payload []byte) error {
	agent := &Agent{Key: key, Payload: payload}
	b.Lock()
	defer b.Unlock()

	if _, ok := b.Agents[agent.Key]; ok {
		return ErrKeyExisted
	}

	b.Agents[agent.Key] = agent
	b.hashAgent(agent)
	return nil
}

func (b *Bucket) Delete(key string) error {
	agent := &Agent{Key: key, Payload: nil}
	b.Lock()
	defer b.Unlock()

	delete(b.Agents, agent.Key)
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

func (b *Bucket) Match(key string) (string, []byte, error) {
	crc := b.hash(key)
	b.RLock()
	defer b.RUnlock()

	if point, ok := b.circle.Match(crc); ok {
		index := uint32(b.circle[point])
		return b.rows[index].Key, b.rows[index].Payload, nil
	}
	return "", nil, ErrNoResultMatched
}
