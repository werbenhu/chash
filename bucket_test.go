// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBucket(t *testing.T) {
	bucket := NewBucket("test", 10000)
	assert.NotNil(t, bucket)
	assert.NotNil(t, bucket.Agents)
	assert.NotNil(t, bucket.circle)
	assert.NotNil(t, bucket.rows)
	assert.Equal(t, "test", bucket.Name)
	assert.Equal(t, 10000, bucket.NumberOfReplicas)
}

func TestBucketInit(t *testing.T) {
	bucket := &Bucket{}
	bucket.Init()
	assert.NotNil(t, bucket.Agents)
	assert.NotNil(t, bucket.circle)
	assert.NotNil(t, bucket.rows)
}

func TestBucketInsert(t *testing.T) {
	bucket := NewBucket("test", 10000)

	key, payload := "192.168.1.100:1883", []byte("werbenhu100")
	bucket.Insert(key, payload)

	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))

	assert.Equal(t, key, bucket.Agents[key].Key)
	assert.Equal(t, payload, bucket.Agents[key].Payload)
}

func TestBucketDelete(t *testing.T) {
	bucket := NewBucket("test", 10000)

	key, payload := "192.168.1.100:1883", []byte("werbenhu100")
	bucket.Insert(key, payload)
	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))
	assert.Equal(t, key, bucket.Agents[key].Key)
	assert.Equal(t, payload, bucket.Agents[key].Payload)

	bucket.Delete(key)
	assert.Equal(t, 0, len(bucket.circle))
	assert.Equal(t, 0, len(bucket.rows))
	assert.Equal(t, 0, len(bucket.Agents))
}

func TestBucketMatch(t *testing.T) {
	bucket := NewBucket("test", 10000)

	setKey, setPayload := "192.168.1.100:1883", []byte("werbenhu100")
	bucket.Insert(setKey, setPayload)
	assert.Equal(t, 10000, len(bucket.circle))
	assert.Equal(t, 10000, len(bucket.rows))
	assert.Equal(t, 1, len(bucket.Agents))

	key, payload, err := bucket.Match("werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")
	assert.Equal(t, setKey, key)
	assert.Equal(t, setPayload, payload)
}

func TestBucketAll(t *testing.T) {
	bucket := NewBucket("test", 10000)
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

func BenchmarkBucketHash(b *testing.B) {
	b.ResetTimer()
	key := "192.168.1.100:1883"
	bucket := NewBucket("test", 10000)
	for i := 0; i < b.N; i++ {
		bucket.hash(bucket.virtualKey(key, i))
	}
}

func BenchmarkBucketMatch(b *testing.B) {
	bucket := NewBucket("test", 10000)
	bucket.Insert("192.168.1.100:1883", []byte("werbenhu100"))
	bucket.Insert("192.168.1.101:1883", []byte("werbenhu101"))

	b.ReportAllocs()
	b.ResetTimer()

	key := "xxxxx"
	for i := 0; i < b.N; i++ {
		bucket.Match(key)
	}
}
