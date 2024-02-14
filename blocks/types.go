package blocks

import (
	"flag-example/crypto/ecdsa/common"
)

type CurieBlock struct {
	dummyData []byte
}

type SignedCurieBlockForOG struct {
	body      *CurieBlock
	signature common.Signature
}

type SignedCurieBlockForNG struct {
	signature common.Signature
}
