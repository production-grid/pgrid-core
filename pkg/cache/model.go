package cache

//Provider defines the cache provider contract
type Provider interface {
	Put(key string, payload interface{}) error
	PutWithExpiration(key string, payload interface{}, exp int32) error
	Add(key string, payload interface{}) error
	AddWithExpiration(key string, payload interface{}, exp int32) error
	IsPresent(key string) (bool, error)
	Get(key string, resultPtr interface{}) (bool, error)
	Delete(key string) error
}
