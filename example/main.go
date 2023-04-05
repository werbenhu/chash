// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package main

import (
	"fmt"

	"github.com/werbenhu/chash"
)

func main() {

	// create a group named db
	dbHash, _ := chash.CreateGroup("db", 10000)

	// insert three mysql server elements
	dbHash.Insert("192.168.1.100:3306", []byte("mysql0-info"))
	dbHash.Insert("192.168.1.101:3306", []byte("mysql1-info"))
	dbHash.Insert("192.168.1.102:3306", []byte("mysql2-info"))

	// create a group named redis
	redisHash, _ := chash.CreateGroup("redis", 10000)

	// insert three redis server elements
	redisHash.Insert("192.168.1.100:6379", []byte("redis0-info"))
	redisHash.Insert("192.168.1.101:6379", []byte("redis1-info"))
	redisHash.Insert("192.168.1.102:6379", []byte("redis2-info"))

	// get the mysql server close to where the user to in the circle by userid
	user1DbHost, user1Dbinfo, err := dbHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched db host:%s, info:%s\n", user1DbHost, user1Dbinfo)

	// get the mysql server close to where the user to in the circle by userid
	user1RedisHost, user1RedisInfo, err := redisHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched redis host:%s, info:%s\n", user1RedisHost, user1RedisInfo)

	// delete a mysql server element from the group's circle
	dbHash.Delete("192.168.1.101:3306")
	user1DbHost, user1Dbinfo, err = dbHash.Match("user-id-1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("user-id-1 matched db host:%s, info:%s\n", user1DbHost, user1Dbinfo)
}
