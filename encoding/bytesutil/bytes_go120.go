package bytesutil

func ToBytes32(x []byte) [32]byte {
	return [32]byte(PadTo(x, 32))
}
