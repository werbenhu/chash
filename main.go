//
//  @File : main.go
//	@Author : WerBen
//  @Email : 289594665@qq.com
//	@Time : 2021/1/29 10:52 
//	@Desc : TODO ...
//

package main

import (
	"ConsistentHashing/chash"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	cHash := chash.NewCHash()
	cHash.InitNodes([]chash.HashNode{
		{Name: "class_one", Data: make(map[string]interface{})},
		{Name: "class_two", Data: make(map[string]interface{})},
		{Name: "class_three", Data: make(map[string]interface{})},
		{Name: "class_four", Data: make(map[string]interface{})},
	})

	r := gin.Default()
	// print all nodes distributed across the circle
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, cHash)
	})

	// get the hash value of name
	r.GET("/hash", func(c *gin.Context) {
		name := c.Query("name")
		c.JSON(200, cHash.GetHashValue(name))
	})

	// delete node by name
	r.GET("/del", func(c *gin.Context) {
		name := c.Query("name")
		cHash.DelNode(chash.HashNode{
			Name: name,
			Data: make(map[string]interface{}),
		})
		c.JSON(200, cHash)
	})

	// add node
	r.GET("/add", func(c *gin.Context) {
		name := c.Query("name")
		cHash.DelNode(chash.HashNode{
			Name: name,
			Data: make(map[string]interface{}),
		})
		c.JSON(200, cHash)
	})

	// match the node to belong by name
	r.GET("/match", func(c *gin.Context) {
		key := c.Query("name")
		c.JSON(200, cHash.Match(key))
	})

	// match the node to belong by hash value
	r.GET("/match_code", func(c *gin.Context) {
		key := c.Query("code")
		code, _ := strconv.Atoi(key)
		c.JSON(200, cHash.MatchCode(uint32(code)))
	})
	r.Run()
}
