package slugcraft

// Set adds a slug to the in-memory cache.
func (c *Cache) Set(slug string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Store[slug] = 0
}

// Get checks if a slug object exist in the cache
func (c *Cache) Get(slug string) bool {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	_, exist := c.Store[slug]
	return exist
}

// Del removes a slug from the cache.
func (c *Cache) Del(slug string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.Store, slug)
}
