// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {
}

func (t *testHandler) OnAgentInsert(name string, agent *Agent) error {
	return nil
}

func (t *testHandler) OnAgentDelete(name string, agent *Agent) error {
	return nil
}

func (t *testHandler) OnBucketCreate(bucket *Bucket) error {
	return nil
}

func (t *testHandler) OnBucketRemove(bucket *Bucket) error {
	return nil
}

func TestNewBucket(t *testing.T) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)
	assert.NotNil(t, bucket)
	assert.NotNil(t, bucket.Agents)
	assert.NotNil(t, bucket.circle)
	assert.NotNil(t, bucket.rows)
	assert.Equal(t, "test", bucket.Name)
	assert.Equal(t, 10000, bucket.NumberOfReplicas)
	assert.Equal(t, handler, bucket.handler)
}

func TestBucketInit(t *testing.T) {
	bucket := &Bucket{}
	bucket.Init()
	assert.NotNil(t, bucket.Agents)
	assert.NotNil(t, bucket.circle)
	assert.NotNil(t, bucket.rows)
}

func TestBucketInsert(t *testing.T) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)

	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))

	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))
}

func TestBucketDelete(t *testing.T) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)
	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))

	bucket.Delete("192.168.1.100:1883")
	assert.Equal(t, 0, len(bucket.circle))
	assert.Equal(t, 0, len(bucket.rows))
	assert.Equal(t, 0, len(bucket.Agents))
}

func TestBucketMatch(t *testing.T) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)
	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))

	key, payload, err := bucket.Match("werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")
	assert.Equal(t, "192.168.1.100:1883", key)
	assert.Equal(t, []byte("werbenhu100"), payload)
}

func TestBucketAll(t *testing.T) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)
	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))

	bucket.Insert("192.168.1.101:1883", []byte("werbenhu101"))
	assert.Equal(t, 20000, len(bucket.circle))
	assert.Equal(t, 20000, len(bucket.rows))
	assert.Equal(t, 2, len(bucket.Agents))

	key, payload, err := bucket.Match("werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")

	bucket.Delete(key)
	if key == "192.168.1.100:1883" {
		key2, payload2, err := bucket.Match("werbenhuxxxxx")
		assert.Nil(t, err)
		assert.NotNil(t, payload2)
		assert.NotEqual(t, key2, "")
		assert.Equal(t, "192.168.1.101:1883", key2)
		assert.Equal(t, []byte("werbenhu101"), payload2)
	} else {
		key2, payload2, err := bucket.Match("werbenhuxxxxx")
		assert.Nil(t, err)
		assert.NotNil(t, payload2)
		assert.NotEqual(t, key2, "")
		assert.Equal(t, "192.168.1.100:1883", key2)
		assert.Equal(t, []byte("werbenhu100"), payload2)
	}
}

func BenchmarkHash(b *testing.B) {
	b.ResetTimer()
	handler := &testHandler{}
	key := "192.168.1.100:1883"
	bucket := NewBucket("test", 10000, handler)
	for i := 0; i < b.N; i++ {
		bucket.hash(bucket.virtualKey(key, i))
	}
}

func BenchmarkMatch(b *testing.B) {
	handler := &testHandler{}
	bucket := NewBucket("test", 10000, handler)
	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	bucket.Insert("192.168.1.101:1883", []byte("werbenhu101"))

	b.ReportAllocs()
	b.ResetTimer()

	key := "xxxxx"
	for i := 0; i < b.N; i++ {
		bucket.Match(key)
	}
}
