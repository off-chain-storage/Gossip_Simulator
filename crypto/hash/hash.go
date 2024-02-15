package hash

import (
	"hash"
	"sync"

	"github.com/minio/sha256-simd"
)

var sha256Pool = sync.Pool{New: func() interface{} {
	return sha256.New()
}}

func Hash(data []byte) []byte {
	h, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		h = sha256.New()
	}
	defer sha256Pool.Put(h)
	h.Reset()

	return h.Sum(nil)
}
