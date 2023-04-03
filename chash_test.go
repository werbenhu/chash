package chash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHash(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)
	assert.NotNil(t, hash)
	assert.NotNil(t, hash.buckets)
	assert.Equal(t, handler, hash.handler)
}

func TestGetBucket(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	err := hash.CreateBucket("b1", 2000)
	assert.Nil(t, err)
	bucket1, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	bucket2, err := hash.GetBucket("b2")
	assert.Nil(t, bucket2)
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestCreateBucket(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	err := hash.CreateBucket("b1", 2000)
	assert.Nil(t, err)
	bucket1, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	err = hash.CreateBucket("b1", 3000)
	assert.Equal(t, ErrBucketExisted, err)
	bucket3, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket3)
	assert.Equal(t, bucket1, bucket3)

	err = hash.CreateBucket("b2", 1000)
	assert.Nil(t, err)
	bucket2, err := hash.GetBucket("b2")
	assert.NotNil(t, bucket2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(hash.buckets))

	assert.Equal(t, bucket1, hash.buckets[bucket1.Name])
	assert.Equal(t, bucket2, hash.buckets[bucket2.Name])
}

func TestRemoveBucket(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	err := hash.CreateBucket("b1", 2000)
	assert.Nil(t, err)
	bucket1, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)

	err = hash.CreateBucket("b2", 1000)
	assert.Nil(t, err)
	bucket2, err := hash.GetBucket("b2")
	assert.NotNil(t, bucket2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(hash.buckets))

	hash.RemoveBucket("b1")
	assert.Equal(t, 1, len(hash.buckets))

	_, err = hash.GetBucket("b1")
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestInsertAgent(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	hash.CreateBucket("b1", 10000)
	err := hash.InsertAgent("b1", "192.168.1.101:8080", []byte("101"))
	assert.Nil(t, err)

	err = hash.InsertAgent("b1", "192.168.1.102:8080", []byte("102"))
	assert.Nil(t, err)

	bucket1, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)
	assert.Equal(t, bucket1.Agents, map[string]*Agent{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("102"),
		},
	})

	hash.CreateBucket("b2", 10000)
	err = hash.InsertAgent("b2", "192.168.2.101:8080", []byte("201"))
	assert.Nil(t, err)
	err = hash.InsertAgent("b2", "192.168.2.102:8080", []byte("202"))
	assert.Nil(t, err)
	err = hash.InsertAgent("b2", "192.168.2.103:8080", []byte("203"))
	assert.Nil(t, err)

	bucket2, err := hash.GetBucket("b2")
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)
	assert.Equal(t, bucket2.Agents, map[string]*Agent{
		"192.168.2.101:8080": {
			Key:     "192.168.2.101:8080",
			Payload: []byte("201"),
		},
		"192.168.2.102:8080": {
			Key:     "192.168.2.102:8080",
			Payload: []byte("202"),
		},
		"192.168.2.103:8080": {
			Key:     "192.168.2.103:8080",
			Payload: []byte("203"),
		},
	})

	err = hash.InsertAgent("b3", "192.168.1.101:8080", []byte("101"))
	assert.Equal(t, ErrBucketNotFound, err)
}

func TestDelete(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	hash.CreateBucket("b1", 10000)
	hash.InsertAgent("b1", "192.168.1.101:8080", []byte("101"))
	hash.InsertAgent("b1", "192.168.1.102:8080", []byte("102"))

	bucket, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket)

	assert.Equal(t, 2, len(bucket.Agents))
	err = hash.DeleteAgent("b1", "192.168.1.101:8080")
	assert.Equal(t, 1, len(bucket.Agents))
	assert.Nil(t, err)

	err = hash.DeleteAgent("b1", "192.168.1.102:8080")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bucket.Agents))

	err = hash.DeleteAgent("b1", "192.168.1.101:8080")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bucket.Agents))
}

func TestSerialize(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	hash.CreateBucket("b1", 2000)
	hash.CreateBucket("b2", 1000)

	hash.InsertAgent("b1", "192.168.1.101:8080", []byte("101"))
	hash.InsertAgent("b1", "192.168.1.102:8080", []byte("102"))

	hash.InsertAgent("b2", "192.168.2.101:8080", []byte("201"))
	hash.InsertAgent("b2", "192.168.2.102:8080", []byte("202"))

	bs, err := hash.Serialize()
	expert := `{"b1":{"Name":"b1","NumberOfReplicas":2000,"Agents":{"192.168.1.101:8080":{"Key":"192.168.1.101:8080","Payload":"MTAx"},"192.168.1.102:8080":{"Key":"192.168.1.102:8080","Payload":"MTAy"}}},"b2":{"Name":"b2","NumberOfReplicas":1000,"Agents":{"192.168.2.101:8080":{"Key":"192.168.2.101:8080","Payload":"MjAx"},"192.168.2.102:8080":{"Key":"192.168.2.102:8080","Payload":"MjAy"}}}}`

	assert.Nil(t, err)
	assert.Equal(t, expert, string(bs))
}

func TestRestore(t *testing.T) {
	handler := &testHandler{}
	hash := New()
	hash.SetHandler(handler)

	data := []byte(`{"b1":{"Name":"b1","NumberOfReplicas":2000,"Agents":{"192.168.1.101:8080":{"Key":"192.168.1.101:8080","Payload":"MTAx"},"192.168.1.102:8080":{"Key":"192.168.1.102:8080","Payload":"MTAy"}}},"b2":{"Name":"b2","NumberOfReplicas":1000,"Agents":{"192.168.2.101:8080":{"Key":"192.168.2.101:8080","Payload":"MjAx"},"192.168.2.102:8080":{"Key":"192.168.2.102:8080","Payload":"MjAy"}}}}`)
	err := hash.Restore(data)
	assert.Nil(t, err)

	bucket1, err := hash.GetBucket("b1")
	assert.Nil(t, err)
	assert.NotNil(t, bucket1)
	assert.Equal(t, 4000, len(bucket1.rows))
	assert.Equal(t, bucket1.Agents, map[string]*Agent{
		"192.168.1.101:8080": {
			Key:     "192.168.1.101:8080",
			Payload: []byte("101"),
		},
		"192.168.1.102:8080": {
			Key:     "192.168.1.102:8080",
			Payload: []byte("102"),
		},
	})

	bucket2, err := hash.GetBucket("b2")
	assert.Nil(t, err)
	assert.NotNil(t, bucket2)
	assert.Equal(t, 2000, len(bucket2.rows))
	assert.Equal(t, bucket2.Agents, map[string]*Agent{
		"192.168.2.101:8080": {
			Key:     "192.168.2.101:8080",
			Payload: []byte("201"),
		},
		"192.168.2.102:8080": {
			Key:     "192.168.2.102:8080",
			Payload: []byte("202"),
		},
	})

	wrongData := append(data, []byte("--werbenhu--")...)
	err = hash.Restore(wrongData)
	assert.NotNil(t, err)
}
