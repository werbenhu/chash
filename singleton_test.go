// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCHashSingletonCreateBucket(t *testing.T) {
	singleton = nil
	bucket, err := CreateBucket("werbenhu1", 2000)
	assert.NotNil(t, singleton)
	assert.Nil(t, err)
	assert.NotNil(t, bucket)
	assert.Equal(t, "werbenhu1", bucket.Name)

	bucket, err = CreateBucket("werbenhu1", 3000)
	assert.Equal(t, ErrBucketExisted, err)
	assert.Nil(t, bucket)
}

func TestCHashSingletonGetBucket(t *testing.T) {
	singleton = nil
	bucket1, err := CreateBucket("werbenhu2", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	bucket1, err = GetBucket("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	bucket2, err := GetBucket("werbenhu3")
	assert.Nil(t, bucket2)
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestCHashSingletonSerialize(t *testing.T) {
	singleton = nil
	CreateBucket("werbenhu1", 2000)
	CreateBucket("werbenhu2", 1000)

	singleton.InsertAgent("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	singleton.InsertAgent("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	singleton.InsertAgent("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	singleton.InsertAgent("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))

	bs, err := Serialize()
	expert := `{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"agents":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"agents":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`

	assert.Nil(t, err)
	assert.Equal(t, expert, string(bs))
}

func TestCHashSingletonRestore(t *testing.T) {
	singleton = nil
	data := []byte(`{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"agents":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"agents":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`)
	err := Restore(data)
	assert.Nil(t, err)

	bucket1, err := singleton.GetBucket("werbenhu1")
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

	bucket2, err := singleton.GetBucket("werbenhu2")
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
	err = Restore(wrongData)
	assert.NotNil(t, err)
}
