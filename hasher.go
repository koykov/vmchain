package vmchain

type Hasher interface {
	Sum64(string) uint64
}

const (
	offset64 = uint64(14695981039346656037)
	prime64  = uint64(1099511628211)
)

// FNV-2 implementation
type defaultHasher struct{}

func (defaultHasher) Sum64(p string) uint64 {
	n := len(p)
	if n == 0 {
		return 0
	}
	h := offset64
	_ = p[n-1]
	for i := 0; i < n; i++ {
		h *= prime64
		h ^= uint64(p[i])
	}
	return h
}
