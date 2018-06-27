package storetest

func BuildJobID(ID uint8) [16]byte {
	return [16]byte{ID}
}
