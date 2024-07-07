package tcpip

import (
	"encoding/binary"
	"math"
	"testing"
)

func Test_ProcessState(t *testing.T) {

	t.Run("test proceed false byte", func(t *testing.T) {

		bytes := []byte{
			0, 1, 2, 3,
		}

		proceed, _, err := processState(&bytes, false)

		if err != nil {
			t.Error(err)
		}

		if proceed {
			t.Error("got proceed as true with byte 0")
		}
	})

	t.Run("test proceed true byte", func(t *testing.T) {

		bytes := []byte{
			1, 0, 12,
		}

		proceed, _, err := processState(&bytes, false)

		if err != nil {
			t.Error(err)
		}

		if !proceed {
			t.Error("got proceed as false with byte 0")
		}
	})

	t.Run("test valid dtype", func(t *testing.T) {

		bytes := []byte{
			1, 22, 10, 11,
		}

		_, dtype, err := processState(&bytes, true)

		if err != nil {
			t.Error(err)
		}

		if dtype != 22 {
			t.Error("got proceed as false with byte 0")
		}
	})

	t.Run("test invalid dtype", func(t *testing.T) {

		bytes := []byte{
			1, 42, 10, 11,
		}

		_, _, err := processState(&bytes, true)

		if err == nil {
			t.Error("accepeted invalid dtype")
		}
	})

}

func Test_CheckBytes(t *testing.T) {

	t.Run("test valid datasize", func(t *testing.T) {

		byt := []byte{
			1, 1, 1, 2, 3, 19, 19, 20,
		}
		err := checkBytes(64, &byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}
	})

	t.Run("test invalid datasize", func(t *testing.T) {

		byt := []byte{
			1, 1, 1, 2, 3, 19, 19, 20,
		}
		err := checkBytes(24, &byt)

		if err == nil {
			t.Error("accepted invalid datasize")
		}
	})

	t.Run("test invalid datasize", func(t *testing.T) {

		byt := []byte{
			1, 1, 1, 2, 3, 19, 19,
		}
		err := checkBytes(32, &byt)

		if err == nil {
			t.Error("accepted invalid byte len")
		}
	})
}

func isEq[n comparable](s1, s2 []n) bool {

	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true

}

func Test_ExtractUint8(t *testing.T) {

	t.Run("test uint8", func(t *testing.T) {

		byt := []byte{
			1, 1, 1, 2, 3, 19, 19, 20,
		}

		s1, err := extractUint8(&byt)

		s2 := []uint8{
			1, 1, 1, 2, 3, 19, 19, 20,
		}

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[uint8](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractInt8(t *testing.T) {

	t.Run("test uint8", func(t *testing.T) {

		byt := []byte{
			1, 1, 1, 2, 3, 0xFF, 19, 20,
		}

		s1, err := extractInt8(&byt)

		s2 := []int8{
			1, 1, 1, 2, 3, -1, 19, 20,
		}

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[int8](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractUint16(t *testing.T) {

	t.Run("test uint16", func(t *testing.T) {

		s2 := []uint16{
			300, 400, 500, 10_000, 50_000,
		}

		byt := make([]byte, 0, len(s2)*2)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint16(byt, i)
		}

		s1, err := extractUint16(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[uint16](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractInt16(t *testing.T) {

	t.Run("test int16", func(t *testing.T) {

		s2 := []int16{
			300, 400, 500, -10_000,
		}

		byt := make([]byte, 0, len(s2)*2)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint16(byt, uint16(i))
		}

		s1, err := extractInt16(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[int16](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractUint32(t *testing.T) {

	t.Run("test uint32", func(t *testing.T) {

		s2 := []uint32{
			1_000_000, 100_000_000, 1_000_000_000,
		}

		byt := make([]byte, 0, len(s2)*4)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint32(byt, i)
		}

		s1, err := extractUint32(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[uint32](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractInt32(t *testing.T) {

	t.Run("test int32", func(t *testing.T) {

		s2 := []int32{
			1_000_000, 100_000_000, -1_000_000_000,
		}

		byt := make([]byte, 0, len(s2)*4)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint32(byt, uint32(i))
		}

		s1, err := extractInt32(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[int32](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractUint64(t *testing.T) {

	t.Run("test uint64", func(t *testing.T) {

		s2 := []uint64{
			18446744073709, 18446749073709, 18486749073709,
		}

		byt := make([]byte, 0, len(s2)*8)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint64(byt, i)
		}

		s1, err := extractUint64(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[uint64](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractInt64(t *testing.T) {

	t.Run("test int64", func(t *testing.T) {

		s2 := []int64{
			18446744073709, 18446749073709, -18486749073709,
		}

		byt := make([]byte, 0, len(s2)*8)

		for _, i := range s2 {
			byt = binary.LittleEndian.AppendUint64(byt, uint64(i))
		}

		s1, err := extractInt64(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[int64](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractFloat32(t *testing.T) {

	t.Run("test float32", func(t *testing.T) {

		s2 := []float32{
			1_000_000.12321312, 1_000_000.12321312, -1_000_000.12321312,
		}

		byt := make([]byte, 0, len(s2)*4)

		for _, i := range s2 {
			ui := math.Float32bits(i)
			byt = binary.LittleEndian.AppendUint32(byt, ui)
		}

		s1, err := extractFloat32(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[float32](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}

func Test_ExtractFloat64(t *testing.T) {

	t.Run("test float64", func(t *testing.T) {

		s2 := []float64{
			1_000_000.12321312, 1_000_000.12321312, -1_000_000.12321312,
		}

		byt := make([]byte, 0, len(s2)*8)

		for _, i := range s2 {
			ui := math.Float64bits(i)
			byt = binary.LittleEndian.AppendUint64(byt, ui)
		}

		s1, err := extractFloat64(&byt)

		if err != nil {
			t.Errorf("recieved error %v", err)
		}

		if !isEq[float64](*s1, s2) {
			t.Error("failed to convert slice")
		}
	})
}
