package store

type Store interface {
	Clear()
	Delete(key string)
	Get(key string) []byte
	Set(key string, value []byte)
}