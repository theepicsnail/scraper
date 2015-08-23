package main

import "encoding/json"
import "os"

type Cache struct {
	filename string
	data     map[string]string
}

func NewCache(filename string) *Cache {
	return &Cache{
		filename,
		make(map[string]string),
	}
}

func (c *Cache) Save() error {
	f, err := os.Create(c.filename)
	defer f.Close()

	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(c.data)
}

func (c *Cache) Load() error {
	f, err := os.Open(c.filename)
	defer f.Close()
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(&cache.data)
}

func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Put(key, val string) {
	c.data[key] = val
}
