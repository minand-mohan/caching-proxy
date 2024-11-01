/*
* Copyright © 2024 Minand Manomohanan <minand.nell.mohan@gmail.com>
 */
package server

import (
	"fmt"
	"sync"
)

type Response struct {
	Code int
	Body interface{}
}

type Cache struct {
	CacheMutex sync.Mutex
	UrlMap     map[string]*Response
}

func (c *Cache) CheckAndReturnSavedResponse(url string) (*Response, bool) {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()
	if resp, ok := c.UrlMap[url]; ok {
		fmt.Println("Fetching value from cache...")
		return resp, true
	}

	return nil, false
}

func (c *Cache) PutNewEntryInCache(url string, code int, body interface{}) {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()
	resp := &Response{
		Code: code,
		Body: body,
	}
	c.UrlMap[url] = resp
	fmt.Println(c.UrlMap)
}
