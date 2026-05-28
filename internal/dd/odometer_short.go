package dd

import "fmt"

// UnmarshalOdometerShort unmarshals an OdometerShort value from a 3-byte slice.
//
// The data type `OdometerShort` is specified in the Data Dictionary,
// Section 2.113.
//
// ASN.1 Definition:
//
//	OdometerShort ::= INTEGER (0..2^24-1)
//
// Binary Layout (3 bytes):
//   - Big-endian unsigned 24-bit integer
//
// The all-ones value (`0xFFFFFF`, i.e. 2^24-1) is the spec-defined "value
// not available" sentinel and is returned as `0` so it cannot be mistaken
// for a real odometer reading.
func (opts UnmarshalOptions) UnmarshalOdometerShort(data []byte) (int32, error) {
	const lenOdometerShort = 3
	if len(data) != lenOdometerShort {
		return 0, fmt.Errorf("invalid data length for OdometerShort: got %d, want %d", len(data), lenOdometerShort)
	}
	v := int32(data[0])<<16 | int32(data[1])<<8 | int32(data[2])
	if v == 0xFFFFFF {
		return 0, nil
	}
	return v, nil
}
