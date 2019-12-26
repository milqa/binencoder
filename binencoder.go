package binencoder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type Encoder struct {
	w         io.Writer
	byteOrder binary.ByteOrder
}

func NewEncoder(w io.Writer, byteOrder binary.ByteOrder) *Encoder {
	return &Encoder{
		w:         w,
		byteOrder: byteOrder,
	}
}

func (enc *Encoder) Encode(data interface{}, bytesLen int) error {
	if bytesLen == -1 {
		return nil
	}
	var err error
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Array, reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			err = enc.Encode(v.Index(i).Interface(), bytesLen)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			fieldType := v.Type().Field(i)
			tag := decodeTags(fieldType.Tag.Get("len"), bytesLen)
			err = enc.Encode(v.Field(i).Interface(), tag)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		return enc.Encode(v.Elem().Interface(), bytesLen)
	default:
		by, err := encodeBaseType(data)
		if err != nil {
			log.Printf("[encodeBaseType] Error: %s", err)
			return nil
		}
		// fmt.Printf("Data: %s Len: %v\n", data, bytesLen)
		if bytesLen != 0 {
			delta := bytesLen - len(by)
			if delta < 0 {
				return errors.New("StringLenErr")
			}
			byDelta := make([]byte, delta)
			if enc.byteOrder == binary.LittleEndian {
				by = append(by, byDelta...)
			} else {
				by = append(byDelta, by...)
			}

		}
		err = binary.Write(enc.w, enc.byteOrder, by)
		if err != nil {
			return err
		}
	}
	return err
}

func encodeBaseType(data interface{}) ([]byte, error) {
	b := make([]byte, 0)
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Bool:
		if v.Bool() {
			return append(b, 0x01), nil
		}
		return append(b, 0x00), nil

	case reflect.Uint8:
		return append(b, uint8(v.Uint())), nil
	case reflect.Uint16:
		val := uint16(v.Uint())
		return append(b,
			uint8(val),
			uint8(val>>8),
		), nil

	case reflect.Int16:
		val := uint16(v.Int())
		return append(b,
			uint8(val),
			uint8(val>>8),
		), nil

	case reflect.Uint32:
		val := uint32(v.Uint())
		return append(b,
			uint8(val>>0),
			uint8(val>>8),
			uint8(val>>16),
			uint8(val>>24),
		), nil

	case reflect.Int32:
		val := uint32(v.Int())
		return append(b,
			uint8(val>>0),
			uint8(val>>8),
			uint8(val>>16),
			uint8(val>>24),
		), nil

	case reflect.Uint64:
		val := v.Uint()
		return append(b,
			uint8(val>>0),
			uint8(val>>8),
			uint8(val>>16),
			uint8(val>>24),
			uint8(val>>32),
			uint8(val>>40),
			uint8(val>>48),
			uint8(val>>56),
		), nil

	case reflect.Int64:
		val := uint64(v.Int())
		return append(b,
			uint8(val>>0),
			uint8(val>>8),
			uint8(val>>16),
			uint8(val>>24),
			uint8(val>>32),
			uint8(val>>40),
			uint8(val>>48),
			uint8(val>>56),
		), nil

	case reflect.String:
		s := v.String()
		b = append(b, s...)
		return b, nil

	default:
		return b, fmt.Errorf("unsupported type: " + v.Type().Kind().String())
	}
}

func decodeTags(tag string, defaultTag int) int {
	if tag == "-" {
		return -1
	}
	ans, err := strconv.Atoi(tag)
	if err != nil || ans < 0 {
		return defaultTag
	}
	return ans
}
