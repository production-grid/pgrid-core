package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

//LocalCache is a simple in memory cache implementation for testing.
//Not recommended for production.  Not scalable and will leak memory.
type LocalCache struct {
	Map map[string]cacheEntry
}

type cacheEntry struct {
	CacheTime  time.Time
	ExpiryTime *time.Time
	Value      []byte
}

func (entry *cacheEntry) isExpired() bool {

	if entry.ExpiryTime == nil {
		return false
	}

	return entry.ExpiryTime.Before(time.Now())

}

func (cache *LocalCache) ensureInitialized() {
	if cache.Map == nil {
		cache.Map = make(map[string]cacheEntry)
	}
}

//Put adds or overrides a value in the cache
func (cache *LocalCache) Put(key string, payload interface{}) error {

	cache.ensureInitialized()

	cache.commit(key, payload, nil)

	return nil

}

func (cache *LocalCache) serialize(value interface{}) ([]byte, error) {

	return json.Marshal(value)

}

func (cache *LocalCache) commit(key string, payload interface{}, expiry *time.Time) error {

	content, err := cache.serialize(payload)

	if err != nil {
		return err
	}

	entry := cacheEntry{
		CacheTime:  time.Now(),
		ExpiryTime: expiry,
		Value:      content,
	}

	cache.Map[key] = entry

	return nil

}

//PutWithExpiration adds or overrides a value in the cache with an expiration
func (cache *LocalCache) PutWithExpiration(key string, payload interface{}, exp int32) error {

	cache.ensureInitialized()

	expiry := time.Now().Add(time.Duration(exp) * time.Second)

	cache.commit(key, payload, &expiry)

	return nil

}

//Add adds a value in the cache
func (cache *LocalCache) Add(key string, payload interface{}) error {

	cache.ensureInitialized()

	present, err := cache.IsPresent(key)

	if err != nil {
		return err
	}

	if present {
		return fmt.Errorf("key %v already cached", key)
	}

	cache.commit(key, payload, nil)

	return nil

}

//IsPresent returns true if the given key is in the cache
func (cache *LocalCache) IsPresent(key string) (bool, error) {

	cache.ensureInitialized()

	cacheEntry, ok := cache.Map[key]

	if ok && !cacheEntry.isExpired() {
		if cacheEntry.isExpired() {
			//delete expired entries once detected
			go cache.Delete(key)
		}
		return true, nil
	}

	return false, nil

}

//Delete removes a value from the cache
func (cache *LocalCache) Delete(key string) error {

	cache.ensureInitialized()

	_, ok := cache.Map[key]

	if ok {
		delete(cache.Map, key)
	}

	return nil

}

//AddWithExpiration adds a value in the cache with an expiration
func (cache *LocalCache) AddWithExpiration(key string, payload interface{}, exp int32) error {

	cache.ensureInitialized()

	present, err := cache.IsPresent(key)

	if err != nil {
		return err
	}

	if present {
		return fmt.Errorf("key %v already cached", key)
	}

	expiry := time.Now().Add(time.Duration(exp) * time.Second)

	cache.commit(key, payload, &expiry)

	return nil

}

//Get returns true if the given key is in the cache
func (cache *LocalCache) Get(key string, resultPtr interface{}) (bool, error) {

	cache.ensureInitialized()

	cacheEntry, ok := cache.Map[key]

	if ok && !cacheEntry.isExpired() {
		if cacheEntry.isExpired() {
			//delete expired entries once detected
			go cache.Delete(key)
		}

		err := json.Unmarshal(cacheEntry.Value, resultPtr)

		if err != nil {
			return false, err
		}

		return true, nil

	}

	return false, nil

}
