package limiter

type Algorithm interface {
	IsAllowed(string, string) bool
	Stop()
}

type Limiter struct {
	algo Algorithm
}

func New(algo Algorithm) *Limiter {
	return &Limiter{
		algo: algo,
	}
}

func (l *Limiter) Stop() {
	l.algo.Stop()
}

func (l *Limiter) Can(id, key string) bool {
	return l.algo.IsAllowed(id, key)
}
