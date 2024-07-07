package tcpip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const (
	i8   uint = 0
	i16  uint = 1
	i32  uint = 2
	i64  uint = 3
	ui8  uint = 10
	ui16 uint = 11
	ui32 uint = 12
	ui64 uint = 13
	f32  uint = 22
	f64  uint = 23
)

func processState(packet *[]byte, isFirst bool) (proceed bool, dtype uint8, err error) {
	/*
		first byte represents message life:
			0: stop
			1: is continue

		in the first packet you will add the data type in the second byte
			A 2 digit base 10 number will represent the data type

			the first digit represents the size
				0: 8bit
				1: 16bit
				2: 32bit
				3: 64bit

			the second digit represets the type:
				0: int
				1: uint
				2: float

			For example:
				4: int64
				24: float64
				13: uint32

			illegal values: 20 and 22

		all following bytes then represent numeric values in the array
	*/

	fb := (*packet)[0]

	if fb == 1 {
		proceed = true
	}

	if isFirst {

		valid := [10]uint8{
			uint8(i8), uint8(i16), uint8(i32), uint8(i64),
			uint8(ui8), uint8(ui16), uint8(ui32), uint8(ui64),
			uint8(f32), uint8(f64),
		}

		sb := (*packet)[1]

		var passed bool

		for _, validNum := range valid {
			if sb == validNum {
				passed = true
				break
			}
		}

		if !passed {
			// TODO: use proper error
			err = errors.New("invalid dtype code")
			return
		}

		dtype = sb
	}

	return
}

func checkBytes(dataSize uint, bytes *[]byte) error {

	exp := math.Log(float64(dataSize)) / math.Log(2)

	if exp < 3 || exp > 6 {
		return fmt.Errorf("%d is not a valid size", dataSize)
	}

	if uint(len(*bytes)*8)%dataSize != 0 {
		return errors.New("improper amount of bytes recieved")
	}

	return nil
}

func extractUint8(
	packet *[]byte,
) (*[]uint8, error) {

	if err := checkBytes(8, packet); err != nil {
		return nil, err
	}

	sequence := make([]uint8, len(*packet))

	for i := 0; i < len(sequence); i++ {
		sequence[i] = uint8((*packet)[i])
	}

	return &sequence, nil
}

func extractInt8(
	packet *[]byte,
) (*[]int8, error) {

	if err := checkBytes(8, packet); err != nil {
		return nil, err
	}

	sequence := make([]int8, len(*packet))

	for i := 0; i < len(sequence); i++ {
		sequence[i] = int8((*packet)[i])
	}

	return &sequence, nil
}

func extractUint16(
	packet *[]byte,
) (*[]uint16, error) {

	if err := checkBytes(16, packet); err != nil {
		return nil, err
	}

	jump := 2

	sequence := make([]uint16, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			binary.LittleEndian.Uint16((*packet)[i:i+jump]),
		)
	}

	return &sequence, nil
}

func extractInt16(
	packet *[]byte,
) (*[]int16, error) {

	if err := checkBytes(16, packet); err != nil {
		return nil, err
	}

	jump := 2

	sequence := make([]int16, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			int16(binary.LittleEndian.Uint16((*packet)[i:i+jump])),
		)
	}

	return &sequence, nil
}

func extractUint32(
	packet *[]byte,
) (*[]uint32, error) {

	if err := checkBytes(32, packet); err != nil {
		return nil, err
	}

	jump := 4

	sequence := make([]uint32, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			binary.LittleEndian.Uint32((*packet)[i:i+jump]),
		)
	}

	return &sequence, nil
}

func extractInt32(
	packet *[]byte,
) (*[]int32, error) {

	if err := checkBytes(32, packet); err != nil {
		return nil, err
	}

	jump := 4

	sequence := make([]int32, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			int32(binary.LittleEndian.Uint32((*packet)[i:i+jump])),
		)
	}

	return &sequence, nil
}

func extractUint64(
	packet *[]byte,
) (*[]uint64, error) {

	if err := checkBytes(64, packet); err != nil {
		return nil, err
	}

	jump := 8

	sequence := make([]uint64, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			binary.LittleEndian.Uint64((*packet)[i:i+jump]),
		)
	}

	return &sequence, nil
}

func extractInt64(
	packet *[]byte,
) (*[]int64, error) {

	if err := checkBytes(64, packet); err != nil {
		return nil, err
	}

	jump := 8

	sequence := make([]int64, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		sequence = append(
			sequence,
			int64(binary.LittleEndian.Uint64((*packet)[i:i+jump])),
		)
	}

	return &sequence, nil
}

func extractFloat32(
	packet *[]byte,
) (*[]float32, error) {

	if err := checkBytes(32, packet); err != nil {
		return nil, err
	}

	jump := 4

	sequence := make([]float32, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		floatAsUint := binary.LittleEndian.Uint32((*packet)[i : i+jump])
		sequence = append(
			sequence,
			math.Float32frombits(floatAsUint),
		)
	}

	return &sequence, nil
}

func extractFloat64(
	packet *[]byte,
) (*[]float64, error) {

	if err := checkBytes(64, packet); err != nil {
		return nil, err
	}

	jump := 8

	sequence := make([]float64, 0, len(*packet)/jump)

	for i := 0; i < len(*packet); i += jump {
		floatAsUint := binary.LittleEndian.Uint64((*packet)[i : i+jump])
		sequence = append(
			sequence,
			math.Float64frombits(floatAsUint),
		)
	}

	return &sequence, nil
}
