package maprepo

import (
	cmap "github.com/orcaman/concurrent-map"
)

type DBMap interface {
	Get(key string) (interface{}, error)
	Insert(key string, value interface{}) error
	Upsert(key string, value interface{}) error
	Delete(key string) error
	Clear()
}

type Client struct {
	dbMap *cmap.ConcurrentMap
}

func NewClient() *Client {
	dbMap := cmap.New()

	return &Client{
		dbMap: &dbMap,
	}
}

func (c *Client) Get(key string) (interface{}, error) {
	v, exists := c.dbMap.Get(key)
	if !exists {
		return nil, NewNotFoundError()
	}
	return v, nil
}

func (c *Client) Insert(key string, value interface{}) bool {
	return c.dbMap.SetIfAbsent(key, value)
}

func (c *Client) Upsert(key string, value interface{}) {
	c.dbMap.Set(key, value)
}

func (c *Client) Delete(key string) (interface{}, bool) {
	v, exists := c.dbMap.Pop(key)
	if exists {
		return v, true
	}
	return nil, false
}

func (c *Client) CountAll() int {
	return c.dbMap.Count()
}

func (c *Client) Keys() []string {
	return c.dbMap.Keys()
}

func (c *Client) ClearAll() {
	c.dbMap.Clear()
}
