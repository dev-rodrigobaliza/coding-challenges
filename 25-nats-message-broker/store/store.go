package store

type Store[D any] interface {
	Clear()
	Delete(key string)
	Get(key string) D
	Set(key string, value D)
	Len() int
}