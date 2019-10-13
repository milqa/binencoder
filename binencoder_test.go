package binencoder_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/milQA/binencoder"
)

type inData struct {
	byteOrder binary.ByteOrder
	data      InStruct
}

type outData struct {
	answer []byte
}

type testData struct {
	encoder *binencoder.Encoder
	in      inData
	out     outData
}

type InStruct struct {
	InM       [3]uint16
	InSlice   []uint32
	InPoint   *[]uint32
	InBool    bool
	InBool2   bool
	InUint8   uint8
	InUint16  uint16
	InUint32  uint32
	InInt32   int32
	InUint64  uint64
	InInt64   int64
	InString  string
	InMap     map[int]int
	InString2 string `len:"10"`
	InString3 string `len:"-"`
	InString4 string `len:"3"`
}

var (
	dataForTests = []testData{
		testData{
			in: inData{
				byteOrder: binary.LittleEndian,
				data: InStruct{
					InM:      [3]uint16{1, 2, 3},
					InSlice:  []uint32{4, 5},
					InPoint:  &[]uint32{6, 7},
					InBool:   true,
					InBool2:  false,
					InUint8:  15,
					InUint16: 255,
					InUint32: 255,
					InInt32:  255,
					InUint64: 255,
					InInt64:  255,
					InString: "test",
					InMap: map[int]int{
						2: 2,
						3: 3,
					},
					InString2: "test",
					InString3: "test",
					InString4: "test",
				},
			},
			out: outData{
				answer: []byte{
					1, 0, 2, 0, 3, 0,
					4, 0, 0, 0, 5, 0, 0, 0,
					6, 0, 0, 0, 7, 0, 0, 0,
					1,
					0,
					15,
					255, 0,
					255, 0, 0, 0,
					255, 0, 0, 0,
					255, 0, 0, 0, 0, 0, 0, 0,
					255, 0, 0, 0, 0, 0, 0, 0,
					116, 101, 115, 116,
					116, 101, 115, 116, 0, 0, 0, 0, 0, 0,
				},
			},
		},
		/* need fix
		*	testData{
		*		in: inData{
		*			byteOrder: binary.BigEndian,
		*			data: InStruct{
		*				InBool:   true,
		*				InUint16: 255,
		*			},
		*		},
		*		out: outData{
		*			answer: []byte{
		*				1,
		*				0, 255,
		*			},
		*		},
		*	},
		 */
	}
)

func TestOne(t *testing.T) {
	for _, data := range dataForTests {
		buf := new(bytes.Buffer)
		encoder := binencoder.NewEncoder(buf, data.in.byteOrder)
		err := encoder.Encode(data.in.data, 0)
		equalErr(t, err.Error(), "StringLenErr")
		equalByte(t, buf.Bytes(), data.out.answer)
	}
}

func equalByte(t *testing.T, answerByte, testDataByte []byte) {
	answer, testData := fmt.Sprintf("[% x]", answerByte), fmt.Sprintf("[% x]", testDataByte)
	if answer != testData {
		t.Errorf("We have:\n%s\n got:\n%s\n", testData, answer)
	}
}

func equalErr(t *testing.T, answerErr, testDataErr string) {
	if answerErr != testDataErr {
		t.Errorf("We have:\n%s\n got:\n%s\n", answerErr, testDataErr)
	}
}
