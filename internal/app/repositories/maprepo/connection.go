package maprepo

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map"
)

type DBMap interface {
	Get(key string) (interface{}, bool)
	Insert(key string, value interface{})
	Delete(key string)
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

func (c *Client) Get(key string) (interface{}, bool) {
	return c.dbMap.Get(key)
}

func (c *Client) Insert(key string, value interface{}) bool {
	return c.dbMap.SetIfAbsent(key, value)
}

func (c *Client) Delete(key string) (interface{}, bool) {
	v, exists := c.dbMap.Pop(key)
	if exists {
		return v, true
	}
	return nil, false
}

func (c *Client) ClearAll() {
	c.dbMap.Clear()
}

type MapDBError struct {
	msg string
}

func (m MapDBError) Error() string {
	return fmt.Sprint(m.msg)
}

var (
	ErrNotFound       = MapDBError{"entity not found"}
	ErrDuplicateEntry = MapDBError{"entity already exists"}
)
