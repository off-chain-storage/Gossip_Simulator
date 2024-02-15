package hash

import (
	"github.com/ethereum/go-ethereum/crypto"
)

// var sha256Pool = sync.Pool{New: func() interface{} {
// 	return sha256.New()
// }}

func Hash(data []byte) []byte {
	// h, ok := sha256Pool.Get().(hash.Hash)
	// if !ok {
	// 	h = sha256.New()
	// }
	// defer sha256Pool.Put(h)
	// h.Reset()

	// return h.Sum(nil)

	hash := crypto.Keccak256Hash(data)
	return hash.Bytes()
}
