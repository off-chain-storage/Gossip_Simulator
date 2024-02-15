package hash

import (
	"hash"
	"sync"

	"github.com/minio/sha256-simd"
	"github.com/sirupsen/logrus"
)

var sha256Pool = sync.Pool{New: func() interface{} {
	return sha256.New()
}}

func Hash(data []byte) []byte {
	logrus.Info("@@ Before Hashing @@", len(data))

	h, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		h = sha256.New()
	}
	defer sha256Pool.Put(h)
	h.Reset()

	logrus.Info("@@ After Hashing @@", len(h.Sum(nil)))

	return h.Sum(nil)
}
