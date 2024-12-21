package argument_encode

import "github.com/holiman/uint256"

func DecodeFuncHash(input []byte) ([]byte, []byte) {
	return input[:4], input[4:]
}

func DecodeStringInput(input []byte, idx int) string {
	// TODO: find out how to update Uint64 to Uint256 for start and len
	start := uint256.NewInt(0).SetBytes(input[idx*32 : idx*32+32]).Uint64()
	len := uint256.NewInt(0).SetBytes(input[start : start+32]).Uint64()
	str := input[start+32 : start+32+len]
	return string(str)
}
