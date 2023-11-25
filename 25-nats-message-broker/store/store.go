package store

type Store interface {
	Clear()
	Delete(key string)
	Get(key string) string
	Set(key string, value string)
}