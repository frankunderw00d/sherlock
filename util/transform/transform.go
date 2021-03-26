package transform

const (
	FirByte = 1 << 8
	SecByte = 1 << 16
	ThiByte = 1 << 24
)

func IntToBytes(v int) []byte {
	a := make([]byte, 4)

	a[0] = uint8(v / ThiByte)
	a[1] = uint8(v % ThiByte / SecByte)
	a[2] = uint8(v % ThiByte % SecByte / FirByte)
	a[3] = uint8(v % ThiByte % SecByte % FirByte)

	return a
}

func BytesToInt(v []byte) int {
	a := 0

	a += int(v[0]) * ThiByte
	a += int(v[1]) * SecByte
	a += int(v[2]) * FirByte
	a += int(v[3])
	return a
}
