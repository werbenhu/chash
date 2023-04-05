// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHash(t *testing.T) {
	hash := New()
	assert.NotNil(t, hash)
	assert.NotNil(t, hash.buckets)
}

func TestCHashGetBucket(t *testing.T) {
	hash := New()
	bucket1, err := hash.CreateBucket("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	existing, err := hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, existing)
	assert.Equal(t, existing, bucket1)

	bucket2, err := hash.GetBucket("werbenhu2")
	assert.Nil(t, bucket2)
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestCHashCreateBucket(t *testing.T) {
	hash := New()
	bucket1, err := hash.CreateBucket("werbenhu1", 2000)
	assert.Nil(t, err)
	existing, err := hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)
	assert.Equal(t, "werbenhu1", bucket1.Name)
	assert.Equal(t, bucket1, existing)

	existing, err = hash.CreateBucket("werbenhu1", 3000)
	assert.Equal(t, ErrBucketExisted, err)
	assert.Equal(t, bucket1, existing)

	existing, err = hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, existing)
	assert.Equal(t, bucket1, existing)

	bucket2, err := hash.CreateBucket("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)

	existing, err = hash.GetBucket("werbenhu2")
	assert.NotNil(t, existing)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(hash.buckets))
	assert.Equal(t, "werbenhu2", existing.Name)

	assert.Equal(t, bucket1, hash.buckets[bucket1.Name])
	assert.Equal(t, bucket2, hash.buckets[bucket2.Name])
}

func TestCHashRemoveBucket(t *testing.T) {
	hash := New()
	bucket1, err := hash.CreateBucket("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	bucket2, err := hash.CreateBucket("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)

	assert.Equal(t, 2, len(hash.buckets))
	hash.RemoveBucket("werbenhu1")
	assert.Equal(t, 1, len(hash.buckets))

	_, err = hash.GetBucket("werbenhu1")
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestCHashInsert(t *testing.T) {
	hash := New()
	hash.CreateBucket("werbenhu1", 10000)
	err := hash.InsertAgent("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	assert.Nil(t, err)

	err = hash.InsertAgent("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))
	assert.Nil(t, err)

	bucket1, err := hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)
	assert.Equal(t, bucket1.Agents, map[string]*Agent{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("werbenhu101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("werbenhu102"),
		},
	})

	hash.CreateBucket("werbenhu2", 10000)
	err = hash.InsertAgent("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	assert.Nil(t, err)
	err = hash.InsertAgent("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))
	assert.Nil(t, err)
	err = hash.InsertAgent("werbenhu2", "192.168.2.103:8080", []byte("werbenhu203"))
	assert.Nil(t, err)

	bucket2, err := hash.GetBucket("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)
	assert.Equal(t, bucket2.Agents, map[string]*Agent{
		"192.168.2.101:8080": {
			Key:     "192.168.2.101:8080",
			Payload: []byte("werbenhu201"),
		},
		"192.168.2.102:8080": {
			Key:     "192.168.2.102:8080",
			Payload: []byte("werbenhu202"),
		},
		"192.168.2.103:8080": {
			Key:     "192.168.2.103:8080",
			Payload: []byte("werbenhu203"),
		},
	})

	err = hash.InsertAgent("b3", "192.168.1.101:8080", []byte("werbenhu101"))
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestCHashDeleteAgent(t *testing.T) {
	hash := New()
	hash.CreateBucket("werbenhu1", 10000)
	hash.InsertAgent("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	hash.InsertAgent("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	bucket, err := hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket)

	assert.Equal(t, 2, len(bucket.Agents))
	err = hash.DeleteAgent("werbenhu1", "192.168.1.101:8080")
	assert.Equal(t, 1, len(bucket.Agents))
	assert.Nil(t, err)

	err = hash.DeleteAgent("werbenhu1", "192.168.1.102:8080")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bucket.Agents))

	err = hash.DeleteAgent("werbenhu1", "192.168.1.101:8080")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bucket.Agents))
}

func TestCHashSerialize(t *testing.T) {
	hash := New()
	hash.CreateBucket("werbenhu1", 2000)
	hash.CreateBucket("werbenhu2", 1000)

	hash.InsertAgent("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	hash.InsertAgent("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	hash.InsertAgent("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	hash.InsertAgent("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))

	bs, err := hash.Serialize()
	expert := `{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"agents":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"agents":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`

	assert.Nil(t, err)
	assert.Equal(t, expert, string(bs))
}

func TestCHashRestore(t *testing.T) {
	hash := New()
	data := []byte(`{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"agents":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"agents":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`)
	err := hash.Restore(data)
	assert.Nil(t, err)

	bucket1, err := hash.GetBucket("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)
	assert.Equal(t, 4000, len(bucket1.rows))
	assert.Equal(t, bucket1.Agents, map[string]*Agent{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("werbenhu101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("werbenhu102"),
		},
	})

	bucket2, err := hash.GetBucket("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)
	assert.Equal(t, 2000, len(bucket2.rows))
	assert.Equal(t, bucket2.Agents, map[string]*Agent{
		"192.168.2.101:8080": {
			Key:     "192.168.2.101:8080",
			Payload: []byte("werbenhu201"),
		},
		"192.168.2.102:8080": {
			Key:     "192.168.2.102:8080",
			Payload: []byte("werbenhu202"),
		},
	})

	wrongData := append(data, []byte("--werbenhu--")...)
	err = hash.Restore(wrongData)
	assert.NotNil(t, err)
}
