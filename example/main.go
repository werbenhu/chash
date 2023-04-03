// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package main

import (
	"fmt"

	"github.com/werbenhu/chash"
)

func main() {

	dbHash, err := chash.CreateBucket("db-consistent-hash", 10000)
	if err != nil {
		panic(err)
	}
	dbHash.Insert("192.168.1.100:3306", []byte("mysql0-info"))
	dbHash.Insert("192.168.1.101:3306", []byte("mysql1-info"))
	dbHash.Insert("192.168.1.102:3306", []byte("mysql2-info"))

	redisHash, err := chash.CreateBucket("redis-consistent-hash", 10000)
	if err != nil {
		panic(err)
	}
	redisHash.Insert("192.168.1.100:6379", []byte("redis0-info"))
	redisHash.Insert("192.168.1.101:6379", []byte("redis1-info"))
	redisHash.Insert("192.168.1.102:6379", []byte("medis2-info"))

	user1DbHost, user1Dbinfo, err := dbHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched db host:%s, info:%s\n", user1DbHost, user1Dbinfo)

	user1RedisHost, user1RedisInfo, err := redisHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched redis host:%s, info:%s\n", user1RedisHost, user1RedisInfo)

	dbHash.Delete("192.168.1.101:3306")
	user1DbHost, user1Dbinfo, err = dbHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched db host:%s, info:%s\n", user1DbHost, user1Dbinfo)
}
