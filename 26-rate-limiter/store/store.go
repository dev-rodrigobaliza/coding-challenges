package store

type Store interface {
	Clear()
	Delete(key string)
	Has(key string) bool
	Get(key string) int
	Inc(key string, delta int) bool
	IncAll(delta int)
	Restore(delta int)
	Set(key string, value int)
	Len() int
	GetAllKeys() []string
}