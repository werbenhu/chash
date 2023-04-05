// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGroup(t *testing.T) {
	group := NewGroup("test", 10000)
	assert.NotNil(t, group)
	assert.NotNil(t, group.Agents)
	assert.NotNil(t, group.circle)
	assert.NotNil(t, group.rows)
	assert.Equal(t, "test", group.Name)
	assert.Equal(t, 10000, group.NumberOfReplicas)
}

func TestGroupInit(t *testing.T) {
	group := &Group{}
	group.Init()
	assert.NotNil(t, group.Agents)
	assert.NotNil(t, group.circle)
	assert.NotNil(t, group.rows)
}

func TestGroupInsert(t *testing.T) {
	group := NewGroup("test", 10000)

	key, payload := "192.168.1.100:1883", []byte("werbenhu100")
	err := group.Insert(key, payload)

	assert.Nil(t, err)
	assert.Equal(t, 10000, len(group.circle))
	assert.Equal(t, 10000, len(group.rows))
	assert.Equal(t, 1, len(group.Agents))

	assert.Equal(t, key, group.Agents[key].Key)
	assert.Equal(t, payload, group.Agents[key].Payload)

	err = group.Insert(key, payload)
	assert.Equal(t, ErrKeyExisted, err)
}

func TestGroupDelete(t *testing.T) {
	group := NewGroup("test", 10000)

	key, payload := "192.168.1.100:1883", []byte("werbenhu100")
	group.Insert(key, payload)
	assert.Equal(t, 10000, len(group.circle))
	assert.Equal(t, 10000, len(group.rows))
	assert.Equal(t, 1, len(group.Agents))
	assert.Equal(t, key, group.Agents[key].Key)
	assert.Equal(t, payload, group.Agents[key].Payload)

	group.Delete(key)
	assert.Equal(t, 0, len(group.circle))
	assert.Equal(t, 0, len(group.rows))
	assert.Equal(t, 0, len(group.Agents))
}

func TestGroupMatch(t *testing.T) {
	group := NewGroup("test", 10000)

	setKey, setPayload := "192.168.1.100:1883", []byte("werbenhu100")
	group.Insert(setKey, setPayload)
	assert.Equal(t, 10000, len(group.circle))
	assert.Equal(t, 10000, len(group.rows))
	assert.Equal(t, 1, len(group.Agents))

	key, payload, err := group.Match("werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")
	assert.Equal(t, setKey, key)
	assert.Equal(t, setPayload, payload)
}

func TestGroupAll(t *testing.T) {
	group := NewGroup("test", 10000)
	group.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	assert.Equal(t, 10000, len(group.circle))
	assert.Equal(t, 10000, len(group.rows))
	assert.Equal(t, 1, len(group.Agents))

	group.Insert("192.168.1.101:1883", []byte("werbenhu101"))
	assert.Equal(t, 20000, len(group.circle))
	assert.Equal(t, 20000, len(group.rows))
	assert.Equal(t, 2, len(group.Agents))

	key, payload, err := group.Match("werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")

	group.Delete(key)
	if key == "192.168.1.100:1883" {
		key2, payload2, err := group.Match("werbenhuxxxxx")
		assert.Nil(t, err)
		assert.NotNil(t, payload2)
		assert.NotEqual(t, key2, "")
		assert.Equal(t, "192.168.1.101:1883", key2)
		assert.Equal(t, []byte("werbenhu101"), payload2)
	} else {
		key2, payload2, err := group.Match("werbenhuxxxxx")
		assert.Nil(t, err)
		assert.NotNil(t, payload2)
		assert.NotEqual(t, key2, "")
		assert.Equal(t, "192.168.1.100:1883", key2)
		assert.Equal(t, []byte("werbenhu100"), payload2)
	}
}

func BenchmarkGroupHash(b *testing.B) {
	b.ResetTimer()
	key := "192.168.1.100:1883"
	group := NewGroup("test", 10000)
	for i := 0; i < b.N; i++ {
		group.hash(group.virtualKey(key, i))
	}
}

func BenchmarkGroupMatch(b *testing.B) {
	group := NewGroup("test", 10000)
	group.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	group.Insert("192.168.1.101:1883", []byte("werbenhu101"))

	b.ReportAllocs()
	b.ResetTimer()

	key := "xxxxx"
	for i := 0; i < b.N; i++ {
		group.Match(key)
	}
}
