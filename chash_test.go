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
	assert.NotNil(t, hash.groups)
}

func TestCHashGetGroup(t *testing.T) {
	hash := New()
	group1, err := hash.CreateGroup("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	existing, err := hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, existing)
	assert.Equal(t, existing, group1)

	group2, err := hash.GetGroup("werbenhu2")
	assert.Nil(t, group2)
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashCreateGroup(t *testing.T) {
	hash := New()
	group1, err := hash.CreateGroup("werbenhu1", 2000)
	assert.Nil(t, err)
	existing, err := hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, group1)
	assert.Equal(t, "werbenhu1", group1.Name)
	assert.Equal(t, group1, existing)

	existing, err = hash.CreateGroup("werbenhu1", 3000)
	assert.Equal(t, ErrGroupExisted, err)
	assert.Equal(t, group1, existing)

	existing, err = hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, existing)
	assert.Equal(t, group1, existing)

	group2, err := hash.CreateGroup("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, group2)

	existing, err = hash.GetGroup("werbenhu2")
	assert.NotNil(t, existing)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(hash.groups))
	assert.Equal(t, "werbenhu2", existing.Name)

	assert.Equal(t, group1, hash.groups[group1.Name])
	assert.Equal(t, group2, hash.groups[group2.Name])
}

func TestCHashRemoveGroup(t *testing.T) {
	hash := New()
	group1, err := hash.CreateGroup("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	group2, err := hash.CreateGroup("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, group2)

	assert.Equal(t, 2, len(hash.groups))
	hash.RemoveGroup("werbenhu1")
	assert.Equal(t, 1, len(hash.groups))

	_, err = hash.GetGroup("werbenhu1")
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashRemoveAllGroup(t *testing.T) {
	hash := New()
	group1, err := hash.CreateGroup("werbenhu1", 2000)
	assert.Nil(t, err)
	assert.NotNil(t, group1)

	group2, err := hash.CreateGroup("werbenhu2", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, group2)

	assert.Equal(t, 2, len(hash.groups))
	hash.RemoveAllGroup()
	assert.Equal(t, 0, len(hash.groups))
}

func TestCHashInsert(t *testing.T) {
	hash := New()
	hash.CreateGroup("werbenhu1", 10000)
	err := hash.Insert("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	assert.Nil(t, err)

	err = hash.Insert("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))
	assert.Nil(t, err)

	group1, err := hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, group1)
	assert.Equal(t, group1.Elements, map[string]*Element{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("werbenhu101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("werbenhu102"),
		},
	})

	hash.CreateGroup("werbenhu2", 10000)
	err = hash.Insert("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	assert.Nil(t, err)
	err = hash.Insert("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))
	assert.Nil(t, err)
	err = hash.Insert("werbenhu2", "192.168.2.103:8080", []byte("werbenhu203"))
	assert.Nil(t, err)

	group2, err := hash.GetGroup("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, group2)
	assert.Equal(t, group2.Elements, map[string]*Element{
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

	err = hash.Insert("b3", "192.168.1.101:8080", []byte("werbenhu101"))
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashDeleteElement(t *testing.T) {
	hash := New()
	hash.CreateGroup("werbenhu1", 10000)
	hash.Insert("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	hash.Insert("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	group, err := hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, group)

	assert.Equal(t, 2, len(group.Elements))
	err = hash.Delete("werbenhu1", "192.168.1.101:8080")
	assert.Equal(t, 1, len(group.Elements))
	assert.Nil(t, err)

	err = hash.Delete("werbenhu1", "192.168.1.102:8080")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(group.Elements))

	err = hash.Delete("werbenhu2", "192.168.1.101:8080")
	assert.Equal(t, ErrGroupNotFound, err)
}

func TestCHashMatch(t *testing.T) {
	hash := New()

	key, payload, err := hash.Match("test", "xxx")
	assert.Equal(t, "", key)
	assert.Nil(t, payload)
	assert.Equal(t, ErrGroupNotFound, err)

	group, err := hash.CreateGroup("test", 10000)
	assert.Nil(t, err)

	_, _, err = hash.Match("test", "werbenhuxxxxx")
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoResultMatched, err)

	setKey, setPayload := "192.168.1.100:1883", []byte("werbenhu100")
	group.Insert(setKey, setPayload)
	assert.Equal(t, 10000, len(group.circle))
	assert.Equal(t, 10000, len(group.rows))
	assert.Equal(t, 1, len(group.Elements))

	key, payload, err = hash.Match("test", "werbenhuxxxxx")
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.NotEqual(t, key, "")
	assert.Equal(t, setKey, key)
	assert.Equal(t, setPayload, payload)
}

func TestCHashSerialize(t *testing.T) {
	hash := New()
	hash.CreateGroup("werbenhu1", 2000)
	hash.CreateGroup("werbenhu2", 1000)

	hash.Insert("werbenhu1", "192.168.1.101:8080", []byte("werbenhu101"))
	hash.Insert("werbenhu1", "192.168.1.102:8080", []byte("werbenhu102"))

	hash.Insert("werbenhu2", "192.168.2.101:8080", []byte("werbenhu201"))
	hash.Insert("werbenhu2", "192.168.2.102:8080", []byte("werbenhu202"))

	bs, err := hash.Serialize()
	expert := `{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"elements":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"elements":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`

	assert.Nil(t, err)
	assert.Equal(t, expert, string(bs))
}

func TestCHashRestore(t *testing.T) {
	hash := New()
	data := []byte(`{"werbenhu1":{"name":"werbenhu1","numberOfReplicas":2000,"elements":{"192.168.1.101:8080":{"key":"192.168.1.101:8080","payload":"d2VyYmVuaHUxMDE="},"192.168.1.102:8080":{"key":"192.168.1.102:8080","payload":"d2VyYmVuaHUxMDI="}}},"werbenhu2":{"name":"werbenhu2","numberOfReplicas":1000,"elements":{"192.168.2.101:8080":{"key":"192.168.2.101:8080","payload":"d2VyYmVuaHUyMDE="},"192.168.2.102:8080":{"key":"192.168.2.102:8080","payload":"d2VyYmVuaHUyMDI="}}}}`)
	err := hash.Restore(data)
	assert.Nil(t, err)

	group1, err := hash.GetGroup("werbenhu1")
	assert.Nil(t, err)
	assert.NotNil(t, group1)
	assert.Equal(t, 4000, len(group1.rows))
	assert.Equal(t, group1.Elements, map[string]*Element{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("werbenhu101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("werbenhu102"),
		},
	})

	group2, err := hash.GetGroup("werbenhu2")
	assert.Nil(t, err)
	assert.NotNil(t, group2)
	assert.Equal(t, 2000, len(group2.rows))
	assert.Equal(t, group2.Elements, map[string]*Element{
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
